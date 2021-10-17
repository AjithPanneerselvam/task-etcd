package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	log "github.com/sirupsen/logrus"
)

const (
	ClaimsKeyUserID = "userID"
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return k.name
}

var (
	TokenCtxKey = &contextKey{"Token"}
)

type JWTAuth struct {
	secretKey          []byte
	signatureAlgorithm jwa.SignatureAlgorithm
	expiryDuration     time.Duration
	verifier           jwt.ParseOption
}

func NewJWTAuth(secretKey string, expiryDuration time.Duration) *JWTAuth {
	return &JWTAuth{
		signatureAlgorithm: jwa.HS256,
		secretKey:          []byte(secretKey),
		expiryDuration:     expiryDuration,
		verifier:           jwt.WithVerify(jwa.HS256, []byte(secretKey)),
	}
}

func (j *JWTAuth) CreateToken(claims map[string]interface{}) (string, error) {
	token := jwt.New()
	for key, value := range claims {
		token.Set(key, value)
	}

	// standard claims
	currentTime := time.Now().UTC().Unix()
	err := token.Set(jwt.IssuedAtKey, currentTime)
	if err != nil {
		return "", err
	}

	err = token.Set(jwt.ExpirationKey, int64(currentTime+int64(j.expiryDuration.Seconds())))
	if err != nil {
		return "", err
	}

	signedToken, err := jwt.Sign(token, j.signatureAlgorithm, j.secretKey)
	if err != nil {
		return "", err
	}

	return string(signedToken), nil
}

func (j *JWTAuth) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := FetchBearerToken(r)

		if tokenString == "" {
			log.Error("error as authorization token is empty")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse([]byte(tokenString), j.verifier)
		if err != nil {
			log.Errorf("error verifying token: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.Debug("parsed token: %v", token)

		if err := jwt.Validate(token); err != nil {
			log.Errorf("error validating token: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), TokenCtxKey, token)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FetchBearerToken fetches the bearer token from the request header
func FetchBearerToken(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}

	return ""
}

func FetchClaimValFromCtx(ctx context.Context, claimKey string) (interface{}, error) {
	token, err := FetchTokenFromCtx(ctx)
	if err != nil {
		return "", err
	}

	claims := make(map[string]interface{})

	if token != nil {
		claims, err = token.AsMap(context.Background())
		if err != nil {
			return "", err
		}
	}

	claimVal, ok := claims[claimKey]
	if !ok {
		return "", errors.New("")
	}

	return claimVal, nil
}

func FetchTokenFromCtx(ctx context.Context) (jwt.Token, error) {
	token, ok := ctx.Value(TokenCtxKey).(jwt.Token)
	if !ok {
		return nil, errors.New("error type asserting token into jwt token")
	}

	return token, nil
}
