package main

import (
	"encoding/json"
	"log"
	"net/http"

	"go-workflow/config"
	"go-workflow/controllers"
	"go-workflow/services"
	"go-workflow/websocket"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func main() {
	// Database bağlantısı
	db := config.ConnectDB()

	// WebSocket hub'ını oluştur
	hub := websocket.NewHub()
	go hub.Run()

	// Servisleri oluştur
	workflowService := services.NewWorkflowService(db.DB, hub)

	// Controller'ları oluştur
	workflowController := controllers.NewWorkflowController(workflowService)

	// Router'ı oluştur
	r := chi.NewRouter()

	// Middleware'leri ekle
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	// Ana route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response := Response{Message: "Hoş geldiniz!"}
		json.NewEncoder(w).Encode(response)
	})

	// WebSocket endpoint'i
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// WebSocket bağlantısını yükselt
		conn, err := websocket.Upgrade(w, r)
		if err != nil {
			log.Printf("WebSocket yükseltme hatası: %v", err)
			return
		}

		// Yeni client oluştur
		client := &websocket.Client{
			Hub:  hub,
			Conn: conn,
			Send: make(chan []byte, 256),
		}

		// Client'ı hub'a kaydet
		hub.Register <- client

		// Client'ın mesajlarını dinle
		go client.WritePump()
		go client.ReadPump()
	})

	// Workflow route'larını kaydet
	workflowController.RegisterRoutes(r)

	// Sunucuyu başlat
	port := ":8080"
	println("Sunucu " + port + " portunda başlatılıyor...")
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}
