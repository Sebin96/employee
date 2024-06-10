package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type CustomRouter struct {
	*mux.Router
	DB *sql.DB `json:"db,omitempty"`
}

// NewCustomRouter creates a new CustomRouter instance with the provided router and database connection
func NewCustomRouter(router *mux.Router, db *sql.DB) *CustomRouter {
	return &CustomRouter{
		Router: router,
		DB:     db,
	}
}

func main() {

	db := initDB()
	defer db.Close()

	// Create a new router
	r := mux.NewRouter()

	// Setup routes
	customRouter := NewCustomRouter(r, db)

	// Setup routes
	customRouter.SetupRouter()

	// Start the HTTP server
	port := ":8080"
	fmt.Printf("Server started at %s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
