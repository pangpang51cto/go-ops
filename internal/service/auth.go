package service

import (
	"context"
	v1 "go-ops/api/v1"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	jwt "github.com/gogf/gf-jwt/v2"
)

var authService *jwt.GfJWTMiddleware

func Auth() *jwt.GfJWTMiddleware {
	return authService
}

func init() {
	auth := jwt.New(&jwt.GfJWTMiddleware{
		Realm:           "test zone",
		Key:             []byte("secret key"),
		Timeout:         time.Minute * 10,
		MaxRefresh:      time.Minute * 20,
		IdentityKey:     "uid",
		TokenLookup:     "header: GO-OPS-X-TOKEN, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
		Authenticator:   Authenticator,
		Unauthorized:    Unauthorized,
		PayloadFunc:     PayloadFunc,
		IdentityHandler: IdentityHandler,
	})
	authService = auth
}

// PayloadFunc is a callback function that will be called during login.
// Using this function it is possible to add additional payload data to the webtoken.
// The data is then made available during requests via c.Get("JWT_PAYLOAD").
// Note that the payload is not encrypted.
// The attributes mentioned on jwt.io can't be used as keys for the map.
// Optional, by default no additional data will be set.
func PayloadFunc(data interface{}) jwt.MapClaims {
	claims := jwt.MapClaims{}
	params := data.(map[string]interface{})
	if len(params) > 0 {
		for k, v := range params {
			claims[k] = v
		}
	}
	return claims
}

// IdentityHandler get the identity from JWT and set the identity for every request
// Using this function, by r.GetParam("id") get identity
func IdentityHandler(ctx context.Context) interface{} {
	claims := jwt.ExtractClaims(ctx)
	return claims[authService.IdentityKey]
}

// Unauthorized is used to define customized Unauthorized callback function.
func Unauthorized(ctx context.Context, code int, message string) {

	r := g.RequestFromCtx(ctx)

	r.Response.WriteStatus(code, g.Map{
		"code":    code,
		"message": message,
	})
	r.ExitAll()
}

// Authenticator is used to validate login parameters.
// It must return user data as user identifier, it will be stored in Claim Array.
// if your identityKey is 'id', your user data must have 'id'
// Check error (e) to determine the appropriate error message.
func Authenticator(ctx context.Context) (interface{}, error) {
	var (
		r  = g.RequestFromCtx(ctx)
		in v1.AuthLoginReq
	)
	if err := r.Parse(&in); err != nil {
		return "", err
	}

	user, err := User().SignIn(ctx, &in)

	if err != nil {
		return nil, err
	}

	if user != nil {
		return g.Map{
			"uid":      user.Uid,
			"username": user.Username,
		}, nil
	}

	return nil, jwt.ErrFailedAuthentication
}
