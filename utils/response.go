package util

import "github.com/gin-gonic/gin"

type Response struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

type Meta struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type ErrorResponse struct {
	Data interface{} `json:"data"`
	Meta ErrorMeta   `json:"meta"`
}

type ErrorMeta struct {
	StatusCode int         `json:"status_code"`
	Message    interface{} `json:"message"`
}

func APIResponse(ctx *gin.Context, Message string, StatusCode int, Data interface{}) {

	jsonResponse := Response{
		Data: Data,
		Meta: Meta{
			StatusCode: StatusCode,
			Message:    Message,
		},
	}

	if StatusCode >= 400 {
		ctx.JSON(StatusCode, jsonResponse)
		defer ctx.AbortWithStatus(StatusCode)
	} else {
		ctx.JSON(StatusCode, jsonResponse)
	}
}

func ExecuteErrorResponse(ctx *gin.Context, StatusCode int, err interface{}) {
	errResponse := ErrorResponse{
		Data: nil,
		Meta: ErrorMeta{
			StatusCode: StatusCode,
			Message:    err,
		},
	}
	ctx.JSON(StatusCode, errResponse)
	defer ctx.AbortWithStatus(StatusCode)
}
