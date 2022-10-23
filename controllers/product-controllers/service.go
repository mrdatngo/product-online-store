package product_controllers

import (
	"fmt"
	model "github.com/mrdatngo/gin-products/models"
	"time"
)

type Service interface {
	SearchProducts(input *InputSearchProducts) ([]*model.EntityProduct, string)
	GetProduct(input *InputGetProduct) (*model.EntityProduct, string)
}

type service struct {
	repository Repository
}

func NewProductService(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) SearchProducts(input *InputSearchProducts) ([]*model.EntityProduct, string) {
	products, errCode := s.repository.SearchProductsRepository(input)
	return products, errCode
}

func (s *service) GetProduct(input *InputGetProduct) (*model.EntityProduct, string) {
	product, errCode := s.repository.GetProductByIdRepository(input)
	return product, errCode
}

func (s *service) UpdateEmailService() {

	// query to search engines

	c1 := make(chan string)
	c2 := make(chan string)
	c3 := make(chan string)

	go func() {
		c1 <- SearchGoogle("")
	}()
	go func() {
		c2 <- SearchYahoo("")
	}()
	go func() {
		c3 <- SearchBing("")
	}()

	result := ""
	for i := 0; i < 3; i++ {
		select {
		case msg1 := <-c1:
			result += msg1 + "_"
		case msg2 := <-c2:
			result += msg2 + "_"
		case msg3 := <-c3:
			result += msg3 + "_"
		}
	}
	fmt.Println(result)
}

func SearchGoogle(msg string) string {
	time.Sleep(1 * time.Second)
	return "google"
}

func SearchYahoo(msg string) string {
	time.Sleep(1 * time.Second)
	return "yahoo"
}

func SearchBing(msg string) string {
	time.Sleep(2 * time.Second)
	return "bing"
}
