package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"githb.com/demo-employee-api/internal/config"
	employee "githb.com/demo-employee-api/internal/employee"
	"githb.com/demo-employee-api/pkg/db"
	"github.com/gorilla/mux"
)

func main() {
	config, err := config.Load()
	if err != nil {
		panic(err)
	}

	db, err := db.NewDb(config)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connection Success", db)

	router := mux.NewRouter()

	employee.RegisterHandlers(
		config,
		router,
		employee.NewService(employee.NewRepository(db)),
	)

	srv := &http.Server{
		Addr:    ":" + fmt.Sprintf("%v", config.Server.Port),
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("cannot listen and serve")
		os.Exit(0)
	}

}
