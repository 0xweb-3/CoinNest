package handler

import (
	"context"
	"github.com/0xweb-3/CoinNest/proto"
	"github.com/rogpeppe/go-internal/cache"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type EthRepo struct {
	db    *gorm.DB
	cache *cache.Cache
	log   *zap.SugaredLogger
}

func NewUserRepo(db *gorm.DB) *EthRepo {
	return &EthRepo{
		db:  db,
		log: zap.S(),
	}
}

// GetUserById 获取账号信息
func (r *EthRepo) GetUserById(ctx context.Context, userId uint64) (*proto.UserInfo, error) {
	return nil, nil
}
