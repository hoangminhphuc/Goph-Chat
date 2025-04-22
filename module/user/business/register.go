package business

import (
	"context"
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


type registerBusiness struct {
	repo   RegisterRepo
	hasher utils.Hasher
}

func NewRegisterBusiness(repo RegisterRepo, hasher utils.Hasher) *registerBusiness {
	return &registerBusiness{repo: repo, hasher: hasher}
}

func (rb *registerBusiness) Register(ctx context.Context, data *dto.UserRegister) error {
	user, _ := rb.repo.FindUser(ctx, map[string]interface{}{"email": data.Email})

	
	if user != nil {
		return common.NewError("email already exists", http.StatusBadRequest)
	}


	salt, _ := utils.GenerateSalt(40)

	data.Password = rb.hasher.Hash(data.Password + salt)
	data.Salt = salt
	data.Role = "user"

	if err := rb.repo.CreateUser(ctx, data); err != nil {
		return common.ErrDB(err)
	}

	return nil
}
