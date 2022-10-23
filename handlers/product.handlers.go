package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/mrdatngo/gin-products/connectors"
	_const "github.com/mrdatngo/gin-products/const"
	product_controllers "github.com/mrdatngo/gin-products/controllers/product-controllers"
	model "github.com/mrdatngo/gin-products/models"
	util "github.com/mrdatngo/gin-products/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type handler struct {
	service       product_controllers.Service
	kafkaProducer *connectors.KafkaProducer
	redisClient   *redis.Client
}

func NewProductHandler(service product_controllers.Service, kafkaProducer *connectors.KafkaProducer, redisClient *redis.Client) *handler {
	return &handler{service: service, kafkaProducer: kafkaProducer, redisClient: redisClient}
}

func (h *handler) SearchProductsHandler(ctx *gin.Context) {
	inputSearch := product_controllers.InputSearchProducts{}
	err := ctx.ShouldBindQuery(&inputSearch)

	if inputSearch.Page == 0 {
		inputSearch.Page = 1
	}
	if inputSearch.PageSize == 0 {
		inputSearch.PageSize = 20
	}
	if err != nil {
		util.ExecuteErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	if inputSearch.MaxPrice > 0 && inputSearch.MinPrice > inputSearch.MaxPrice {
		errData := util.FormatResponseErrorWithMessage("min_price & max_price", "min_price must smaller than max_price!")
		util.ExecuteErrorResponse(ctx, http.StatusBadRequest, errData)
		return
	}
	go func() {
		data := util.GetMessageFromInputProductSearch(0, &inputSearch)
		for _, item := range data {
			h.kafkaProducer.SendMessageData("user_log", item)
		}
	}()
	products, errCode := h.service.SearchProducts(&inputSearch)
	if errCode != "" {
		switch errCode {
		case _const.UNKNOWN_PANIC:
			util.APIResponse(ctx, "Something went wrong!", http.StatusInternalServerError, nil)
			return
		default:
			logrus.Errorf("Error on search products, errCode: %v", errCode)
			util.ExecuteErrorResponse(ctx, http.StatusBadRequest, "Something went wrong")
			return
		}
	}
	util.APIResponse(ctx, "OK", http.StatusOK, map[string]interface{}{
		"products": products,
	})
}

func (h *handler) GetProductHandler(ctx *gin.Context) {
	inputGet := product_controllers.InputGetProduct{}

	if id, err := strconv.Atoi(ctx.Param("product_id")); err != nil {
		errData := util.FormatResponseErrorWithMessage("product_id", "Invalid product id!")
		util.ExecuteErrorResponse(ctx, http.StatusBadRequest, errData)
		return
	} else {
		inputGet.ProductID = int64(id)
	}

	go func() {
		h.kafkaProducer.SendMessageData("user_log", util.GetMessageFromInputGetProduct(0, &inputGet))
	}()

	// check valid cache product
	if h.redisClient != nil {
		val, err := h.redisClient.Get(ctx, fmt.Sprintf("productId_%v", inputGet.ProductID)).Result()
		if err != nil {
			logrus.Errorf("Cannot get product id %v in cache, process to db!", inputGet.ProductID)
		} else {
			product := &model.EntityProduct{}
			err = json.Unmarshal([]byte(val), product)
			if err != nil {
				logrus.Errorf("Error in parse get product id %v in cache, process to db!", inputGet.ProductID)
			} else {
				util.APIResponse(ctx, "OK", http.StatusOK, map[string]interface{}{
					"product": product,
				})
				return
			}
		}
	} else {
		logrus.Infof("redisClient has beeen nil!")
	}

	product, errCode := h.service.GetProduct(&inputGet)
	if errCode != "" {
		switch errCode {
		case _const.UNKNOWN_PANIC:
			util.APIResponse(ctx, "Something went wrong!", http.StatusInternalServerError, nil)
			return
		case _const.PRODUCT_NOT_FOUND:
			util.APIResponse(ctx, "No product found!", http.StatusNotFound, nil)
			return
		default:
			logrus.Errorf("Error on search products, errCode: %v", errCode)
			util.ExecuteErrorResponse(ctx, http.StatusInternalServerError, "Something went wrong")
			return
		}
	}
	// cache product
	if h.redisClient != nil {
		productRaw, err := json.Marshal(product)
		if err != nil {
			logrus.Errorf("Error in Marshal product id %v!", inputGet.ProductID)
		} else {
			h.redisClient.Set(ctx, fmt.Sprintf("productId_%v", inputGet.ProductID), productRaw, 24*time.Hour)
			logrus.Infof("Save product id %v to redis cache", fmt.Sprintf("productId_%v", inputGet.ProductID))
		}
	}
	util.APIResponse(ctx, "OK", http.StatusOK, map[string]interface{}{
		"product": product,
	})
}
