package main

import (
	"fmt"
	"log"
	"net/http"

	"hash-service/backend/internal/handler"
	"hash-service/backend/internal/service"
)

func main() {
	svc := service.NewHashService()
	h := handler.NewHashHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/hash", h.HandleHash)
	mux.HandleFunc("/api/health", h.HandleHealth)

	srv := handler.CORSMiddleware(mux)

	addr := ":8080"
	fmt.Printf("Server starting on %s\n", addr)
	if err := http.ListenAndServe(addr, srv); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
