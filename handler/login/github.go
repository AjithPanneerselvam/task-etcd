package login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AjithPanneerselvam/task-etcd/auth"
	"github.com/AjithPanneerselvam/task-etcd/client/github"
	log "github.com/sirupsen/logrus"
)

const (
	githubRedirectURLFormat = "%s?client_id=%s&redirect_uri=%s"
)

type GithubLoginHandler struct {
	githubClient      *github.Client
	githubCallbackURL string

	jwtAuthenticator        *auth.JWTAuth
	loginSuccessRedirectURL string
}

type UserInfo struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

func NewGithubLoginHandler(githubClient *github.Client, githubCallbackURL string, jwtAuthenticator *auth.JWTAuth, loginSuccessRedirectURL string) *GithubLoginHandler {
	return &GithubLoginHandler{
		githubClient:            githubClient,
		githubCallbackURL:       githubCallbackURL,
		loginSuccessRedirectURL: loginSuccessRedirectURL,
		jwtAuthenticator:        jwtAuthenticator,
	}
}

func (g *GithubLoginHandler) Home(w http.ResponseWriter, r *http.Request) {
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
	log.Infof("user %v of id %v signed in", userInfo.Name, userInfo.ID)

	claims := map[string]interface{}{
		"userID": userInfo.ID,
	}

	_, jwtTokenString, err := g.jwtAuthenticator.CreateToken(claims)
	if err != nil {
		log.Error("error creating jwt token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var uInfo = UserInfo{
		ID:    strconv.Itoa(userInfo.ID),
		Token: jwtTokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(uInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
