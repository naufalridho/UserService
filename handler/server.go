package handler

import "github.com/SawitProRecruitment/UserService/internal/config"

type Server struct {
	AuthUsecase AuthUsecase
	Config      config.Config
}

func NewServer(au AuthUsecase, cfg config.Config) *Server {
	return &Server{
		AuthUsecase: au,
		Config:      cfg,
	}
}
