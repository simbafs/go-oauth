package server

import (
	"github.com/gin-gonic/gin"
	"github.com/simba-fs/go-oauth/types"

	"github.com/simba-fs/go-oauth/server/handler"
)

func New(config *types.Config) {
	app := gin.Default()

	h := handler.Handler{Config: config}

	app.Use(func(ctx *gin.Context){
		ctx.Set("config", config)
	})

	app.GET("/github/login", h.Login)
	app.GET("/github/callback", h.GithubCallback)

	app.Run(config.Addr)
}
