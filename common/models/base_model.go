package models

import (
	"net/http"
	"time"

	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/common/utils"
)

type BaseModel struct {
	ID        int        `json:"-" gorm:"column:id;"`
	FakeID  	string  	`json:"id" gorm:"column:-;"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at;"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;"`
}

func (bModel *BaseModel) Mask() {
	uid, err := utils.EncodeID(uint64(bModel.ID))

	if err != nil {
		common.WrapError(err, "failed to encode ID", http.StatusInternalServerError)
	}

	bModel.FakeID = uid
}