// Command testserver запускает локальный тестовый HTTP-сервер для отладки linkcheck.
//
// Ручки:
//
//	/status/{code}            — сразу отвечает указанным кодом
//	/delay/{ms}               — отвечает 200 после задержки в ms миллисекунд
//	/delay/{ms}/status/{code} — задержка + произвольный код
//	/timeout                  — никогда не отвечает
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"linkchecker/internal/testserver"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "listen address")
	flag.Parse()

	srv, err := testserver.NewOn(*addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("testserver listening on", srv.URL)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()

	// рвём клиентские соединения, иначе Close ждёт зависшие обработчики (/timeout)
	srv.CloseClientConnections()
	srv.Close()
}
