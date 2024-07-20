package database

import (
	"context"
	"github.com/cegielkowski/mba-golang-client-server-api/internal/entity"
	"gorm.io/gorm"
)

type DollarInterface interface {
	Create(ctx context.Context, dollar *entity.Dollar) error
}

type Dollar struct {
	DB *gorm.DB
}

func NewDollar(db *gorm.DB) *Dollar {
	return &Dollar{DB: db}
}

func (d *Dollar) Create(ctx context.Context, dollar *entity.Dollar) error {
	return d.DB.WithContext(ctx).Create(dollar).Error
}
