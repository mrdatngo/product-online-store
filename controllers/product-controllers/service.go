package product_controllers

import (
	model "github.com/mrdatngo/gin-products/models"
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
