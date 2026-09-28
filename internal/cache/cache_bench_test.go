package cache

// TODO (Этап 6): go test -bench для трёх реализаций (MutexCache,
// RWMutexCache, SyncMapCache) на разных соотношениях чтений/записей
// (например 100/0, 90/10, 50/50). Вывод должен показать: sync.Map выигрывает
// только на read-heavy нагрузке с редкой записью, в остальных случаях
// map+RWMutex быстрее. Результаты запиши в README.
