package main

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	pb "github.com/Zach-Johnson/tempus/proto/api/v1/tempus"
	storage "github.com/Zach-Johnson/tempus/server/db"
	"github.com/Zach-Johnson/tempus/server/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

//go:embed frontend/dist/*
var embeddedFiles embed.FS

var (
	dbPath    string
	grpcPort  int
	httpPort  int
	enableTLS bool
)

func main() {
	flag.BoolVar(&enableTLS, "tls", false, "Enable TLS for gRPC server")
	flag.IntVar(&httpPort, "http-port", 8080, "HTTP server port")
	flag.IntVar(&grpcPort, "grpc-port", 9090, "gRPC server port")
	flag.StringVar(&dbPath, "db-path", "./data/tempus.db", "Path to SQLite database file")
	flag.Parse()

	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting tempus application")

	dbPathEnv, ok := os.LookupEnv("DB_PATH")
	if ok {
		dbPath = dbPathEnv
	}

	// Initialize database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Ping database to verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Set pragmas for SQLite
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatalf("Failed to set foreign_keys pragma: %v", err)
	}
	if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		log.Fatalf("Failed to set journal_mode pragma: %v", err)
	}

	// Initialize the storage layer
	store := storage.NewSQLiteStore(db)

	if os.Getenv("RUN_MIGRATIONS") == "true" {
		if err := storage.RunMigrations(store.GetDB()); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("Migrations ran successfully")
	}

	// Start the gRPC server
	go startGRPCServer(store)

	// Start the HTTP server (gRPC-Gateway)
	go startHTTPServer()

	// Wait for interrupt signal
	waitForShutdown()
}

func startGRPCServer(store *storage.SQLiteStore) {
	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	categoryService := handlers.NewCategoryHandler(store.GetDB())
	tagService := handlers.NewTagService(store.GetDB())
	exerciseService := handlers.NewExerciseHandler(store.GetDB())
	practiceSessionService := handlers.NewPracticeSessionHandler(store.GetDB())
	exerciseHistoryService := handlers.NewExerciseHistoryHandler(store.GetDB())

	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	pb.RegisterTagServiceServer(grpcServer, tagService)
	pb.RegisterExerciseServiceServer(grpcServer, exerciseService)
	pb.RegisterPracticeSessionServiceServer(grpcServer, practiceSessionService)
	pb.RegisterExerciseHistoryServiceServer(grpcServer, exerciseHistoryService)

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	// Start listening
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("gRPC server listening on :%d", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func startHTTPServer() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create a client connection to the gRPC server
	conn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%d", grpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to dial gRPC server: %v", err)
	}
	defer conn.Close()

	// Register gRPC-Gateway
	gwmux := runtime.NewServeMux()
	if err := pb.RegisterCategoryServiceHandler(ctx, gwmux, conn); err != nil {
		log.Fatalf("Failed to register gateway for CategoryService: %v", err)
	}
	if err := pb.RegisterTagServiceHandler(ctx, gwmux, conn); err != nil {
		log.Fatalf("Failed to register gateway for TagService: %v", err)
	}
	if err := pb.RegisterExerciseServiceHandler(ctx, gwmux, conn); err != nil {
		log.Fatalf("Failed to register gateway for ExerciseService: %v", err)
	}
	if err := pb.RegisterPracticeSessionServiceHandler(ctx, gwmux, conn); err != nil {
		log.Fatalf("Failed to register gateway for PracticeSessionService: %v", err)
	}
	if err := pb.RegisterExerciseHistoryServiceHandler(ctx, gwmux, conn); err != nil {
		log.Fatalf("Failed to register gateway for ExerciseHistoryService: %v", err)
	}

	staticFS, err := fs.Sub(embeddedFiles, "frontend/dist")
	if err != nil {
		log.Fatalf("failed to open file system")
	}
	mux := http.NewServeMux()

	// API routes
	mux.Handle("/api/", http.StripPrefix("/api", middleware(gwmux)))

	// Frontend (embedded)
	mux.Handle("/", spaHandlerFS(staticFS, "index.html"))

	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: basicAuthMiddleware(mux),
	}

	// Start HTTP server
	log.Printf("HTTP server listening on :%d", httpPort)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}

func spaHandlerFS(staticFS fs.FS, indexPath string) http.Handler {
	fileServer := http.FileServer(http.FS(staticFS))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath := r.URL.Path
		filePath := path.Clean(requestedPath[1:]) // remove leading slash

		f, err := staticFS.Open(filePath)
		if err == nil {
			defer f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}

		// Fallback to index.html
		indexFile, err := staticFS.Open(indexPath)
		if err != nil {
			http.Error(w, "index.html not found", http.StatusInternalServerError)
			return
		}
		defer indexFile.Close()

		// Read the file into memory
		content, err := io.ReadAll(indexFile)
		if err != nil {
			http.Error(w, "failed to read index.html", http.StatusInternalServerError)
			return
		}

		// Use a bytes.Reader to satisfy io.ReadSeeker
		http.ServeContent(w, r, indexPath, time.Now(), bytes.NewReader(content))
	})
}

type statusRec struct {
	http.ResponseWriter
	status int
}

func (r *statusRec) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// Add middleware for REST API
func middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		start := time.Now()
		wrapped := &statusRec{ResponseWriter: w, status: http.StatusOK}

		handler.ServeHTTP(wrapped, r)

		if r.URL.Path != "/healthz" {
			log.Println("request",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start),
				"status", wrapped.status,
			)
		}
	})
}

func basicAuthMiddleware(next http.Handler) http.Handler {
	if os.Getenv("ENV") != "prod" {
		// Disable auth in non-production
		return next
	}

	user := os.Getenv("BASIC_AUTH_USER")
	pass := os.Getenv("BASIC_AUTH_PASS")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || u != user || p != pass {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func waitForShutdown() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	log.Printf("Received signal %v, shutting down...", sig)

	time.Sleep(time.Second) // Give time for graceful shutdown
	log.Println("Server shutdown complete")
}
