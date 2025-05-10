package main

import (
	"loadBalance/api"
	"loadBalance/config"
	"loadBalance/internal/algorithm"
	"loadBalance/internal/utils/server"
	"log"
	"log/slog"
	"os"
)

const configPath = "./config.yaml"

func main() {
	cfg := config.Init(configPath)

	alg := algorithm.Init(cfg.AlgorithmType, cfg.BackendAddresses)

	mdl := api.RateLimitedMiddleware(cfg.Bucket.Rate, cfg.Bucket.Capacity)

	proxy := api.NewProxyHandler(alg)

	httpServer := server.New(":8080", mdl(proxy))

	log.Println("Started server on port 8080")
	httpServer.Start()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

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
