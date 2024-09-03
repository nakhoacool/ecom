package api

import (
	"database/sql"
	"ecom/service/product"
	"ecom/service/user"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{addr: addr, db: db}
}

func (server *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(server.db)
	userHandler := user.NewHandler(userStore)
	userHandler.UserRoutes(subrouter)

	productStore := product.NewStore(server.db)
	productHandler := product.NewHandler(productStore)
	productHandler.ProductRoutes(subrouter)

	log.Println("Listening on", server.addr)
	return http.ListenAndServe(server.addr, router)
}
