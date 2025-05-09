package main

import (
	"loadBalance/api"
	"loadBalance/config"
	"loadBalance/internal/algorithm"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
)

const configPath = "./config.yaml"

func main() {
	cfg := config.Init(configPath)

	alg := algorithm.Init(cfg.AlgorithmType, cfg.BackendAddresses)

	//mux := http.NewServeMux()
	//mux.Handle("/api/status")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("Malformed address: %s", r.RemoteAddr)
		}

		log.Printf("Client IP: %s %s %s", ip, r.Method, r.URL.Path)

		target, err := alg.GetNextServer()
		if err != nil {
			log.Printf("alg.GetNextServer: %v", err)
			http.Error(w, "No backend server are alive", http.StatusServiceUnavailable)
			return
		}

		proxy := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = target.Scheme
				req.URL.Host = target.Host
			},
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
				http.Error(w, "Backend unavailable: "+err.Error(), http.StatusServiceUnavailable)
			},
		}
		proxy.ServeHTTP(w, r)
	})

	mdl := api.RateLimitedMiddleware(cfg.Bucket.Rate, cfg.Bucket.Capacity)

	log.Println("Started server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mdl(handler)))
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
