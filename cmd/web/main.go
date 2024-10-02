package main

import (
	"log"
	"net/http"

	"github.com/MeherKandukuri/studioClasses_API/routes"
)

const portNumber = ":8080"

func main() {
	log.Println("We are starting on port number:", portNumber)
	router := routes.Routes()
	server := &http.Server{
		Addr:    portNumber,
		Handler: router,
	}

	
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
