package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/namsral/flag"
)

const (
	envHost = "DB_HOST"
	envUser = "DB_USER"
	envPass = "DB_PASS"
)

type service struct {
	router *mux.Router
	db     FeedDB
}

func main() {
	log.Println("service started")

	var host, user, pass string

	flag.StringVar(&host, envHost, "", "host:port for postgres")
	flag.StringVar(&user, envUser, "", "username for postgres")
	flag.StringVar(&pass, envPass, "", "password for postgres")

	flag.Parse()

	if len(user) == 0 {
		log.Fatal("required auth credentials not set. quitting.")
	}

	dbInstance, err := InitPostgresDB(host, user, pass)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to database.")

	srv := service{
		router: mux.NewRouter(),
		db:     dbInstance,
	}

	srv.start()
}

func (srv service) start() {
	srv.router.HandleFunc("/category", srv.categoryHandler)
	srv.router.HandleFunc("/news", srv.newsHandler)
	srv.router.HandleFunc("/news/{category:[0-9]+}", srv.newsByCategoryHandler)

	server := &http.Server{Addr: ":8880", Handler: srv.router}

	go server.ListenAndServe()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("unable to stop gracefully server: %v", err)
	}

	srv.db.close()

	log.Println("service stopped.")
}
