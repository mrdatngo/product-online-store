package util

import (
	"encoding/json"
	"fmt"
	_const "github.com/mrdatngo/gin-products/const"
	product_controllers "github.com/mrdatngo/gin-products/controllers/product-controllers"
	model "github.com/mrdatngo/gin-products/models"
	"github.com/sirupsen/logrus"
	"strconv"
)

func GetMessageFromInputProductSearch(userID int64, input *product_controllers.InputSearchProducts) []string {
	var results []string
	var logSearching = &model.EntityUserLog{
		UserID:     userID,
		EventType:  _const.EVENT_TYPE_SEARCHING,
		DetailLogs: nil,
	}
	var logDetails []*model.EntityDetailLog
	if len(input.Name) > 0 {
		logDetails = append(logDetails, &model.EntityDetailLog{
			DataType: _const.DATA_TYPE_STRING,
			Param:    _const.PARAM_PRODUCT_NAME,
			Value:    input.Name,
		})
	}
	if input.MinPrice > 0 {
		logDetails = append(logDetails, &model.EntityDetailLog{
			DataType: _const.DATA_TYPE_NUMBER,
			Param:    _const.PARAM_MIN_PRICE,
			Value:    fmt.Sprintf("%v", input.MinPrice),
		})
	}
	if input.MaxPrice > 0 {
		logDetails = append(logDetails, &model.EntityDetailLog{
			DataType: _const.DATA_TYPE_NUMBER,
			Param:    _const.PARAM_MAX_PRICE,
			Value:    fmt.Sprintf("%v", input.MaxPrice),
		})
	}
	logSearching.DetailLogs = logDetails

	dataBytes, err := json.Marshal(logSearching)
	if err != nil {
		logrus.Errorf("Error on marshall searching log, err: %v", err)
	}
	results = append(results, string(dataBytes))

	if len(input.Branch) > 0 {
		logFiltering := &model.EntityUserLog{
			UserID:    userID,
			EventType: _const.EVENT_TYPE_FILTERING,
			DetailLogs: []*model.EntityDetailLog{
				{
					DataType: _const.DATA_TYPE_NUMBER,
					Param:    _const.PARAM_BRANCH_ID,
					Value:    input.Branch,
				},
			},
		}
		dataBytes, err := json.Marshal(logFiltering)
		if err != nil {
			logrus.Errorf("Error on marshall searching log, err: %v", err)
		}
		results = append(results, string(dataBytes))
	}
	return results
}

func GetMessageFromInputGetProduct(userID int64, input *product_controllers.InputGetProduct) string {
	viewLogs := &model.EntityUserLog{
		UserID:    userID,
		EventType: _const.EVENT_TYPE_VIEWING,
		DetailLogs: []*model.EntityDetailLog{
			{
				DataType: _const.DATA_TYPE_NUMBER,
				Param:    _const.PARAM_PRODUCT_ID,
				Value:    strconv.Itoa(int(input.ProductID)),
			},
		},
	}
	dataBytes, err := json.Marshal(viewLogs)
	if err != nil {
		logrus.Errorf("Error on marshall searching log, err: %v", err)
	}
	return string(dataBytes)
}
