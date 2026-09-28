// Package testserver предоставляет локальный httptest.Server с ручками,
// которые отвечают с разной задержкой и разными кодами — чтобы отлаживать
// linkcheck без интернета и без риска получить бан.
package testserver

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"
)

// New поднимает httptest.Server со следующими ручками:
//
//	/status/{code}        — сразу отвечает указанным кодом
//	/delay/{ms}            — отвечает 200 после задержки в ms миллисекунд
//	/delay/{ms}/status/{code} — задержка + произвольный код
//	/timeout                — никогда не отвечает (для проверки per-req таймаута)
func New() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/status/", func(w http.ResponseWriter, r *http.Request) {
		code, err := strconv.Atoi(r.URL.Path[len("/status/"):])
		if err != nil {
			code = http.StatusOK
		}
		w.WriteHeader(code)
	})

	mux.HandleFunc("/delay/", func(w http.ResponseWriter, r *http.Request) {
		ms, err := strconv.Atoi(r.URL.Path[len("/delay/"):])
		if err != nil {
			ms = 0
		}
		time.Sleep(time.Duration(ms) * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/timeout", func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	})

	return httptest.NewServer(mux)
}
