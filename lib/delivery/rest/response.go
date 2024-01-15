package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type response struct {
	Meta responseMeta `json:"meta,omitempty"`
	Data interface{}  `json:"data,omitempty"`
}

type responseMeta struct {
	Error string `json:"error,omitempty"`
}

func responseOK(c *gin.Context, body interface{}) {
	c.JSON(http.StatusOK, response{
		Data: body,
	})
}

func responseBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, response{
		Meta: responseMeta{
			Error: message,
		},
	})
}

func responseInternalServer(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, response{
		Meta: responseMeta{
			Error: message,
		},
	})
}
