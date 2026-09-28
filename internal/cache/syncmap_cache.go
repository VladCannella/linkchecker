package cache

// TODO (Этап 6): SyncMapCache — обёртка над sync.Map.
// Логическая гонка вида "if !Has(k) { Set(k,v) }" лечится через LoadOrStore.
