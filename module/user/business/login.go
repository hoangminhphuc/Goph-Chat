package business

import (
	"context"
	"fmt"
	"time"

	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/common/utils"
	"github.com/hoangminhphuc/goph-chat/module/user/dto"
	"github.com/hoangminhphuc/goph-chat/module/user/model"
)

type LoginRepo interface {
	FindUser(ctx context.Context, conditions map[string]interface{}) (*model.User, error)
}

type loginBusiness struct {
	repo LoginRepo
	hasher utils.Hasher
	expiry int
	secret string
}

func NewLoginBusiness (repo LoginRepo, hasher utils.Hasher, expiry int, secret string) *loginBusiness {
	return &loginBusiness{
		repo: repo, 
		hasher: hasher,
		expiry: expiry,
		secret: secret,
	}
}

func (lb *loginBusiness) Login(ctx context.Context, data *dto.UserLogin) (string, error) {
	start := time.Now()
	
	user, _ := lb.repo.FindUser(ctx, map[string]interface{}{"email": data.Email})

	if user == nil {
		return "", common.ErrNotFound("user", data.Email)
	}

	if !lb.hasher.Compare(user.Password, data.Password + user.Salt) {
		return "", model.ErrEmailOrPasswordInvalid
	}

	userPayload := utils.Payload{
		UID: user.ID,
		URole: user.Role.String(),
	}

	token, err :=utils.GenerateToken(userPayload, lb.expiry, lb.secret)

	if err != nil {
		return "", common.ErrInvalidRequest(err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Login took %s\n", elapsed)

	return token, nil


}