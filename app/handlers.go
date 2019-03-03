package app

import (
	"github.com/geekfil/zoom-api-service/telegram"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/labstack/echo"
	"net/http"
)

func (app App) handlers() {

	apiGroup := app.Echo.Group("api/")
	apiGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if app.Config.Token == ctx.QueryParam("token") {
				return next(ctx)
			}
			return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		}
	})

	telegramGroup := apiGroup.Group("/telegram")

	telegramGroup.GET("/send", func(ctx echo.Context) error {
		var tg = ctx.Get("tg").(*telegram.Telegram)
		var text string
		if text = ctx.QueryParam("text"); len(text) == 0 {
			return echo.NewHTTPError(400, "text is required")
		}
		if msg, err := tg.Bot.Send(tgbotapi.NewMessage(tg.Config.ChatId, text)); err != nil {
			return echo.NewHTTPError(400, err)
		} else {
			return ctx.JSON(200, msg)
		}

	})
	telegramGroup.GET("/send/errors", func(ctx echo.Context) error {
		var tg = ctx.Get("tg").(*telegram.Telegram)
		return ctx.JSON(200, tg.SendErrors)
	})
	telegramGroup.GET("/send/errors/clear", func(ctx echo.Context) error {
		var tg = ctx.Get("tg").(*telegram.Telegram)
		tg.SendErrors = []telegram.SendError{}
		return ctx.JSON(200, map[string]string{
			"message": "OK",
		})
	})
}
