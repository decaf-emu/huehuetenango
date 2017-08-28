package api

import (
	"context"
	"crypto/rand"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/go-github/github"
	"github.com/labstack/echo"
	"github.com/oklog/ulid"
	"golang.org/x/oauth2"
)

type JWTClaims struct {
	AuthToken   string        `json:"auth_token"`
	GithubToken *oauth2.Token `json:"github_token"`
	Name        string        `json:"name"`
	AvatarURL   string        `json:"avatar_url"`
	jwt.StandardClaims
}

func (a *api) Login(c echo.Context) error {
	authToken, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		return err
	}

	url := a.authConfig.AuthCodeURL(authToken.String(), oauth2.AccessTypeOnline)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTClaims{
		AuthToken: authToken.String(),
	})

	t, err := token.SignedString([]byte(a.jwtSigningSecret))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, struct {
		URL   string `json:"url"`
		Token string `json:"token"`
	}{
		URL:   url,
		Token: t,
	})
}

type loginCallbackRequest struct {
	State string `json:"state"`
	Code  string `json:"code"`
}

func (a *api) LoginCallback(c echo.Context) error {
	request := &loginCallbackRequest{}
	if err := c.Bind(request); err != nil {
		return err
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)

	if claims.AuthToken != request.State {
		return c.NoContent(http.StatusBadRequest)
	}

	oauthToken, err := a.authConfig.Exchange(oauth2.NoContext, request.Code)
	if err != nil {
		return err
	}

	if !oauthToken.Valid() {
		return c.NoContent(http.StatusBadRequest)
	}

	client := github.NewClient(a.authConfig.Client(oauth2.NoContext, oauthToken))
	githubUser, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		return err
	}

	claims = &JWTClaims{
		GithubToken: oauthToken,
	}

	if githubUser != nil {
		if githubUser.Name != nil {
			claims.Name = *githubUser.Name
		}
		if githubUser.AvatarURL != nil {
			claims.AvatarURL = *githubUser.AvatarURL
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(a.jwtSigningSecret))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: t,
	})
}
