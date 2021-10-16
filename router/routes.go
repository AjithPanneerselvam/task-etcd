package router

import (
	"fmt"
	"time"

	"github.com/AjithPanneerselvam/task-etcd/auth"
	"github.com/AjithPanneerselvam/task-etcd/client/github"
	"github.com/AjithPanneerselvam/task-etcd/config"
	"github.com/AjithPanneerselvam/task-etcd/handler/login"
	"github.com/AjithPanneerselvam/task-etcd/handler/task"
	"github.com/AjithPanneerselvam/task-etcd/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	GithubCallbackURLFormat       = "http://%s:%s/login/github/callback"
	LoginSuccessRedirectURLFormat = "http://%s:%s/home"
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
	githubCallbackURL := fmt.Sprintf(GithubCallbackURLFormat, config.HostName, config.ListenPort)
	loginSuccessRedirectURL := fmt.Sprintf(LoginSuccessRedirectURLFormat, config.HostName, config.ListenPort)

	githubClient := github.New(config.GithubOAuthURL, config.GithubAPIURL, config.GithubClientID,
		config.GithubClientSecret, config.GithubTimeoutInSec)

	jwtAuthenticator := auth.NewJWTAuth(config.JWTSecretyKey, time.Minute*time.Duration(config.JWTExpiryInMins))

	githubLoginHandler := login.NewGithubLoginHandler(githubClient, githubCallbackURL,
		jwtAuthenticator, loginSuccessRedirectURL)
	taskHandler := task.NewTaskHandler()

	r.Use(middleware.Logger)

	r.Get("/", githubLoginHandler.Home)

	// login routes
	r.Route("/login", func(r chi.Router) {
		r.Get("/github", githubLoginHandler.Login)
		r.Get("/github/callback", githubLoginHandler.Callback)
	})

	// task routes
	r.Group(func(r chi.Router) {
		r.Use(jwtAuthenticator.Authenticator)

		r.Route("/task", func(r chi.Router) {
			r.Get("/create", taskHandler.CreateTask)
		})
	})
}
