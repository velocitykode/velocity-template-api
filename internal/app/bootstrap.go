package app

import (
	"net/http"
	"os"

	"{{MODULE_NAME}}/config"

	"github.com/joho/godotenv"
	"github.com/velocitykode/velocity/pkg/auth"
	"github.com/velocitykode/velocity/pkg/auth/drivers/guards"
	"github.com/velocitykode/velocity/pkg/crypto"
	"github.com/velocitykode/velocity/pkg/log"
	"github.com/velocitykode/velocity/pkg/orm"
)

func init() {
	godotenv.Load()
}

// Run starts the application
func Run() {
	log.Info("Velocity API Server")

	if err := initialize(); err != nil {
		log.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}

	httpKernel := New()
	httpKernel.Bootstrap()

	port := config.GetPort()
	log.Info("Server starting", "port", port)

	if err := http.ListenAndServe(":"+port, httpKernel.Handler()); err != nil {
		log.Error("Server failed to start", "error", err)
	}
}

// initialize bootstraps all application services
func initialize() error {
	if err := initCrypto(); err != nil {
		return err
	}
	if err := orm.InitFromEnv(); err != nil {
		return err
	}
	return initAuth()
}

func initCrypto() error {
	key := config.GetCryptoKey()
	if key != "" {
		return crypto.Init(crypto.Config{
			Key:    key,
			Cipher: config.GetCryptoCipher(),
		})
	}
	return nil
}

func initAuth() error {
	manager, err := auth.GetManager()
	if err != nil {
		return err
	}

	sessionConfig := auth.NewSessionConfigFromEnv()
	provider := auth.NewORMUserProvider(config.GetAuthModel())
	sessionGuard, err := guards.NewSessionGuard(provider, sessionConfig)
	if err != nil {
		return err
	}

	manager.RegisterGuard(config.GetAuthGuard(), sessionGuard)
	return nil
}
