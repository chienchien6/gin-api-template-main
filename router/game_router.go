package router

import (
	v1 "RCSP/api/v1"
	"github.com/gin-gonic/gin"
)

func InitGameRouter(engine *gin.Engine) {
	gameGroup := engine.Group("game")
	{
		gameGroup.GET("/bet/:user", v1.GetUserBalance)
		gameGroup.GET("/bet/:user/:amount", v1.Bet)
		gameGroup.GET("/prize", v1.GetCurrentPrize)
		gameGroup.GET("/bets", v1.GetUserBets)
	}

}
