package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

type Client struct {
	*http.Client
	oAuthURL     string
	apiURL       string
	clientID     string
	clientSecret string
	accessToken  string
}

func New(oAuthURL string, apiURL string, clientID string, clientSecret string, timeoutInSec int32) *Client {
	return &Client{
		oAuthURL:     oAuthURL,
		apiURL:       apiURL,
		clientID:     clientID,
		clientSecret: clientSecret,
		Client: &http.Client{
			Timeout: time.Duration(timeoutInSec) * time.Second,
		},
	}
}

func (c *Client) GetAccessToken(ctx context.Context, code string) (string, error) {
	var accessTokenRequest = AccessTokenRequest{
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
		Code:         code,
	}

	body, err := json.Marshal(accessTokenRequest)
	if err != nil {
		return "", errors.Wrap(err, "error marshalling access token request")
	}

	oAuthAccessTokenURL := fmt.Sprintf("%s/access_token", c.oAuthURL)

	oAuthAccessTokenParsedURL, err := url.Parse(oAuthAccessTokenURL)
	if err != nil {
		return "", errors.Wrapf(err, "error parsing oauth access token url :%v", oAuthAccessTokenURL)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, oAuthAccessTokenParsedURL.String(),
		bytes.NewReader(body))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "error making request")
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "error reading response body")
	}

	var accessTokenResponse AccessTokenResponse
	err = json.Unmarshal(respBody, &accessTokenResponse)
	if err != nil {
		return "", errors.Wrap(err, "error unmarshalling access token response body")
	}

	return accessTokenResponse.AccessToken, nil
}

func (c *Client) GetRedirectAuthorizeURL(ctx context.Context, callbackURL string) (string, error) {
	authorizeURL := fmt.Sprintf("%s/authorize?client_id=%s&redirect_uri=%s", c.oAuthURL, c.clientID, callbackURL)

	authorizeParsedURL, err := url.Parse(authorizeURL)
	if err != nil {
		return "", errors.Wrap(err, "error parsing github authorize url")
	}

	return authorizeParsedURL.String(), nil
}

func (c *Client) SetAccessToken(ctx context.Context, accessToken string) {
	c.accessToken = accessToken
}

func (c *Client) GetUserInfo(ctx context.Context) (*UserInfo, error) {
	githubUserInfoURL := fmt.Sprintf("%s/user", c.apiURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubUserInfoURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating github user info request")
	}

	authToken := fmt.Sprintf("token %s", c.accessToken)
	req.Header.Set("Authorization", authToken)

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error making request")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading response body")
	}

	var userInfo UserInfo
	err = json.Unmarshal(respBody, &userInfo)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling user info body")
	}

	return &userInfo, nil
}
