package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwt"
	Jwt "github.com/simba-fs/go-oauth/jwt"
	"github.com/simba-fs/go-oauth/types"
)

var (
	ErrCannotGetConfig = errors.New("cannot get config in gin context")
)

// getState return state
func getState() string {
	// TODO: random generate state
	return "login"
}

type Handler struct {
	Config *types.Config
}

func (h *Handler) GithubCallback(ctx *gin.Context) {
	
	code := ctx.Query("code")
	c, ok := ctx.Get("config")
	if !ok {
		panic(ErrCannotGetConfig)
	}

	config := c.(*types.Config)

	// get user token
	res, err := http.PostForm("https://github.com/login/oauth/access_token", url.Values{
		"client_id":     {h.Config.ClientID},
		"client_secret": {h.Config.ClientSecret},
		"code":          {code},
		"state":         {getState()},
	})

	if err != nil {
		ctx.Data(http.StatusBadRequest, "text/plain", []byte(err.Error()))
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		ctx.Data(http.StatusBadRequest, "text/plain", []byte(err.Error()))
		return
	}

	value, err := url.ParseQuery(string(body))
	if err != nil {
		ctx.Data(http.StatusBadRequest, "text/plain", []byte(err.Error()))
		return
	}
	userToken := value.Get("access_token")

	// get email
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		ctx.Data(http.StatusBadRequest, "text/plain", []byte(err.Error()))
		return
	}

	req.Header.Add("Authorization", "token "+userToken)
	res, err = client.Do(req)
	if err != nil {
		ctx.Data(http.StatusBadRequest, "text/plain", []byte(err.Error()))
		return
	}

	if res.StatusCode != http.StatusOK {
		ctx.JSON(res.StatusCode, map[string]string{
			"data": res.Status,
		})
		return
	}

	emails := []types.Email{}
	json.NewDecoder(res.Body).Decode(&emails)

	email := ""
	for _, v := range emails {
		if v.Primary {
			email = v.Email
		}
	}

	token := jwt.New()
	token.Set("email", email)
	// token.Set("username", )
	signed, err := Jwt.Sign(&token, config)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, map[string]string{
		"data": string(signed),
	})
}

func (h *Handler) Login(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&state=%s&allow_signup=%s",
		h.Config.ClientID,
		getState(),
		h.Config.AllowSignup,
	))
}
