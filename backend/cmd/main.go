package main

import (
	"context"
	"log"
	"net/http"

	"github.com/deside01/hotel_notebook/backend/internal/config"
	"github.com/deside01/hotel_notebook/backend/internal/db"
	"github.com/deside01/hotel_notebook/backend/internal/handlers"
	"github.com/deside01/hotel_notebook/backend/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5"
)

func main() {
	cfg, _ := config.Load()
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, cfg.DB_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	querier := db.New(conn)

	authHandler := handlers.NewAuthHandler(querier, cfg.SecretKey)
	guestRecordHandler := handlers.NewGuestRecordHandler(querier)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Post("/auth/registration", authHandler.Registration)
	r.Post("/auth/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware([]byte(cfg.SecretKey)))

		r.Post("/guest_records", guestRecordHandler.Create)
		r.Get("/guest_records", guestRecordHandler.GetAll)
		r.Get("/guest_records/{id}", guestRecordHandler.GetOne)
		r.Get("/guest_records/search", guestRecordHandler.Search)
		r.Patch("/guest_records/{id}", guestRecordHandler.Update)
		r.Delete("/guest_records/{id}", guestRecordHandler.Delete)
	})

	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}
