package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "jwtauth context value " + k.name
}

var (
	TokenCtxKey = &contextKey{"Token"}
	ErrorCtxKey = &contextKey{"Error"}
)

type JWTAuth struct {
	secretKey          []byte
	signatureAlgorithm jwa.SignatureAlgorithm

	expiryDuration time.Duration
}

func NewJWTAuth(secretKey string, expiryDuration time.Duration) *JWTAuth {
	return &JWTAuth{
		signatureAlgorithm: jwa.HS256,
		secretKey:          []byte(secretKey),
		expiryDuration:     expiryDuration,
	}
}

func (j *JWTAuth) CreateToken(claims map[string]interface{}) (jwt.Token, string, error) {
	token := jwt.New()

	for key, value := range claims {
		token.Set(key, value)
	}

	currentTime := time.Now().UTC().Unix()

	// standard claims
	err := token.Set(jwt.IssuedAtKey, currentTime)
	if err != nil {
		return nil, "", err
	}

	err = token.Set(jwt.ExpirationKey, currentTime+int64(j.expiryDuration))
	if err != nil {
		return nil, "", err
	}

	signedToken, err := jwt.Sign(token, j.signatureAlgorithm, j.secretKey)
	if err != nil {
		return nil, "", err
	}

	return token, string(signedToken), nil
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := FromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func FromContext(ctx context.Context) (jwt.Token, map[string]interface{}, error) {
	token, ok := ctx.Value(TokenCtxKey).(jwt.Token)
	if !ok {
		return nil, nil, errors.New("error type asserting token into jwt token")
	}

	var err error
	var claims map[string]interface{}

	if token != nil {
		claims, err = token.AsMap(context.Background())
		if err != nil {
			return token, nil, err
		}
	} else {
		claims = map[string]interface{}{}
	}

	err, ok = ctx.Value(ErrorCtxKey).(error)

	return token, claims, err
}
