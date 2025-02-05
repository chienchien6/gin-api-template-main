package v1

import (
	"RCSP/data"
	"RCSP/service"
	"github.com/gin-gonic/gin"
)

func GetUserBalance(c *gin.Context) {
	service.GetUserBalance(c)
}

func Bet(c *gin.Context) {
	service.Bet(c)
}

func GetCurrentPrize(c *gin.Context) {
	data.GetCurrentPrize(c)

}

func GetUserBets(c *gin.Context) {
	data.GetUserBets(c)
}
