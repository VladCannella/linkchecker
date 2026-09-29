// Package checker содержит основную логику проверки ссылок:
// Result, worker pool, ограничение скорости и работу с context.
package checker

import (
	"context"
	"net/http"
	"net/url"
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

func RunPool(ctx context.Context, urls []string, opts Options) <-chan Result {

	jobs := make(chan string)
	results := make(chan Result)
	sem := NewSemaphore(opts.PerHost)

	var tickC <-chan time.Time
	var ticker *time.Ticker

	if opts.RPS > 0 {
		ticker = time.NewTicker(time.Second / time.Duration(opts.RPS))
		tickC = ticker.C
	}

	var wg sync.WaitGroup

	go func() {
		defer close(jobs)
		for _, url := range urls {
			select {
			case <-ctx.Done():
				return
			case jobs <- url:
				continue
			}
		}

	}()

	for i := 0; i < opts.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for rawUrl := range jobs {
				host, err := url.Parse(rawUrl)
				if err != nil {
					continue
				}

				if opts.RPS > 0 {
					select {
					case <-tickC:
					case <-ctx.Done():
						continue
					}
				}

				err = sem.Acquire(ctx, host.Host)
				if err != nil {
					continue
				}
				checkUrl := checkone(rawUrl)
				sem.Release(host.Host)

				select {
				case <-ctx.Done():
					return
				case results <- checkUrl:
					continue
				}
			}

		}()
	}

	go func() {
		wg.Wait()
		close(results)
		if ticker != nil {
			defer ticker.Stop()
		}
	}()

	return results

}

type Semaphore struct {
	mu    sync.Mutex
	host  map[string]chan struct{}
	limit int
}

func NewSemaphore(limit int) *Semaphore {
	sem := make(map[string]chan struct{})
	return &Semaphore{
		host:  sem,
		limit: limit,
	}
}

func (s *Semaphore) Acquire(ctx context.Context, host string) error {
	s.mu.Lock()

	sem, ok := s.host[host]
	if !ok {
		sem = make(chan struct{}, s.limit)
		s.host[host] = sem
	}
	s.mu.Unlock()

	select {
	case sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}

}

func (s *Semaphore) Release(host string) {

	s.mu.Lock()
	defer s.mu.Unlock()

	sem, ok := s.host[host]
	if !ok {
		return
	}

	<-sem
}

// TODO (Этап 3): семафор chan struct{} на хост + rate limiter на time.Ticker.

// TODO (Этап 4): context.WithTimeout на весь прогон и на каждый запрос отдельно;
// запрос через http.NewRequestWithContext; errors.Is(err, context.DeadlineExceeded/Canceled).

// TODO (Этап 5): select для ожидания результата/отмены, неблокирующей отправки
// в переполненный канал (atomic-счётчик dropped) и таймаута отдельной операции.
