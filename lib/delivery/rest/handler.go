package rest

import (
	"encoding/json"
	modelmonopoly "github.com/bayu-aditya/bayu-aditya-backend/lib/core/model/monopoly"
	"github.com/bayu-aditya/bayu-aditya-backend/lib/core/util"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

func (h *Handler) HealthCheck(c *gin.Context) {
	responseOK(c, map[string]string{
		"time": time.Now().Format(time.RFC3339),
	})
}

type monopolyCreateRoomRequestJson struct {
	InitialBalance int64 `json:"initial_balance"`
	Player         struct {
		ID   string `json:"id" binding:"required"`
		Name string `json:"name" binding:"required"`
	} `json:"player" binding:"required"`
}

type monopolyCreateRoomResponseJson struct {
	RoomID   string `json:"room_id"`
	RoomPass string `json:"room_pass"`
}

func (h *Handler) MonopolyCreateRoom(c *gin.Context) {
	requestJson := monopolyCreateRoomRequestJson{}
	responseJson := monopolyCreateRoomResponseJson{}

	if err := c.BindJSON(&requestJson); err != nil {
		responseBadRequest(c, err.Error())
		return
	}

	player := modelmonopoly.StatePlayer{
		ID:   requestJson.Player.ID,
		Name: requestJson.Player.Name,
	}

	roomID, roomPass, err := h.usecaseMonopolyBanking.CreateRoom(c, player, requestJson.InitialBalance)
	if err != nil {
		responseInternalServer(c, err.Error())
		return
	}

	responseJson.RoomID = roomID
	responseJson.RoomPass = roomPass

	responseOK(c, responseJson)
}

type monopolyJoinRoomRequestJson struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	RoomID     string `json:"room_id"`
	RoomPass   string `json:"room_pass"`
}

func (h *Handler) MonopolyJoinRoom(c *gin.Context) {
	requestJson := monopolyJoinRoomRequestJson{}

	if err := c.BindJSON(&requestJson); err != nil {
		responseBadRequest(c, err.Error())
		return
	}

	playerID := requestJson.PlayerID
	playerName := requestJson.PlayerName
	roomID := requestJson.RoomID
	roomPass := requestJson.RoomPass

	if err := h.usecaseMonopolyBanking.JoinRoom(c, playerID, playerName, roomID, roomPass); err != nil {
		responseInternalServer(c, err.Error())
		return
	}

	responseOK(c, nil)
}

type monopolyGetStateRequestQuery struct {
	PlayerID string `form:"player_id" binding:"required"`
}

type monopolyGetStateResponseJson struct {
	Player  modelmonopoly.StatePlayer   `json:"player"`
	Players []modelmonopoly.StatePlayer `json:"players"`
	Logs    []modelmonopoly.StateLog    `json:"logs"`
}

func (h *Handler) MonopolyGetState(c *gin.Context) {
	requestQuery := monopolyGetStateRequestQuery{}
	responseJson := monopolyGetStateResponseJson{}
	roomID := c.Param("room_id")

	if err := c.BindQuery(&requestQuery); err != nil {
		responseBadRequest(c, err.Error())
		return
	}

	state, err := h.usecaseMonopolyBanking.GetState(c, requestQuery.PlayerID, roomID)
	if err != nil {
		responseInternalServer(c, err.Error())
		return
	}

	for i, player := range state.Players {
		if player.ID == requestQuery.PlayerID {
			responseJson.Player = player
			state.Players = util.ArrayRemove(state.Players, i)
		}
	}

	responseJson.Players = state.Players
	responseJson.Logs = state.Logs

	responseOK(c, responseJson)
}

type monopolyCreateTransactionRequestJson struct {
	PlayerID       string `json:"player_id"`
	TargetPlayerID string `json:"target_player_id"`
	Type           string `json:"type"`
	Amount         int64  `json:"amount"`
}

func (h *Handler) MonopolyCreateTransaction(c *gin.Context) {
	requestJson := monopolyCreateTransactionRequestJson{}

	if err := c.BindJSON(&requestJson); err != nil {
		responseBadRequest(c, err.Error())
		return
	}

	roomID := c.Param("room_id")
	playerID := requestJson.PlayerID
	targetPlayerID := requestJson.TargetPlayerID
	amount := requestJson.Amount
	mode := requestJson.Type

	if err := h.usecaseMonopolyBanking.CreateTransaction(c, roomID, playerID, targetPlayerID, amount, mode); err != nil {
		responseInternalServer(c, err.Error())
		return
	}

	responseOK(c, nil)
}

func (h *Handler) MonopolyServerSentEvent(c *gin.Context) {
	timeoutChan := time.After(60 * time.Second)

	roomID := c.Param("room_id")

	stateChan, stop, err := h.usecaseMonopolyBanking.SubscribeState(c, roomID)
	if err != nil {
		return
	}
	defer stop()

	c.Stream(func(w io.Writer) bool {
		select {
		case state := <-stateChan:
			stateJson, err := json.Marshal(state)
			if err != nil {
				return false
			}

			c.SSEvent("message", stateJson)
			return true

		case <-timeoutChan:
			return false
		case <-c.Request.Context().Done():
			return false
		}
	})
}
