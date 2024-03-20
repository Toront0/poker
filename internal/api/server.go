package api

import (
	"github.com/go-chi/chi/v5"
	
	"github.com/Toront0/poker/internal/handlers"

	"github.com/Toront0/poker/internal/handlers/roomHandler"
	"github.com/Toront0/poker/internal/services"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
	"context"
	"net/http"
	"log"
)



type server struct {
	listenAddr string
}

func NewServer(listenAddr string) *server {
	return &server{
		listenAddr: listenAddr,
	}
}

func (s *server) Run() {
	mux := chi.NewRouter()

	c := cors.New(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	mux.Use(c.Handler)
	conn, err := pgxpool.New(context.Background(), "user=postgres.acqmtejslbbhmduhcjtw password=WpQ_4Q4_8Kw_PrY host=aws-0-eu-central-1.pooler.supabase.com port=5432 dbname=postgres")
	
	if err != nil {
		log.Fatal("could not establish connection with DB", err)
		return
	}

	authService := services.NewAuthStore(conn)
	authHandler := handlers.NewAuthHandler(authService)

	mux.Post("/sign-up", authHandler.HandleCreateAccount)
	mux.Post("/login", authHandler.HandleLoginAccount)
	mux.Get("/logout", authHandler.HandleLogout)
	mux.Get("/auth", authHandler.HandleAuthenticate)
	mux.Get("/check-email/{email}", authHandler.HandleCheckEmailExistance)
	mux.Get("/check-username/{username}", authHandler.HandleCheckUsernameExistance)
	mux.Get("/send-email-code/{email}", authHandler.HandleSendEmailCode)
	mux.Post("/verify-email-code", authHandler.HandleVerifyCode)
	mux.Get("/change-password/{password}", authHandler.HandleChangePassword)
	

	gameLobbyService := services.NewGameLobbyStore(conn)
	gameLobbyHandler := roomHandler.NewGameLobbyHandler(gameLobbyService)

	actualGameStore := services.NewActualGameStore(conn)

	mux.Get("/games", gameLobbyHandler.HandleGetAllGames)
	mux.Post("/create-game", gameLobbyHandler.HandleCreateGame)
	mux.Post("/join-game", func(w http.ResponseWriter, r *http.Request) {

		gameLobbyHandler.HandleJoinGame(w, r, mux, actualGameStore)

	})
	mux.Get("/find-game", gameLobbyHandler.HandleFindGames)
	mux.Get("/game-emojies", gameLobbyHandler.HandleGetEmojies)
	mux.HandleFunc("/ws/games", gameLobbyHandler.ServeWs)
	mux.Get("/game-results/{id}", gameLobbyHandler.HandleGetGameResults)

	mux.Get("/start-game", func (w http.ResponseWriter, r *http.Request) {

		gameLobbyHandler.HandleStartActualGame(w, r, mux, actualGameStore)

	})

	userProfileStore := services.NewUserProfileStore(conn)
	userProfileHandler := handlers.NewUserProfileHandler(userProfileStore)

	mux.Get("/imgs", userProfileHandler.HandleGetAllImages)
	mux.Post("/change-img", userProfileHandler.HandleChangeProfileImage)
	mux.Get("/user/{id}", userProfileHandler.HandleGetUserDetail)
	mux.Post("/change-username", userProfileHandler.HandleChangeUsername)
	mux.Get("/find-user", userProfileHandler.HandleFindUsers)
	mux.Get("/user-games/{id}", userProfileHandler.HandleGetUserGames)
	mux.Get("/money-status/{id}", userProfileHandler.HandleCheckPossibilityToGetMoney)
	mux.Get("/get-money/{id}", userProfileHandler.HandleGetFreeMoney)


	http.ListenAndServe(s.listenAddr, mux)

}	
