package product_controllers

import (
	_const "github.com/mrdatngo/gin-products/const"
	"github.com/mrdatngo/gin-products/controllers"
	model "github.com/mrdatngo/gin-products/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Repository interface {
	SearchProductsRepository(input *InputSearchProducts) ([]*model.EntityProduct, string)
	GetProductByIdRepository(input *InputGetProduct) (*model.EntityProduct, string)
}

type repository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) SearchProductsRepository(input *InputSearchProducts) (products []*model.EntityProduct, errCode string) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error("Recovered in SearchProductsRepository", r)
			errCode = _const.UNKNOWN_PANIC
		}
	}()
	if input.SortBy == "" {
		input.SortBy = "updated_at"
		input.SortDirection = "desc"
	}

	db := r.db.Model(&products)
	tx := db.Select("id, name, branch_id, price, img_url, meta_data, active")
	tx.Scopes(controllers.Paginate(input.Page, input.PageSize))
	tx.Order(input.SortBy + " " + input.SortDirection)
	if len(input.Name) > 0 {
		tx.Where("name LIKE ?", `%${input.Name}%`)
	}
	if len(input.Branch) > 0 {
		tx.Where("branch_id = ?", input.Branch)
	}
	if input.MinPrice > 0 {
		tx.Where("price >= ?", input.MinPrice)
	}
	if input.MaxPrice > 0 {
		tx.Where("price <= ?", input.MaxPrice)
	}
	tx.Preload("Branch")
	err := tx.Find(&products).Error
	if err != nil {
		logrus.Errorf("Error: %v", err)
		return nil, err.Error()
	}
	return products, ""
}

func (r *repository) GetProductByIdRepository(inputGet *InputGetProduct) (*model.EntityProduct, string) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error("Recovered in GetProductByIdRepository", r)
		}
	}()
	var product *model.EntityProduct
	db := r.db.Model(&product)
	selectStr := "id, name, branch_id, product.price, img_url, meta_data, active"
	tx := db.Select(selectStr).Where("product.id = ?", inputGet.ProductID)
	tx.Joins("JOIN product_sku as ps on ps.product_id = product.id")
	tx.Preload("ProductSkus")
	tx.Preload("ProductSkus.SkuValues")
	tx.Preload("ProductSkus.SkuValues.Option")
	tx.Preload("ProductSkus.SkuValues.OptionValue")
	tx.Preload("Branch")
	tx.Group(selectStr)

	err := tx.First(&product).Error
	if err != nil {
		logrus.Errorf("Error: %v", err)
		if err.Error() == _const.RECORD_NOT_FOUND {
			return nil, _const.PRODUCT_NOT_FOUND
		} else {
			return nil, err.Error()
		}
	}
	return product, ""
}
