package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"gorest/models"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// go run .\cmd\api\

const version = "0.0.1"

type config struct {
	port int
	env  string
	db   struct {
		host string
	}
}

type ServerStatus struct {
	Version     string `json:"version"`
	Status      string `json:"status"`
	Environment string `json:"environment"`
}

type application struct {
	config config
	logger *log.Logger
	models models.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "environment")
	flag.StringVar(&cfg.db.host, "db.host", "postgres://postgres:Contai123@@localhost:5432/gomoviedb?sslmode=disable", "postgres database host")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := connectDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting API server on port %d in %s mode\n", cfg.port, cfg.env)

	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func connectDB(cfg config) (*sql.DB, error) {
	conn, err := sql.Open("postgres", cfg.db.host)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = conn.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	// var name [5]string
	// err = conn.QueryRow(ctx, "SELECT email FROM user_entity").Scan(&name)
	// if err != nil {
	// 	return nil, err
	// }

	// log.Printf("Connected to database %s\n", name)

	return conn, nil
}
