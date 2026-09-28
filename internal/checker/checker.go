// Package checker содержит основную логику проверки ссылок:
// Result, worker pool, ограничение скорости и работу с context.
package checker

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// Result — итог проверки одного URL.
type Result struct {
	URL      string
	Status   int
	Duration time.Duration
	Err      error
}

// Options задаёт параметры прогона (Этапы 2-4).
type Options struct {
	Workers    int           // размер пула воркеров (Этап 2)
	RPS        int           // запросов в секунду, 0 = без ограничения (Этап 3)
	PerHost    int           // семафор: макс. одновременных запросов к одному хосту (Этап 3)
	PerRequest time.Duration // таймаут одного запроса (Этап 4)
}

var (
	clientOnce sync.Once
	httpClient *http.Client
)

// getClient лениво инициализирует http.Client с настроенным транспортом.
// TODO (Этап 6): sync.Once для ленивой инициализации http.Client.
func getClient() *http.Client {
	clientOnce.Do(func() {
		httpClient = &http.Client{}
	})
	return httpClient
}

// TODO (Этап 1): func CheckAll(urls []string) []Result — горутина на каждый URL,
// sync.WaitGroup, закрывающая горутина go func(){ wg.Wait(); close(results) }().
func checkone(url string) Result {

	start := time.Now()
	resp, err := http.Get(url)
	duration := time.Since(start)

	if err != nil {
		res := Result{
			URL:      url,
			Duration: duration,
			Err:      err,
		}

		return res
	}

	defer resp.Body.Close()
	res := Result{
		URL:      url,
		Status:   resp.StatusCode,
		Duration: duration,
		Err:      err,
	}
	return res
}

func CheckAll(urls []string) []Result {
	var wg sync.WaitGroup
	var results []Result
	resultsChan := make(chan Result)

	for _, url := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()

			ans := checkone(url)
			resultsChan <- ans

		}()
	}
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for r := range resultsChan {
		results = append(results, r)
	}

	return results

}

// TODO (Этап 2): func RunPool(ctx context.Context, urls []string, opts Options) <-chan Result
// — канал jobs, N воркеров (fan-out), запись в общий results (fan-in).

func RunPool(ctx context.Context, urls []string, opts Options) <-chan Result {
	jobs := make(chan string)
	results := make(chan Result)
	var wg sync.WaitGroup

	go func() {
		for _, url := range urls {
			jobs <- url
		}
		close(jobs)
	}()

	for i := 0; i < opts.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range jobs {
				checkUrl := checkone(url)
				results <- checkUrl
			}

		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results

}

// TODO (Этап 3): семафор chan struct{} на хост + rate limiter на time.Ticker.

// TODO (Этап 4): context.WithTimeout на весь прогон и на каждый запрос отдельно;
// запрос через http.NewRequestWithContext; errors.Is(err, context.DeadlineExceeded/Canceled).

// TODO (Этап 5): select для ожидания результата/отмены, неблокирующей отправки
// в переполненный канал (atomic-счётчик dropped) и таймаута отдельной операции.
