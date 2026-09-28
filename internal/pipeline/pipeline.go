// Package pipeline перестраивает обработку ссылок в конвейер стадий
// (Этап 10): gen -> normalize -> dedup -> check -> collect.
// Каждая стадия — функция вида func(ctx, <-chan T) <-chan U и сама создаёт
// и закрывает свой выходной канал (см. Go Blog: "Pipelines and cancellation").
package pipeline

// TODO: func Gen(ctx context.Context, urls []string) <-chan string
// TODO: func Normalize(ctx context.Context, in <-chan string) <-chan string
// TODO: func Dedup(ctx context.Context, in <-chan string) <-chan string
// TODO: func Check(ctx context.Context, in <-chan string, opts checker.Options) <-chan checker.Result
// TODO: func Collect(ctx context.Context, in <-chan checker.Result) []checker.Result
