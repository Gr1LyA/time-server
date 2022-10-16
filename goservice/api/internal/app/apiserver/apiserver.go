package apiserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Cache interface {
	Store(key string, value any) error
	Load(key string) (string, error)
	Close() error
}

func Run(cache Cache) {
	wg := new(sync.WaitGroup)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan struct{})

	// goroutine for set and update data in cache
	wg.Add(1)
	go setAndUpdateTime(wg, cache, sig, done)

	handler := newHandler(cache)
	srv := http.Server{Addr: ":8080", Handler: handler}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Fprintln(os.Stderr, err)
			sig <- os.Interrupt
		}
	}()

	// After receiving signal close con and exit
	<-sig
	fmt.Println("Shutdown server...")
	if err := srv.Shutdown(context.TODO()); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	close(done)
	wg.Wait()
	if err := cache.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println("Complete!")
}

// set and update data in redis
func setAndUpdateTime(wg *sync.WaitGroup, cache Cache, sig chan os.Signal, done chan struct{}) {
	defer wg.Done()
	t := time.Now()

	err := cache.Store("responce", t.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		sig <- os.Interrupt
		return
	}

	for {
		select {
		case <-done:
			return
		default:
			if time.Since(t).Seconds() >= 10 {
				t = time.Now()
				err = cache.Store("responce", t.String())
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					sig <- os.Interrupt
					return
				}
			}
		}
	}
}