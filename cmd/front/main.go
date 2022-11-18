package main

import (
	"context"
	"cryptoshare/cmd/front/handler"
	_ "cryptoshare/conf"
	"cryptoshare/ds"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// to get file line and path when print
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	port := os.Getenv("APP_FRONT_PORT")
	addr := fmt.Sprintf(":%s", port)

	// load datasource
	ds, err := ds.NewDataSource()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	h := handler.NewHandler(
		&handler.HConfig{
			R:  router,
			DS: ds,
		})

	h.Register()

	server := http.Server{
		Addr:           addr,
		Handler:        h.R,
		ReadTimeout:    time.Duration(time.Minute * 3),
		WriteTimeout:   time.Duration(time.Minute * 3),
		MaxHeaderBytes: 10 << 20, //10MB
	}

	go func() {
		log.Println("server started listening on port : ", addr)
		if err := server.ListenAndServe(); err != nil {
			log.Println("server failed to initialized  on port : ", addr)
			log.Fatalf("error on listening :%v\n", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c

	// shutdown close
	if err := server.Shutdown(context.Background()); err != nil {
		log.Println("Failed to shutdown server: ", err.Error())
	}
}
