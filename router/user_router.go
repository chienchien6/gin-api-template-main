/**
 * @Description 用户相关的路由
 **/
package router

import (
	v1 "RCSP/api/v1"
	"RCSP/middleware"

	"github.com/gin-gonic/gin"
)

/**
 * @description: 用户相关的路由
 * @param engine
 */
func InitUserRouter(engine *gin.Engine) {
	// 不需要登录的路由
	noLoginGroup := engine.Group("v1/user")
	{
		// 登录
		noLoginGroup.POST("login", v1.Login)
		// 注册
		noLoginGroup.POST("register", v1.Register)

		noLoginGroup.GET("permit", v1.GetKey)

		noLoginGroup.GET("getExpireTime", v1.ExpireGet)
		noLoginGroup.POST("setExpireTime", v1.ExpireSet)
		noLoginGroup.DELETE("deleteExpireTime", v1.ExpireDelete)
	}
	// 需要登录
	tokenGroup := engine.Group("v1/user").Use(middleware.JWTAuthMiddleware())
	{
		tokenGroup.POST("/detail", v1.GetUser)
		tokenGroup.POST("/contact", v1.Contact)
		tokenGroup.GET("/getMember/:id", v1.GetMember)
		tokenGroup.POST("/createMember", v1.CreateMember)
		tokenGroup.PUT("/updateMember", v1.UpdateMember)
		tokenGroup.DELETE("/deleteMember/:id", v1.DeleteMember)
	}
}
