package service

import (
	"context"
	"github.com/0xweb-3/CoinNest/proto"
	"go.uber.org/zap"
)

type EthRepo interface {
	// GetUserById 获取账号信息
	GetUserById(ctx context.Context, userId uint64) (*proto.UserInfo, error)
}

type EthServer struct {
	proto.UnimplementedUserServer
	userRepo EthRepo
	log      *zap.SugaredLogger
}

func NewEthServer(repo EthRepo) *EthServer {
	return &EthServer{
		userRepo: repo,
		log:      zap.S(),
	}
}

func (s *EthServer) GetUserById(ctx context.Context, req *proto.GetUserByIdReq) (*proto.UserInfo, error) {
	user, err := s.userRepo.GetUserById(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return user, nil
}
