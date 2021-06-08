package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/n404an/gomv/service"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	if err := service.Start(ctx); err != nil {
		log.Printf("failed to start gomv: %+v\n", err)
	}
}
