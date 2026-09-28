// Package server реализует HTTP-сервер с ручкой /check (Этап 9, шаг 2):
// принимает URL, возвращает JSON с результатом проверки, поддерживает
// graceful shutdown через srv.Shutdown(shutdownCtx).
package server

// TODO: http.Server с /check handler.
// TODO: signal.NotifyContext + srv.Shutdown(shutdownCtx) с жёстким таймаутом.
// Разберись: чем Shutdown отличается от Close (Shutdown ждёт завершения
// активных соединений и останавливает listener, Close рвёт всё немедленно).
