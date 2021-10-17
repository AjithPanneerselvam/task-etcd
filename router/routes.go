package router

import (
	"fmt"
	"net/http"
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

func (r *Router) AddRoutes(config *config.Config, taskStore store.TaskStore) {
	githubCallbackURL := fmt.Sprintf(GithubCallbackURLFormat, config.HostName, config.ListenPort)
	loginSuccessRedirectURL := fmt.Sprintf(LoginSuccessRedirectURLFormat, config.HostName, config.ListenPort)

	githubClient := github.New(config.GithubOAuthURL, config.GithubAPIURL, config.GithubClientID,
		config.GithubClientSecret, config.GithubTimeoutInSec)

	jwtAuthenticator := auth.NewJWTAuth(config.JWTSecretyKey, time.Minute*time.Duration(config.JWTExpiryInMins))

	githubLoginHandler := login.NewGithubLoginHandler(githubClient, githubCallbackURL,
		jwtAuthenticator, loginSuccessRedirectURL)
	taskHandler := task.NewTaskHandler(taskStore)

	r.Use(middleware.Logger)

	r.Get("/", githubLoginHandler.Home)

	// login routes
	r.Route("/login", func(r chi.Router) {
		r.Get("/github", githubLoginHandler.Login)
		r.Get("/github/callback", githubLoginHandler.Callback)
	})

	// serve  static  sites
	fileServer := http.FileServer(http.Dir("./static/"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// task routes
	r.Group(func(r chi.Router) {
		r.Use(jwtAuthenticator.Authenticator)

		r.Route("/task", func(r chi.Router) {
			r.Post("/create", taskHandler.CreateTask)
			r.Get("/get/{task-id}", taskHandler.GetTask)
			r.Get("/get/all", taskHandler.GetAllTasks)
			r.Delete("/delete/{task-id}", taskHandler.DeleteTask)
			r.Put("/update/{task-id}", taskHandler.UpdateTask)
		})
	})
}
