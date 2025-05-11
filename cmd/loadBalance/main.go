package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"loadBalance/api"
	"loadBalance/config"
	"loadBalance/internal/algorithm"
	"loadBalance/internal/utils/server"
	"os"
	"os/signal"
	"syscall"
)

const (
	configPath = "./config.yaml"
)

func main() {
	cfg := config.Init(configPath)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	setupLogger(cfg.Logger.Level)

	log.Infof("Initialising algorithm: %s...", cfg.AlgorithmType)
	alg := algorithm.Init(ctx, cfg.AlgorithmType, cfg.BackendAddresses)

	mdl := api.RateLimitedMiddleware(ctx, cfg.Bucket.Rate, cfg.Bucket.Capacity)

	log.Info("Initialising proxy handler...")
	proxy := api.NewProxyHandler(alg)

	log.Info("Started server on port:", cfg.App.Address.Port)
	httpServer := server.New(":"+cfg.App.Address.Port, mdl(proxy))

	select {
	case err := <-httpServer.Notify():
		if err != nil {
			log.Errorf("server exited with error: %v", err)
		} else {
			log.Info("Server exited normally")
		}
	case <-ctx.Done():
		log.Info("Shutdown signal received")
	}

	log.Info("Shutting down HTTP server...")

	if err := httpServer.Shutdown(); err != nil {
		log.Errorf("Error during shutdown: %v", err)
	} else {
		log.Info("Shutdown complete")
	}
}

func setupLogger(level string) {
	logLevel, err := log.ParseLevel(level)
	if err != nil {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(logLevel)
	}

	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	log.SetOutput(os.Stdout)
}

// TODO у меня есть возможность останавливать горутины(токенБакеты, проверка состояния серверов), для корректного закрытия нужно передать управление серверу, чтобы при остановке он закрывал соединения к бд, серверам и работающие на фоне горутины
// TODO проблема состоит в том, что токены генерирует каждый пользователь и нужно будет дать сигнал каждому, так еще и держать это хранилище при себе

// TODO выбрать логер, установить уровень логирования. логирование входящих запросов, событий и ошибок
// TODO tokenBucket: поддерживать настройку лимитов каждого клиента, настройки можно сохранят в бд
// TODO подготовка README с описанием сборки для запуска проекта.
// TODO dockerfile и docker-compose для развертывания сервиса и бд.
// TODO интеграционные тесты с использованием `go test -bench=. -race=.`.
// TODO пример нагрузки через apache bench (ab -n 5000 -c 1000 http://localhost:8080/).

// EXTRA
// TODO несколько алгоритмов(least connection or random)
// TODO Graceful Shutdown.
// TODO CRUD для добавления/удаления клиентов(ip/api-ключей) и настройки лимитов. Пример эндпоинта: POST /clients { "client_id": "user1", "capacity": 100, "rate_per_sec": 10 }
// TODO сохранять состояние клиентов(текущие токены и настройки) в бд
