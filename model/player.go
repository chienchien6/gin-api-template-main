package model

type Player struct {
	Id      string `json:"Id"`
	Balance int    `json:"balance"`
}

type UserBet struct {
	Id     string `json:"Id"`
	Round  int    `json:"round"`  // 局數
	Amount int    `json:"amount"` // 下注金額
}
