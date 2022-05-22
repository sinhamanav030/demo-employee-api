package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"githb.com/demo-employee-api/internal/config"
	employee "githb.com/demo-employee-api/internal/employee"

	"githb.com/demo-employee-api/pkg/db"
	"githb.com/demo-employee-api/pkg/token"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate)
	logger.SetFlags(log.LstdFlags | log.Llongfile)

	config, err := config.Load(logger)
	if err != nil {
		logger.Fatal(err)
	}

	db, err := db.NewDb(config, logger)
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println("Connection Success", db)

	router := mux.NewRouter()

	tokenMaker, err := token.NewJwtMaker(config.Auth.JwtKey)
	if err != nil {
		logger.Fatal(err)
	}

	employee.RegisterHandlers(
		config,
		router,
		employee.NewService(employee.NewRepository(db, logger), logger, tokenMaker),
		logger,
		tokenMaker,
	)
	handler := cors.Default().Handler(router)

	srv := &http.Server{
		Addr:    ":" + fmt.Sprintf("%v", config.Server.Port),
		Handler: handler,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal("cannot listen and serve")
		// os.Exit(0)
	}

}
