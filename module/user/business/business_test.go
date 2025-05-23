package business_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/common/models"
	"github.com/hoangminhphuc/goph-chat/common/utils"
	"github.com/hoangminhphuc/goph-chat/module/user/business"
	"github.com/hoangminhphuc/goph-chat/module/user/dto"
	"github.com/hoangminhphuc/goph-chat/module/user/model"
	"github.com/stretchr/testify/assert"
)

// ==================== TESTING FOR RegisterBusiness ====================

// mockRepo implements RegisterRepo for testing
type mockRepo struct {
    foundUser   *model.User
    findErr     error
    createErr   error
    createdData *dto.UserRegister
}

func (m *mockRepo) FindUser(ctx context.Context, conditions map[string]interface{}) (*model.User, error) {
    return m.foundUser, m.findErr
}

func (m *mockRepo) CreateUser(ctx context.Context, data *dto.UserRegister) error {
    m.createdData = data
    return m.createErr
}

// fakeHasher allows predictable hashing and comparison
type fakeHasher struct{}

func (h *fakeHasher) Hash(s string) string {
    return "hashed:" + s
}

func (h *fakeHasher) Compare(hashedValue, plainText string) bool {
    return hashedValue == h.Hash(plainText)
}

func TestRegisterEmailExists(t *testing.T) {
    // Simulate existing user in repo
    repo := &mockRepo{foundUser: &model.User{Email: "a@b.c"}}
    hasher := &fakeHasher{}
    svc := business.NewRegisterBusiness(repo, hasher)

    err := svc.Register(context.Background(), &dto.UserRegister{Email: "a@b.c", Password: "pass"})
    
    // Expect an AppError because email already exists
    assert.Error(t, err)
    appErr, ok := err.(*common.AppError)
    assert.True(t, ok, "error should be *common.AppError type")
    assert.Equal(t, http.StatusBadRequest, appErr.Code)
    assert.Equal(t, "email already exists", appErr.Message)
}

func TestRegisterSuccess(t *testing.T) {
    // Simulate no existing user
    repo := &mockRepo{foundUser: nil}
    hasher := &fakeHasher{}
    svc := business.NewRegisterBusiness(repo, hasher)

    input := &dto.UserRegister{Email: "x@y.z", Password: "secret"}
    err := svc.Register(context.Background(), input)

    // Expect no error
    assert.NoError(t, err)

    // Ensure CreateUser was called and data set
    created := repo.createdData
    assert.NotNil(t, created, "CreateUser was not called")

    // Role should default to "user"
    assert.Equal(t, "user", created.Role)
    // Salt should be greater or equal to 40
    assert.GreaterOrEqual(t, len(created.Salt), 40)

    // Password should use fakeHasher.Hash
    expectedHash := hasher.Hash("secret" + created.Salt)
    assert.Equal(t, expectedHash, created.Password)
}

func TestRegisterCreateFails(t *testing.T) {
    // Simulate DB error on create
    repo := &mockRepo{foundUser: nil, createErr: errors.New("db error")}
    hasher := &fakeHasher{}
    svc := business.NewRegisterBusiness(repo, hasher)

    err := svc.Register(context.Background(), &dto.UserRegister{Email: "u@v.w", Password: "pwd"})
    // Expect an AppError wrapping the DB error
    assert.Error(t, err)
    appErr, ok := err.(*common.AppError)
    assert.True(t, ok, "error should be *common.AppError type")
    assert.Equal(t, http.StatusInternalServerError, appErr.Code)
    // Underlying error should mention "db error"
    assert.Contains(t, appErr.Error(), "db error")
}

func TestBcryptHashCompare(t *testing.T) {
    hasher := utils.NewBcryptHash()
    password := "my-secret-pw"
    // Hash the password
    hashed := hasher.Hash(password)
    
    // Correct password should compare true
    assert.True(t, hasher.Compare(hashed, password))
    // Wrong password should compare false
    assert.False(t, hasher.Compare(hashed, password+"wrong"))
}

// ==================== TESTING FOR LoginBusiness ====================
type mockLoginRepo struct {
	user *model.User
	err  error
}

func (m *mockLoginRepo) FindUser(ctx context.Context, conditions map[string]interface{}) (*model.User, error) {
	return m.user, m.err
}

func TestLoginUserNotFound(t *testing.T) {
	repo := &mockLoginRepo{user: nil}
	hasher := &fakeHasher{}
	biz := business.NewLoginBusiness(repo, hasher, 3600, "secret")

	_, err := biz.Login(context.Background(), &dto.UserLogin{Email: "missing@x.com", Password: "any"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestLoginPasswordMismatch(t *testing.T) {
	repo := &mockLoginRepo{user: &model.User{
		Email: "ok@x.com",
		Password: "hashed:wrong+salt",
		Salt: "salt",
		Role: model.RoleUser,
		BaseModel: models.BaseModel{
        ID: 123,
    },
	}}
	hasher := &fakeHasher{}
	biz := business.NewLoginBusiness(repo, hasher, 3600, "secret")

	_, err := biz.Login(context.Background(), &dto.UserLogin{Email: "ok@x.com", Password: "right"})
	assert.Error(t, err)
	assert.Equal(t, model.ErrEmailOrPasswordInvalid, err)
}

func TestLoginTokenError(t *testing.T) {
	repo := &mockLoginRepo{user: &model.User{
		Email: "ok@x.com",
		Password: "hashed:right+salt",
		Salt: "salt",
		Role: model.RoleUser,
		BaseModel: models.BaseModel{
        ID: 123,
    },
	}}
	hasher := &fakeHasher{}
	biz := business.NewLoginBusiness(repo, hasher, 3600, "") // empty secret might break token gen

	// temporarily patch utils.GenerateToken if needed to simulate failure
	// or use an interface abstraction for token generation

	_, err := biz.Login(context.Background(), &dto.UserLogin{Email: "ok@x.com", Password: "right"})
	assert.Error(t, err)
}

func TestLoginSuccess(t *testing.T) {
    repo := &mockRepo{
        foundUser: &model.User{
            BaseModel: models.BaseModel{
								ID: 1,
						},
            Email:    "a@b.c",
            Salt:     "mysalt",
            Password: "hashed:mypasswordmysalt",
            Role:     model.RoleUser,
        },
    }
    hasher := &fakeHasher{}
    biz := business.NewLoginBusiness(repo, hasher, 3600, "mysecret")

    input := &dto.UserLogin{Email: "a@b.c", Password: "mypassword"}
    token, err := biz.Login(context.Background(), input)

    assert.NoError(t, err)
    assert.NotEmpty(t, token)
}

