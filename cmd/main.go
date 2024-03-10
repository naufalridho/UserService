package main

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"os"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/internal/config"
	"github.com/SawitProRecruitment/UserService/repository"
	ucase "github.com/SawitProRecruitment/UserService/usecase"
)

func main() {
	e := echo.New()

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	dbDsn := os.Getenv("DATABASE_URL")
	secretKey := os.Getenv("AUTH_SECRET_KEY")

	cfg := config.Config{
		SecretKey: secretKey,
	}

	dbDriver := "postgres"

	d, err := sql.Open(dbDriver, dbDsn)
	if err != nil {
		panic(err)
	}

	db := sqlx.NewDb(d, dbDriver)

	repo := repository.NewUserRepository(db)
	usecase := ucase.NewAuthUsecaseImpl(repo, cfg)

	return handler.NewServer(usecase, cfg)
}
