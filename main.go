package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
}

func (s *Server) OnProbe(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{})
}

func (s *Server) OnWork(ctx *gin.Context) {
	log.Printf("[%s] start work!\n", time.Now().String())
	time.Sleep(300 * time.Millisecond)
	buf, _ := os.Hostname()
	ctx.JSON(http.StatusOK, gin.H{
		"hello": buf,
	})
	log.Printf("[%s] fin work!\n", time.Now().String())
}

func main() {
	wg := &sync.WaitGroup{}
	s := &Server{}
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.GET("/", s.OnProbe)
	r.GET("/work", s.OnWork)

	srv := &http.Server{
		Addr:    ":9090",
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	wg.Add(1)
	go func() {
		defer wg.Done()
		sig := <-quit
		log.Printf("[%s] load signal: %s, will be shutdown server!", time.Now().String(), sig.String())

		time.Sleep(15 * time.Second)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("Server forced to shutdown: ", err)
		}

		log.Println(time.Now().String(), ": ", "FIN shutdown server!!")
	}()

	fmt.Println("run server!!")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	} else {
		log.Println(err)
	}

	wg.Wait()
}
