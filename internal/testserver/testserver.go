// Package testserver предоставляет локальный httptest.Server с ручками,
// которые отвечают с разной задержкой и разными кодами — чтобы отлаживать
// linkcheck без интернета и без риска получить бан.
package testserver

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"
)

// handler собирает ручки тестового сервера:
//
//	/status/{code}            — сразу отвечает указанным кодом
//	/delay/{ms}               — отвечает 200 после задержки в ms миллисекунд
//	/delay/{ms}/status/{code} — задержка + произвольный код
//	/timeout                  — никогда не отвечает (для проверки per-req таймаута)
func handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/status/", func(w http.ResponseWriter, r *http.Request) {
		code, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/status/"))
		if err != nil {
			http.Error(w, "bad status code", http.StatusBadRequest)
			return
		}
		w.WriteHeader(code)
	})

	mux.HandleFunc("/delay/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/delay/"), "/")

		ms, err := strconv.Atoi(parts[0])
		if err != nil {
			http.Error(w, "bad delay", http.StatusBadRequest)
			return
		}

		code := http.StatusOK
		if len(parts) == 3 && parts[1] == "status" {
			code, err = strconv.Atoi(parts[2])
			if err != nil {
				http.Error(w, "bad status code", http.StatusBadRequest)
				return
			}
		}

		timer := time.NewTimer(time.Duration(ms) * time.Millisecond)
		defer timer.Stop()
		select {
		case <-timer.C:
			w.WriteHeader(code)
		case <-r.Context().Done():
			// клиент отключился (отмена, таймаут) — незачем держать обработчик
		}
	})

	mux.HandleFunc("/timeout", func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	})

	return mux
}

// New поднимает сервер на случайном свободном порту (удобно для тестов).
func New() *httptest.Server {
	return httptest.NewServer(handler())
}

// NewOn поднимает сервер на заданном адресе, например "127.0.0.1:8080"
// (удобно для ручной отладки CLI).
func NewOn(addr string) (*httptest.Server, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	srv := httptest.NewUnstartedServer(handler())
	srv.Listener.Close()
	srv.Listener = l
	srv.Start()
	return srv, nil
}
