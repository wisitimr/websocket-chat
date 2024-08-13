package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"websocket-chat/internal/config"
	_handler "websocket-chat/internal/handler"
	"websocket-chat/internal/model"
	mHandler "websocket-chat/internal/model/handler"
	mRepo "websocket-chat/internal/model/repository"
	mService "websocket-chat/internal/model/service"
	_repo "websocket-chat/internal/repository"
	_service "websocket-chat/internal/service"
	"websocket-chat/pkg/redis"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
)

type Server struct {
	cfg    config.HTTPServer
	router *chi.Mux
	logger *logrus.Logger
}

var clients = make(map[*model.Client]bool)
var broadcast = make(chan *model.Chat)

func NewHttpServer(ctx context.Context, cfg config.HTTPServer) (*Server, error) {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05.999999999",
		FullTimestamp:   true,
	})
	s := Server{
		logger: logger,
		cfg:    cfg,
		router: register(ctx, logger),
	}
	return &s, nil
}

func register(ctx context.Context, logger *logrus.Logger) *chi.Mux {
	// mongodb, err := db.Connect(ctx)
	// if err != nil {
	// 	log.Fatalf("Can't establish database connection")
	// }
	redisClient := redis.InitializeRedis()
	// defer redisClient.Close()

	r := chi.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(c.Handler)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	// init service
	repo := mRepo.Repository{
		Redis: _repo.InitRedisRepository(redisClient, logger),
	}

	// init service
	service := mService.Service{
		WebSocket:     _service.InitWsService(repo, logger, broadcast),
		Register:      _service.InitRegisterService(repo, logger),
		Login:         _service.InitLoginService(repo, logger),
		VerifyContact: _service.InitVerifyContactService(repo, logger),
		ChatHistory:   _service.InitChatHistoryService(repo, logger),
		ContactList:   _service.InitContactListService(repo, logger),
	}

	// init handler
	handler := mHandler.Handler{
		// AuthToken: jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil),
		WebSocket:     _handler.InitWsHandler(service, logger),
		Register:      _handler.InitRegisterHandler(service, logger),
		Login:         _handler.InitLoginHandler(service, logger),
		VerifyContact: _handler.InitVerifyContactHandler(service, logger),
		ChatHistory:   _handler.InitChatHistoryHandler(service, logger),
		ContactList:   _handler.InitContactListHandler(service, logger),
	}
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Chat Room!"))
	})
	r.Post("/register", handler.Register.NewUser)
	r.Post("/login", handler.Login.Login)
	r.Post("/verify-contact", handler.VerifyContact.VerifyContact)
	r.Get("/chat-history", handler.ChatHistory.ChatHistory)
	r.Get("/contact-list", handler.ContactList.ContactList)

	go broadcaster()

	r.HandleFunc("/ws", handler.WebSocket.WebSocketConnect)

	return r
}

func (s Server) Start(ctx context.Context) error {
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		Handler:      s.router,
		IdleTimeout:  s.cfg.IdleTimeout,
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
	}

	stopServer := make(chan os.Signal, 1)
	signal.Notify(stopServer, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(stopServer)

	// channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		s.logger.Printf("Server started on port %d", s.cfg.Port)
		serverErrors <- server.ListenAndServe()
	}(&wg)

	// blocking run and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("error: starting rest api server: %w", err)
	case <-stopServer:
		s.logger.Warn("server received stop signal")
		// asking listener to shutdown
		err := server.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("graceful shutdown did not complete: %w", err)
		}
		wg.Wait()
		s.logger.Info("server was shut down gracefully")
	}
	return nil
}

func broadcaster() {
	for {
		message := <-broadcast
		// send to every client that is currently connected
		fmt.Println("new message", message)

		for client := range clients {
			// send message only to involved users
			fmt.Println("username:", client.Username,
				"from:", message.From,
				"to:", message.To)

			if client.Username == message.From || client.Username == message.To {
				err := client.Conn.WriteJSON(message)
				if err != nil {
					log.Printf("Websocket error: %s", err)
					client.Conn.Close()
					delete(clients, client)
				}
			}
		}
	}
}
