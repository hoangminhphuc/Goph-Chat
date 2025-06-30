package business

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/common/models"
	"github.com/hoangminhphuc/goph-chat/module/user/business/mocks"
	"github.com/hoangminhphuc/goph-chat/module/user/dto"
	"github.com/hoangminhphuc/goph-chat/module/user/model"
)

func TestLoginBusiness_Login_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.LoginRepo)
	mockHasher := new(mocks.Hasher)
	expiry := 3600
	secret := "secret"

	user := &model.User{
		BaseModel: models.BaseModel{
			ID: 1,
		},
		Email:    "test@example.com",
		Password: "hashedpassword",
		Salt:     "salt123",
		Role:     model.RoleAdmin,
	}

	loginData := &dto.UserLogin{
		Email:    "test@example.com",
		Password: "password123",
	}

	mockRepo.
		On("FindUser", mock.Anything, map[string]interface{}{"email": loginData.Email}).
		Return(user, nil)

	mockHasher.
		On("Compare", user.Password, loginData.Password+user.Salt).
		Return(true)

	lb := NewLoginBusiness(mockRepo, mockHasher, expiry, secret)

	// Act
	token, err := lb.Login(context.Background(), loginData)

	// Assert
	assert.NoError(t, err)    // check if returning error
	assert.NotEmpty(t, token) // check if token string is not empty

	// Verifies that the mock methods were called exactly as specified.
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestLoginBusiness_Login_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.LoginRepo)
	mockHasher := new(mocks.Hasher)

	loginData := &dto.UserLogin{
		Email:    "notfound@example.com",
		Password: "password123",
	}

	mockRepo.
		On("FindUser", mock.Anything, mock.Anything).
		Return(nil, nil)

	lb := NewLoginBusiness(mockRepo, mockHasher, 3600, "secret")

	token, err := lb.Login(context.Background(), loginData)

	assert.Error(t, err)
	assert.Equal(t, "", token)
	assert.True(t, common.ErrNotFound("user", loginData.Email).Equal(err))
}

func TestLoginBusiness_Login_InvalidPassword(t *testing.T) {
	mockRepo := new(mocks.LoginRepo)
	mockHasher := new(mocks.Hasher)

	user := &model.User{
		BaseModel: models.BaseModel{
			ID: 2,
		},
		Email:    "user@example.com",
		Password: "hashedpassword",
		Salt:     "salt123",
		Role:     model.RoleUser,
	}

	loginData := &dto.UserLogin{
		Email:    "user@example.com",
		Password: "wrongpassword",
	}

	mockRepo.
		On("FindUser", mock.Anything, mock.Anything).
		Return(user, nil)

	mockHasher.
		On("Compare", user.Password, loginData.Password+user.Salt).
		Return(false)

	lb := NewLoginBusiness(mockRepo, mockHasher, 3600, "secret")

	token, err := lb.Login(context.Background(), loginData)

	assert.Error(t, err)
	assert.Equal(t, model.ErrEmailOrPasswordInvalid, err)
	assert.Equal(t, "", token)
}
