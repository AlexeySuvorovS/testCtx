package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func checkSite(ctx context.Context, isDone chan<- struct{}, site string) {
	select {
	case <-ctx.Done():
		return
	default:
		fmt.Println("req start: ", site)
		req, err := http.NewRequestWithContext(ctx, "GET", site, nil)

		if err != nil {
			fmt.Println("request error", err.Error())
			return
		}

		t0 := time.Now()
		resp, err := http.DefaultClient.Do(req)
		defer resp.Body.Close()
		t1 := time.Now()

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("TimeOfRequest:", site, " = ", t1.Sub(t0))
		}

		isDone <- struct{}{}
	}
}

func main() {
	isDone := make(chan struct{})

	sites := []string{"http://google.com", "http://ya.ru", "http://linkedin.com"}

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	for _, url := range sites {
		go checkSite(ctx, isDone, url)
	}

	<-isDone
	cancelCtx()
}
