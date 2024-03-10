package handler_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/handler/mocks"
	"github.com/SawitProRecruitment/UserService/internal/config"
)

func TestServer_NewServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	s := handler.NewServer(mocks.NewMockAuthUsecase(ctrl), config.Config{})

	assert.NotNil(t, s)
}
