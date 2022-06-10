package server

import (
	"fmt"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/server/controllers"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func Run() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	router := controllers.InitRoutes()
	log.Fatal(http.ListenAndServe(":9090", router))
}
