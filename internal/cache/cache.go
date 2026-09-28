// Package cache содержит три реализации одного интерфейса кэша результатов
// по URL (Этап 6): map+Mutex, map+RWMutex и sync.Map — и их бенчмарк.
package cache

import "linkchecker/internal/checker"

// Cache — общий интерфейс для всех трёх реализаций.
type Cache interface {
	Get(url string) (checker.Result, bool)
	Set(url string, r checker.Result)
}

// Правило: структура с sync.Mutex/RWMutex внутри всегда передаётся по
// указателю — копия даст два независимых мьютекса и тихий data race.
