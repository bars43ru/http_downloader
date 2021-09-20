package main

import (
	"bufio"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"downloader_test/internal/service/entity"
	"downloader_test/internal/service/gateway/openapi"
	"downloader_test/internal/service/workers/loader"
	"downloader_test/internal/service/workers/orchestrator"
)

func main() {
	var wg sync.WaitGroup
	ctx := gracefulWorker(context.Background())

	p := orchestrator.New(ctx)
	task := make(chan string)
	tr := make(chan entity.Response)

	runWorker(&wg, func() { httpServerRun(ctx, p) }, "http server")
	runWorker(&wg, func() { p.Sub(task) }, "orchestrator")
	runWorker(&wg, func() { newWorker().Run(task, tr) }, "worker")
	runWorker(&wg, func() { output(tr) }, "output")
	runWorker(&wg, func() { initLoad(ctx, p) }, "read from file")

	<-ctx.Done()
	log.Println("waiting for everyone to finish worker")
	wg.Wait()
	log.Println("application exit")
}

func runWorker(wg *sync.WaitGroup, worker func(), caption string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("[starting]: %s\n", caption)
		worker()
		log.Printf("[stoped]: %s\n", caption)
	}()
}

func httpServerRun(ctx context.Context, p *orchestrator.PubSub) {
	e := echo.New()
	e.Use(
		middleware.Recover(),
		middleware.Logger())
	api := e.Group("/api")

	openapi.New(p).Configure(api)

	go func() {
		<-ctx.Done()
		log.Println("shutdown http server")
		if err := e.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()

	if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
		log.Fatal("shutting down the server")
	}
}

func newWorker() *loader.Worker {
	if goroutine, err := strconv.Atoi(os.Getenv("REQUEST_LIMIT")); err != nil {
		loader.New(loader.WithGoroutine(uint(goroutine)))
	}
	return loader.New(loader.WithGoroutine(10))
}

func gracefulWorker(ctx context.Context) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		oscall := <-c
		log.Printf("system call: %+v", oscall)
		cancel()
	}()
	return ctx
}

func initLoad(ctx context.Context, p *orchestrator.PubSub) {
	f := os.Getenv("FILE")
	if f == "" {
		log.Println("system Environment `FILE` not set")
		return
	}

	file, err := os.Open(f)
	if err != nil {
		log.Printf("open file: %v\n",err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			break
		default:
			p.Pub(scanner.Text())
		}
	}
}

func output(in <-chan entity.Response) {
	for v := range in {
		log.Printf("%s\n", v)
	}
}
