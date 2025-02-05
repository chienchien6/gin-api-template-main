package service

import (
	"RCSP/data"
	"github.com/gin-gonic/gin"
)

func GameServer() {
	data.GameServer()
}

func GetUserBalance(c *gin.Context) {
	data.GetUserBalance(c)
}

func Bet(c *gin.Context) {
	data.Bet(c)
}
