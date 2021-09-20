package loader

type Opt func(w *Worker)

func WithGoroutine(c uint) Opt {
	return func(w *Worker) {
		w.goroutine = c
	}
}
