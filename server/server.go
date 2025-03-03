package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/evandrojr/ratelimiter/configs"
	"github.com/evandrojr/ratelimiter/internal/limiter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kr/pretty"
)

func main() {

	configs.LoadConfig()

	limiter.Init()

	r := chi.NewRouter()

	// Middleware padrão do Chi
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Middleware personalizado
	r.Use(RequestLogger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bem-vindo!"))
	})

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Olá, mundo!"))
	})

	log.Println("Servidor rodando em :8080")
	http.ListenAndServe(":8080", r)
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Método: %s, URL: %s, Host: %s", r.Method, r.URL.Path, r.Host)
		pretty.Println(r.Header, r.RemoteAddr, r.Referer())
		segundoRegistrado := time.Now().Unix()
		ip := strings.Split(r.RemoteAddr, ":")[0]
		validaAcesso := limiter.ValidaAcesso(segundoRegistrado, ip, r.Header["Token"][0])
		log.Println(validaAcesso)
		next.ServeHTTP(w, r)
	})
}
