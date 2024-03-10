package fixtures

import "github.com/SawitProRecruitment/UserService/internal/config"

func Config() config.Config {
	return config.Config{
		SecretKey: "QcRPpsGwuHNAoWvOrWmM",
	}
}
