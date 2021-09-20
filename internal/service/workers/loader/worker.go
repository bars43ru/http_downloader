package loader

import (
	"log"
	"net/http"
	"sync"

	"downloader_test/internal/service/entity"
)

type Worker struct {
	goroutine uint
	client    *http.Client
}

func New(opts ...Opt) *Worker {
	tr := &http.Transport{}
	r := &Worker{
		goroutine: countGoroutine,
		client:    &http.Client{Transport: tr},
	}
	for _, f := range opts {
		f(r)
	}
	log.Println("//---------------------//")
	log.Printf("worker.goroutine = %d\n", r.goroutine)
	log.Println("//---------------------//")
	return r
}

func (w Worker) Run(in <-chan string, out chan<- entity.Response) {
	limiter := make(chan struct{}, w.goroutine)
	wg := &sync.WaitGroup{}
	for v := range in {
		limiter <- struct{}{}
		wg.Add(1)
		v := v
		go func() {
			defer wg.Done()
			defer func() { <-limiter }()
			ret := entity.Response{Url: v}
			r, err := w.client.Get(v)
			if err != nil {
				ret.Err = err
			} else {
				ret.StatusCode = r.StatusCode
				ret.Status = r.Status
			}
			out <- ret
		}()
	}
	wg.Wait()
	close(out)
}
