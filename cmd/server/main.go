package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"githb.com/demo-employee-api/internal/auth"
	"githb.com/demo-employee-api/internal/config"
	employee "githb.com/demo-employee-api/internal/employee"
	// "githb.com/demo-employee-api/internal/healthcheck"
	"githb.com/demo-employee-api/pkg/db"
	// gohandler "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

	// healthcheck.RegisterHandlers(router)

	auth.RegisterHandlers(
		config,
		router,
		auth.NewService(auth.NewRepository(db)),
	)

	employee.RegisterHandlers(
		config,
		router,
		employee.NewService(employee.NewRepository(db)),
	)
	handler := cors.Default().Handler(router)

	// corsHandler := gohandler.CORS(gohandler.AllowedOrigins([]string{"*"}))
	// fmt.Println(config.Server.Cors)
	srv := &http.Server{
		Addr:    ":" + fmt.Sprintf("%v", config.Server.Port),
		Handler: handler,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("cannot listen and serve")
		os.Exit(0)
	}

}
