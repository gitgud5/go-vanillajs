package main

// This is very tricky but very important
import (
	"database/sql"
	"fmt"
	"go-vanillajs/data"
	"go-vanillajs/handlers"
	"go-vanillajs/logger"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

func initializeLogger() *logger.Logger {
	Loginstance, err := logger.NewLogger("movies.log")
	if err != nil {
		log.Fatalf("failed to initialize logger %v", err)
	}

	// but why this line below of defer?
	defer Loginstance.Close()

	return Loginstance

}

func main() {

	// Log initializer
	logInstance := initializeLogger()

	// Environmental variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file was available")
	}

	// Connect to the DB
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	fmt.Println("Hello world!")

	// Initialize repositories
	movieRepo, err := data.NewMovieRepository(db, logInstance)
	if err != nil {
		log.Fatalf("Failed to initialize movie repository: %v", err)
	}

	// Initialize handlers
	movieHandler := handlers.NewMovieHandler(movieRepo, logInstance)
	// authHandler := handlers.NewAuthHandler(userStorage, jwt, logInstance)

	// Set up routes
	// I wrote this for myself to help me understand what is happening in the background
	http.Handle("/api/movies/random", http.HandlerFunc(movieHandler.GetRandomMovies))
	http.HandleFunc("/api/movies/top", movieHandler.GetTopMovies)
	http.HandleFunc("/api/movies/search", movieHandler.SearchMovies)
	http.HandleFunc("/api/movies/", movieHandler.GetMovie)
	http.HandleFunc("/api/genres", movieHandler.GetGenres)
	http.HandleFunc("/api/account/register", movieHandler.GetGenres)
	http.HandleFunc("/api/account/authenticate", movieHandler.GetGenres)

	// This handler is for static files
	http.Handle("/", http.FileServer(http.Dir("public")))

	logInstance.Info("starting server...")

	const addr string = ":8080"
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("not working error happened")
		logInstance.Error("Server failed", err)
	}

}
