package server

import (
	"log"
	"net/http"
	"time"
	"tinder_users/internal/server/routes"

	"github.com/go-chi/chi"
)

//Devolvemos un puntero con nuestro server
type Server struct {
	server *http.Server
}

//Inicializamos el servidor y montamos los endpoints
func New(port string) (*Server, error) {
	//Estructura que funciona de mux
	r := chi.NewRouter()

	//Se monta como raiz la direccion "api"
	r.Mount("/api/tinder", routes.New())

	serv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	//Construimos un server inicializado con el que acabamos de crear
	server := Server{server: serv}
	return &server, nil
}

func (serv *Server) Start() {
	log.Printf("Servidor corriendo")
	log.Fatal(serv.server.ListenAndServe())
}
