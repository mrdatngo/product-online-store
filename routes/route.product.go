package route

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/mrdatngo/gin-products/connectors"
	product_controllers "github.com/mrdatngo/gin-products/controllers/product-controllers"
	product_handlers "github.com/mrdatngo/gin-products/handlers"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
	"strconv"
)

func InitProductRoutes(db *gorm.DB, route *gin.Engine) {

	/**
	@description All Handler Auth
	*/
	productRepository := product_controllers.NewProductRepository(db)
	productService := product_controllers.NewProductService(productRepository)

	kafkaProducer := &connectors.KafkaProducer{}
	err := kafkaProducer.Init()
	if err != nil {
		logrus.Errorf("Init kafka failed!")
	}
	dbStr := os.Getenv("REDIS_DB")
	dbI, err := strconv.Atoi(dbStr)
	if err != nil {
		logrus.Errorf("Error in get DB variable for redis, use default 0, err: %v", err)
		dbI = 0
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       dbI,                         // use default DB
	})
	productHandler := product_handlers.NewProductHandler(productService, kafkaProducer, redisClient)

	/**
	@description All Auth Route
	*/
	groupRoute := route.Group("/api/v1/product")
	groupRoute.GET("/search", productHandler.SearchProductsHandler)
	groupRoute.GET("/:product_id", productHandler.GetProductHandler)
}
