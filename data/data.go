package data

import (
	"RCSP/model"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	RoundSecond    = 60               // 每一局的時間
	DefaultBalance = 1000             // 玩家初始化金額
	UserMember     = "game"           // 儲存所有使用者的Balance   	Redis:`Sorted-Set`	SCORE -> USER
	BetThisRound   = "bet_this_round" // 儲存目前局的下注狀況		 	Redis:`Sorted-Set`  SCORE -> USER
)

var Round = 0
var startTimeThisRound time.Time
var RC *redis.Client

func newClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	log.Println(pong)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}

func init() {
	RC = newClient()

	// 初始化清空所有Redis
	RC.Del(UserMember)
	RC.Del(BetThisRound)

	go GameServer()
}

func GameServer() {
	// 1. 設置隨機種子,使每次隨機結果不同
	rand.Seed(time.Now().UTC().UnixNano())

	// 2. 建立定時器,每隔 RoundSecond 秒執行一次遊戲輪次
	ticker := time.NewTicker(RoundSecond * time.Second)
	go func() {
		for {
			// 3. 增加回合計數
			Round++

			// 4. 記錄本輪開始時間
			startTimeThisRound = time.Now()
			log.Println(startTimeThisRound.Format("2006-01-02 15:04:05"), "\t round", Round, "start")

			// 5. 等待定時器觸發,開始本輪遊戲
			_ = <-ticker.C

			// 6. 獲取當前的獎金池和所有玩家的下注資訊
			var prizePool = getCurrentPrize()
			var userBets = getUserBets()

			// 7. 如果本輪沒有任何玩家下注,則跳過本輪
			if len(userBets) == 0 {
				log.Println("Round", Round, "沒有任何玩家下注")
				continue
			}

			// 8. 隨機選出一個中獎號碼
			winNum := rand.Intn(prizePool + 1)

			// 9. 根據玩家的下注金額,確定中獎玩家
			var winner string
			for _, userBet := range userBets {
				winNum -= userBet.Amount
				if winNum <= 0 {
					winner = userBet.Id
					break
				}
			}
			log.Println("獎金池:", prizePool, "\t 得主:", winner)

			// 10. 將獎金發放給中獎玩家
			RC.ZIncrBy(UserMember, float64(prizePool), winner)

			// 11. 刪除本輪的下注記錄
			RC.Del(BetThisRound)
		}
	}()
}

func GetUserBalance(c *gin.Context) {
	var user model.Player
	user.Id = c.Param("user")
	balance, err := RC.ZScore(UserMember, user.Id).Result()
	if err == redis.Nil { //查無使用者，註冊新帳號
		balance = DefaultBalance
		RC.ZAdd(UserMember, redis.Z{
			Score:  balance,
			Member: user.Id,
		})
	}
	user.Balance = int(balance)
	wrapResponse(c, user, nil)

}

func Bet(c *gin.Context) {
	var user model.Player
	user.Id = c.Param("user")
	amountStr := c.Param("amount")
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		wrapResponse(c, nil, errors.New("下注金額有誤"))
		return
	}

	balance, err := RC.ZScore(UserMember, user.Id).Result()
	if err == redis.Nil {
		wrapResponse(c, nil, errors.New("查無此用戶，請先註冊"))
		return
	}
	user.Balance = int(balance)
	if amount <= 0 {
		wrapResponse(c, nil, errors.New("下注金額需為正整數"))
		return
	}
	if amount > user.Balance {
		wrapResponse(c, nil, errors.New("餘額不足"))
		return
	}

	user.Balance -= amount
	RC.ZIncrBy(UserMember, float64(-amount), user.Id)
	RC.ZIncrBy(BetThisRound, float64(amount), user.Id)

	wrapResponse(c, user, nil)
}

func GetCurrentPrize(c *gin.Context) {
	wrapResponse(c, getCurrentPrize(), nil)
}

func GetUserBets(c *gin.Context) {
	UserBets := getUserBets()
	if len(UserBets) == 0 {
		wrapResponse(c, nil, errors.New("目前沒有任何記錄"))
		return
	}
	wrapResponse(c, UserBets, nil)
}

func getCurrentPrize() (prizePool int) {
	bets, _ := RC.ZRangeWithScores(BetThisRound, 0, -1).Result()
	for _, bet := range bets {
		var userBet model.UserBet
		userBet.Amount = int(bet.Score)
		prizePool += userBet.Amount
	}
	return
}

func getUserBets() (userBets []model.UserBet) {
	bets, _ := RC.ZRangeWithScores(BetThisRound, 0, -1).Result()
	for _, bet := range bets {
		var userBet model.UserBet
		userBet.Id = fmt.Sprintf("%s", bet.Member)
		userBet.Amount = int(bet.Score)
		userBet.Round = Round
		userBets = append(userBets, userBet)
	}
	return
}

func wrapResponse(c *gin.Context, data interface{}, err error) {
	type ret struct {
		Status string      `json:"status"`
		Msg    string      `json:"msg"`
		Data   interface{} `json:"data"`
	}

	d := ret{
		Status: "ok",
		Msg:    "",
		Data:   []struct{}{},
	}

	if data != nil {
		d.Data = data
	}

	if err != nil {
		d.Status = "failed"
		d.Msg = err.Error()
	}

	c.JSON(http.StatusOK, d)
}
