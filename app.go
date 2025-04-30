package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/evandrojr/go-rate-limiter/configs"
	limiterstrategy "github.com/evandrojr/go-rate-limiter/internal/limiter_strategy"
	persistencystrategy "github.com/evandrojr/go-rate-limiter/internal/persistency_strategy"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kr/pretty"
)

const TOKEN_KEY = "APT_KEY"

var estrategiaLimiterEscolhida limiterstrategy.TipoEstrategiaStruct
var estrategiaPersistenciaEscolhida persistencystrategy.TipoEstrategiaStruct

func main() {
	configs.LoadConfig()
	now := time.Now()
	segundoRegistrado := now.Unix()
	var loggerConfig []string
	redisHostname := configs.Config.DBHost
	// loggerConfig = append(loggerConfig, configs.Config.DBPort)

	if os.Getenv("DOCKER_EXECUTION") == "true" {
		redisHostname = "redis-dev"
	}
	loggerConfig = append(loggerConfig, configs.Config.DBHost)

	estrategiaPersistenciaEscolhida.SetStrategy(persistencystrategy.PersistencyStrategyStruct{})
	estrategiaPersistenciaEscolhida.GetStrategy().Init(loggerConfig)

	estrategiaLimiterEscolhida.SetStrategy(limiterstrategy.LimiterStrategyStruct{})
	estrategiaLimiterEscolhida.GetStrategy().Init(segundoRegistrado, configs.Config)

	r := chi.NewRouter()

	// Middleware padrão do Chi
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Middleware personalizado
	r.Use(RequestLimiter)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bem-vindo!"))
	})

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Olá, mundo!"))
	})

	log.Println("Servidor rodando em :8080")
	http.ListenAndServe(":8080", r)
}

func RequestLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Método: %s, URL: %s, Host: %s", r.Method, r.URL.Path, r.Host)
		pretty.Println(r.Header, r.RemoteAddr, r.Referer())
		segundoRegistrado := time.Now().Unix()
		ip := strings.Split(r.RemoteAddr, ":")[0]
		token := ""
		for k, v := range r.Header {
			if strings.EqualFold(k, TOKEN_KEY) {
				token = v[0]
				break
			}
		}
		if err := limiterstrategy.ValidaAcessoPolimorfico(estrategiaLimiterEscolhida.GetStrategy(), segundoRegistrado, ip, token); err != nil {
			log.Println("Acesso negado:", err)
			http.Error(w, err.Error(), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
