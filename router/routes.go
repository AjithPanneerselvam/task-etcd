package router

import (
	"fmt"
	"net/http"

	"github.com/AjithPanneerselvam/todo/client/github"
	"github.com/AjithPanneerselvam/todo/config"
	"github.com/AjithPanneerselvam/todo/handler/login"
	"github.com/AjithPanneerselvam/todo/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	GithubCallbackURLFormat = "http://%s:%s/login/github/callback"
)

type Router struct {
	*chi.Mux
}

func NewRouter() *Router {
	return &Router{
		Mux: chi.NewRouter(),
	}
}

func (r *Router) AddRoutes(config *config.Config, userStore store.UserStore) {
	githubClient := github.New(config.GithubOAuthURL, config.GithubAPIURL, config.GithubClientID,
		config.GithubClientSecret, config.GithubTimeoutInSec)

	githubCallbackURL := fmt.Sprintf(GithubCallbackURLFormat, config.HostName, config.ListenPort)

	githubLoginHandler := login.NewGithubLoginHandler(githubClient, githubCallbackURL, userStore)

	r.Use(middleware.Logger)

	r.Get("/", githubLoginHandler.HomePage)

	r.Route("/login", func(r chi.Router) {
		r.Get("/github", githubLoginHandler.Login)
		r.Get("/github/callback", githubLoginHandler.Callback)
	})

	r.Handle("/home", http.FileServer(http.Dir("./static")))

	/*r.Route("/user", func(r chi.Router) {
		r.Get("/{userId}", userHandler.GetInfo)
	})*/
}
