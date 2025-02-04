package v1

import (
	"RCSP/global"
	"RCSP/model"
	"RCSP/model/response"
	"RCSP/service"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"strconv"
)

func GetKey(ctx *gin.Context) {

	global.GvaLogger.Info(global.GvaConfig.Test.Key)
	response.OkWithData(ctx, global.GvaConfig.Test.Key)
}

type User struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func Contact(ctx *gin.Context) {

	var contactForm model.ContactForm

	if err := ctx.ShouldBindJSON(&contactForm); err != nil {
		global.GvaLogger.Sugar().Errorf("%#v", contactForm)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//global.GvaLogger.Sugar().Infof("%#v", contactForm)
	global.GvaLogger.Sugar().Debug("%#v", contactForm)

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	response.OkWithData(ctx, global.GvaConfig.Jwt.Issuer)
	return
}

func GetMember(ctx *gin.Context) {

	id := ctx.Param("id")
	memberService := service.UserService{}
	user, err := memberService.GetUserByID(id)

	if err != nil {
		global.GvaLogger.Sugar().Errorf("%#v", user)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	global.GvaLogger.Sugar().Debug("%#v", user)
	response.OkWithData(ctx, user)
}

func CreateMember(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		global.GvaLogger.Sugar().Errorf("Error binding JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	memberService := service.UserService{}
	createdUser, err := memberService.Create(&user)
	if err != nil {
		global.GvaLogger.Sugar().Errorf("Error creating user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	global.GvaLogger.Sugar().Debug("User created: %#v", createdUser)
	response.OkWithData(ctx, createdUser)
}

func UpdateMember(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		global.GvaLogger.Sugar().Errorf("Error binding JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	memberService := service.UserService{}
	updatedUser, err := memberService.Update(&user)
	if err != nil {
		global.GvaLogger.Sugar().Errorf("Error updating user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	global.GvaLogger.Sugar().Debug("User updated: %#v", updatedUser)
	response.OkWithData(ctx, updatedUser)
}

func DeleteMember(ctx *gin.Context) {
	id := ctx.Param("id")
	memberService := service.UserService{}
	if err := memberService.Delete(id); err != nil {
		global.GvaLogger.Sugar().Errorf("Error deleting user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	global.GvaLogger.Sugar().Debug("User deleted: %s", id)
	response.Ok(ctx)
}

func ExpireSet(c *gin.Context) {
	key := c.PostForm("key")
	value := c.PostForm("value")
	expireStr := c.PostForm("expire") // 獲取過期時間（秒）

	var expire int64
	var err error
	if expireStr != "" {
		expire, err = strconv.ParseInt(expireStr, 10, 64)
		if err != nil {
			global.GvaLogger.Sugar().Errorf("invalid expire time: %v", err)
			response.Error(c, "invalid expire time")
			return
		}
	}
	redisService := service.NewRedisService(global.GvaConfig.Redis.Addr)

	// 若 expire 為零，不設置過期時間
	if expire <= 0 {
		expire = 0 // 這樣會使鍵值永久存在
	}

	err = redisService.Set(key, value, expire)
	if err != nil {
		// 記錄詳細錯誤日誌
		global.GvaLogger.Sugar().Errorf("Error setting value in Redis: %v", err)
		response.Error(c, "could not set value")
		return
	}

	response.OkWithData(c, "value set successfully")

}

func ExpireGet(c *gin.Context) {
	//key := c.Param("key")
	key := c.Query("key")

	//redisService := service.RedisService{} //這樣僅是空得結構體
	redisService := service.NewRedisService(global.GvaConfig.Redis.Addr)

	value, err := redisService.Get(key)
	if errors.Is(err, redis.Nil) {
		global.GvaLogger.Sugar().Errorf(err.Error())
		response.Error(c, "key does not exist or has expired")
		return
	} else if err != nil {
		global.GvaLogger.Sugar().Errorf("Error getting value from Redis: %v", err)
		response.Error(c, "could not get value")
		return
	}
	global.GvaLogger.Sugar().Debug(value)
	response.OkWithData(c, value)

}

func ExpireDelete(c *gin.Context) {
	key := c.Query("key")

	redisService := service.NewRedisService(global.GvaConfig.Redis.Addr)

	err := redisService.Delete(key)
	if err != nil {
		global.GvaLogger.Sugar().Errorf("Error deleting value from Redis: %v", err)
		response.Error(c, "could not delete value")
		return
	}

	response.Ok(c)
}
