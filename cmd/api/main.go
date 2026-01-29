package main

import (
	"lesson-proj/internal/database"
	"lesson-proj/internal/handlers"
	authService "lesson-proj/internal/services/auth" 
	productService "lesson-proj/internal/services/products"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	serverPort := os.Getenv("SERVER_PORT")

	log.Printf("Starting server on port %s with database %s", serverPort, databaseURL)
	log.Printf("DATABASE_URL from env: %q", databaseURL)
	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// close the database connection when main function exits
	defer db.Close()

	log.Println("Connected to the database successfully")

	productRepository := database.NewProductRepository(db)
	productService := productService.NewProductService(productRepository)
	handler := handlers.NewProductHandler(productService)

	userRepository := database.NewUserRepository(db)
	userService := authService.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	router := http.NewServeMux()
	router.HandleFunc("/products", methodHandler(handler.GetAllProducts, http.MethodGet))
	router.HandleFunc("/products/create", methodHandler(handler.CreateProduct, http.MethodPost))
	router.HandleFunc("/products/", productIDHandler(handler))

	router.HandleFunc("/users", methodHandler(userHandler.GetAllUsers, http.MethodGet))
	router.HandleFunc("/users/create", methodHandler(userHandler.Registration, http.MethodPost))
	router.HandleFunc("/users/", userIDHandler(userHandler))
	router.HandleFunc("/users/auth", methodHandler(userHandler.Authorization, http.MethodPost))

	loggedRouter := loggingMiddleware(router)
	corsHandler := corsMiddleware(loggedRouter)

	srv := &http.Server{
		Addr:    ":" + serverPort,
		Handler: corsHandler,

		// Production-friendly timeouts:
		ReadHeaderTimeout: 5 * time.Second,  // time to read request headers
		ReadTimeout:       10 * time.Second, // time to read the entire request
		WriteTimeout:      20 * time.Second, // time to write the response
		IdleTimeout:       60 * time.Second, // time to wait for the next request
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
