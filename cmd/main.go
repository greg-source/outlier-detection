package main

import (
	"log"
	"os"
	"os/signal"
	"outlier-detection/server"
)

func main() {
	server := server.New()
	port := "8080"
	go func() {
		if err := server.Run(port); err != nil {
			log.Fatalf("error occured while running http server: %s\n", err.Error())
		}
	}()
	log.Printf("application is up and running on port: %s", port)

	systemQuit := make(chan os.Signal, 1)

	signal.Notify(systemQuit, os.Interrupt)

	<-systemQuit
	log.Println("shutdown initialized")
	server.Shutdown()
}
