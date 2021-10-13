package login

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AjithPanneerselvam/task-etcd/auth"
	"github.com/AjithPanneerselvam/task-etcd/client/github"
	"github.com/AjithPanneerselvam/task-etcd/store"
	userstore "github.com/AjithPanneerselvam/task-etcd/store/user"
	log "github.com/sirupsen/logrus"
)

const (
	githubRedirectURLFormat = "%s?client_id=%s&redirect_uri=%s"
)

type GithubLoginHandler struct {
	githubClient      *github.Client
	githubCallbackURL string

	jwtAuthenticator *auth.JWTAuth

	loginSuccessRedirectURL string
	userStore               store.UserStore
}

func NewGithubLoginHandler(githubClient *github.Client, githubCallbackURL string, jwtAuthenticator *auth.JWTAuth,
	userStore store.UserStore, loginSuccessRedirectURL string) *GithubLoginHandler {
	return &GithubLoginHandler{
		githubClient:            githubClient,
		githubCallbackURL:       githubCallbackURL,
		loginSuccessRedirectURL: loginSuccessRedirectURL,
		userStore:               userStore,
	}
}

func (g *GithubLoginHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<a href="/login/github">Github Login</a>`)
}

func (g *GithubLoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	redirectURL, err := g.githubClient.GetRedirectAuthorizeURL(ctx, g.githubCallbackURL)
	if err != nil {
		log.Errorf("error fetching github redirect url: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Infof("login redirecting to URL: %v", redirectURL)
	http.Redirect(w, r, redirectURL, 301)
}

func (g *GithubLoginHandler) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	code, ok := r.URL.Query()["code"]
	if !ok {
		log.Error("error as query param 'code' is missing")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Debugf("Auth code: %v", code)

	githubAccessToken, err := g.githubClient.GetAccessToken(ctx, code[0])
	if err != nil {
		log.Errorf("error fetching github access token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	g.githubClient.SetAccessToken(ctx, githubAccessToken)
	log.Debugf("Github access token: %v", githubAccessToken)

	userInfo, err := g.githubClient.GetUserInfo(ctx)
	if err != nil {
		log.Error("error fetching user info: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Debugf("github user info: %v", userInfo)

	_, err = g.userStore.GetInfoByID(ctx, strconv.Itoa(userInfo.ID))
	if err != nil && err != userstore.ErrUserStoreNoRecord {
		log.Errorf("error fetching user info for id: %v: %v", userInfo.ID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err == userstore.ErrUserStoreNoRecord {
		var user = store.User{
			ID:        userInfo.ID,
			Name:      userInfo.Name,
			Handle:    userInfo.Login,
			Email:     userInfo.Email,
			CreatedAt: time.Now(),
		}
		log.Infof("creating a new user of id: %v, name: %v", user.ID, user.Name)

		err = g.userStore.CreateUser(ctx, user)
		if err != nil {
			log.Errorf("error creating a new user of id: %v : %v", user.ID, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info("created a new user of id: %v, name: %v", user.ID, user.Name)
	}

	log.Infof("user %v of id %v signed in", userInfo.Name, userInfo.ID)
	log.Infof("redirecting to URL: %v", g.loginSuccessRedirectURL)

	http.Redirect(w, r, g.loginSuccessRedirectURL, 301)
}
