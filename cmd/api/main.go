package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vitaliy-ukiru/bank-service/internal/application"
	"github.com/vitaliy-ukiru/bank-service/internal/config"
	"github.com/vitaliy-ukiru/bank-service/internal/infrastructure/repository/accounts"
	"github.com/vitaliy-ukiru/bank-service/internal/transport/webapi"
	"github.com/vitaliy-ukiru/bank-service/internal/transport/webapi/controllers"
	"github.com/vitaliy-ukiru/bank-service/pkg/client/pg"
	"github.com/vitaliy-ukiru/bank-service/pkg/logging"
)

func main() {
	envPath := flag.String("env-path", "", "Path to .env file")
	flag.Parse()

	err := config.LoadConfig(*envPath)
	if err != nil {
		panic(err)
	}

	cfg := config.Get()

	log := logging.New(os.Stdout, cfg.Env == config.EnvDev)

	log.Info("InitConfig", "config initialized", logging.String("env", string(cfg.Env)))

	db, err := pg.New(context.Background(), pg.ConnString(
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Database,
		cfg.Database.Host,
		cfg.Database.Port,
	))

	if err != nil {
		log.Error("InitPostgres", "fail init postgres", err)
		os.Exit(1)
	}

	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Error("PingDatabase", "fail ping database", err)
		os.Exit(1)
	}

	accountsRepository := accounts.NewRepository(db)
	accountService := application.NewAccountService(accountsRepository, accountsRepository)

	accountController := controllers.NewAccountController(accountService)

	apiServer := webapi.New(cfg, accountController, log)

	go func() {
		if err := apiServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("StartServer", "fail run server", err)
			os.Exit(1)
		}
	}()

	{
		quit := make(chan os.Signal, 1)
		signal.Notify(quit,
			syscall.SIGINT,
			syscall.SIGQUIT,
			syscall.SIGKILL,
			syscall.SIGTERM,
		)
		sig := <-quit
		log.Info("shutdown", "shutdown app", logging.String("signal", sig.String()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Error("ShutdownServer", "fail shutdown server", err)
		os.Exit(1)
	}
}
