package business

import (
	"context"
	"log"
	"net/http"

	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/common/utils"
	"github.com/hoangminhphuc/goph-chat/module/user/dto"
	"github.com/hoangminhphuc/goph-chat/module/user/model"
)

type RegisterRepo interface {
	FindUser(ctx context.Context, conditions map[string]interface{}) (*model.User, error)
	CreateUser(ctx context.Context, data *dto.UserRegister) error
}

type Hasher interface {
	Hash(plainPassword string) string 
	Compare(hashedValue, plainText string) bool
}

type registerBusiness struct {
	repo   RegisterRepo
	hasher Hasher
}

func NewRegisterBusiness(repo RegisterRepo, hasher Hasher) registerBusiness {
	return registerBusiness{repo: repo, hasher: hasher}
}

func (rb *registerBusiness) Register(ctx context.Context, data *dto.UserRegister) error {
	user, _ := rb.repo.FindUser(ctx, map[string]interface{}{"email": data.Email})
	if user != nil {
		return common.NewError("email already exists", http.StatusBadRequest)
	}

	var (
		salt string
		err  error
	)

	if salt, err = utils.GenerateSalt(40); err != nil {
		return common.ErrInvalidRequest(err)
	}
	
	log.Println("hello: ", data.Password + salt)
	data.Password = rb.hasher.Hash(data.Password + salt)
	data.Salt = salt
	data.Role = "user"

	if err := rb.repo.CreateUser(ctx, data); err != nil {
		return common.ErrDB(err)
	}

	return nil
}
