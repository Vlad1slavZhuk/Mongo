package server

import (
	"Mongo/internal/pkg/config"
	"Mongo/internal/pkg/graphql/graph"
	"Mongo/internal/pkg/graphql/graph/generated"
	"Mongo/internal/pkg/grpc"
	"Mongo/internal/pkg/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
)

type Server struct {
	router  http.Handler
	config  *config.Config
	service service.InterfaceServer
}

func NewServer() *Server {
	return new(Server)
}

func (s *Server) SetConfig() {
	s.config = config.NewConfigFromEnv()
}

func (s *Server) SetHandlers() {
	s.router = mux.NewRouter()
	s.SetSubRouter()
}

func (s *Server) SetSubRouter() {
	subR, ok := s.router.(*mux.Router)
	if !ok {
		log.Fatal("Expected mux subRouter.")
	}
	subR.HandleFunc("/login", s.Login).Methods(http.MethodPost)                // Login.
	subR.HandleFunc("/signup", s.SignUp).Methods(http.MethodPost)              // Sign Up.
	subR.HandleFunc("/logout", s.Logout).Methods(http.MethodPost)              // Logout.
	subR.HandleFunc("/ads", s.GetAll).Methods(http.MethodGet)                  // Get Ads.
	subR.HandleFunc("/ad", s.Create).Methods(http.MethodPost)                  // Create Ad.
	subR.HandleFunc("/ad/{id:[1-9]\\d*}", s.Get).Methods(http.MethodGet)       // Get Ad.
	subR.HandleFunc("/ad/{id:[1-9]\\d*}", s.Update).Methods(http.MethodPut)    // Update Ad.
	subR.HandleFunc("/ad/{id:[1-9]\\d*}", s.Delete).Methods(http.MethodDelete) // Delete Ad and update ID other Ad.
	gqlSrv := subR.Handle("/query", handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{IGraphqL: s.service}})))

	url, err := gqlSrv.URL()
	if err != nil {
		log.Println("Not Work gql")
	}
	subR.Handle("/gql", playground.Handler("GQL", url.Path))
}

func (s *Server) SetStorage() {
	service := new(service.Service)
	service.SetStorage(grpc.NewGrpcClient())
	s.service = service
}

func (s *Server) Run() {
	server := &http.Server{
		Addr:         ":8000",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		Handler:      s.router,
	}

	//log.Println("HTTP server launching on port " + s.config.Port)
	//log.Fatal(http.ListenAndServe(s.config.Port, s.router))
	go func() {
		log.Fatal(http.ListenAndServe(s.config.PortHTTP, s.router))
	}()
	log.Println("Start server on localhost:", s.config.PortHTTP)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)
	defer close(c)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}
	log.Printf("HTTP server is shutdown...")
}
