// 模拟两种退出方式：ctrl+c 信号退出，某个gr报错退出；
package main

import (
	"context"
	"errors"
	"fmt"
	"html"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

var myHandler = func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "w, Hello, %q", html.EscapeString(r.URL.Path))

	return
}

func closeServers(servers map[string]*http.Server) {
	for serverName, server := range servers {
		if err := server.Shutdown(context.TODO()); err != nil {
			// TODO: 如果shutdown失败呢？
			fmt.Println("shutdown fail", serverName)
		}
	}
}

func main() {

	// NOTE(shenglong): 配置两个httpServer
	mux1 := http.NewServeMux()
	mux1.HandleFunc("/foo", myHandler)
	httpServer1 := http.Server{
		Addr:    "0.0.0.0:8811",
		Handler: mux1,
	}

	mux2 := http.NewServeMux()
	mux2.HandleFunc("/bar", myHandler)
	httpServer2 := http.Server{
		Addr:    "0.0.0.0:8812",
		Handler: mux2,
	}

	// NOTE(shenglong): httpServer注册到serverMap
	serverMap := make(map[string]*http.Server)
	serverMap["server1"] = &httpServer1
	serverMap["server2"] = &httpServer2

	// NOTE(shenglong): errgroup start gr
	g, ctx := errgroup.WithContext(context.Background())

	// server 1
	g.Go(func() error {
		defer fmt.Println("gr1: exit")
		return httpServer1.ListenAndServe()
	})

	// server 2
	g.Go(func() error {
		defer fmt.Println("gr2: exit")
		return httpServer2.ListenAndServe()
	})

	// 模拟gr 报错退出
	g.Go(func() error {
		// mock http server return error
		defer fmt.Println("grErr: exit")

		var exitString string
		select {
		case <-time.After(10 * time.Second):
			exitString = "grErr 主动退出(timeout)"
		case <-ctx.Done():
			exitString = "grErr 被动退出"
		}
		fmt.Println(exitString)
		return errors.New(exitString)
	})

	// 监听errgroup cancel，通知其他gr
	g.Go(func() error {
		defer fmt.Println("grMonitor: exit")
		select {
		case <-ctx.Done():
			closeServers(serverMap)
		}
		return nil
	})

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalCh:
		// 监听信号，通知所有gr退出
		fmt.Println("main: 收到信号，开始退出")
		closeServers(serverMap)

	case <-ctx.Done():
		// 监听errgroup cancel
		fmt.Println("main", ctx.Err())
	}

	if err := g.Wait(); err != nil {
		fmt.Println("main: exit reason:", err.Error())
	}
}
