package rest

import (
	"encoding/json"
	modelmonopoly "github.com/bayu-aditya/bayu-aditya-backend/lib/core/model/monopoly"
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

type monopolyServerSentEventRequestQuery struct {
	RoomID string `form:"room_id" binding:"required"`
}

func (h *Handler) MonopolyServerSentEvent(c *gin.Context) {
	requestQuery := monopolyServerSentEventRequestQuery{}
	timeoutChan := time.After(60 * time.Second)

	if err := c.BindQuery(&requestQuery); err != nil {
		responseBadRequest(c, err.Error())
		return
	}

	roomID := requestQuery.RoomID

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

			c.SSEvent("state", stateJson)
			return true

		case <-timeoutChan:
			return false
		case <-c.Request.Context().Done():
			return false
		}
	})
}
