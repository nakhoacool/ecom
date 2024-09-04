package api

import (
	"database/sql"
	"ecom/middleware"
	"ecom/service/auth"
	"ecom/service/cart"
	"ecom/service/order"
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

	authStore := auth.NewStore()
	userStore := user.NewStore(server.db)
	userHandler := user.NewHandler(userStore, authStore)
	userHandler.UserRoutes(subrouter)

	productStore := product.NewStore(server.db)
	productHandler := product.NewHandler(productStore)
	productHandler.ProductRoutes(subrouter)

	orderStore := order.NewStore(server.db)

	cartHandler := cart.NewHandler(orderStore, productStore)
	cartSubrouter := subrouter.PathPrefix("/cart").Subrouter()
	cartSubrouter.Use(middleware.JWTMiddleware)
	cartHandler.RegisterRoutes(cartSubrouter)

	log.Println("Listening on", server.addr)
	return http.ListenAndServe(server.addr, router)
}
