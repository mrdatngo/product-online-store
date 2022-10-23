package model

import (
	"gorm.io/gorm"
	"time"
)

// EntityProduct product model
type EntityProduct struct {
	ID          int64               `json:"id" gorm:"type:uint;primaryKey;autoIncrement"`
	Name        string              `json:"name" gorm:"type:varchar(255);not null"`
	BranchID    int64               `json:"branch_id" gorm:"type:uint;not null"`
	Active      bool                `json:"active" gorm:"type:bool;default:true"`
	MetaData    string              `json:"meta_data" gorm:"type:varchar(255)"`
	Price       float64             `json:"price"`
	ImgUrl      string              `json:"img_url" gorm:"type:varchar(255)"`
	CreatedAt   time.Time           `json:"-"`
	UpdatedAt   time.Time           `json:"-"`
	Branch      *EntityBranch       `json:"branch" gorm:"foreignKey:BranchID"`
	ProductSkus []*EntityProductSku `json:"product_skus" gorm:"foreignKey:ProductID"`
}

func (entity *EntityProduct) BeforeCreate(db *gorm.DB) error {
	entity.CreatedAt = time.Now().Local()
	return nil
}

func (entity *EntityProduct) BeforeUpdate(db *gorm.DB) error {
	entity.UpdatedAt = time.Now().Local()
	return nil
}

func (EntityProduct) TableName() string {
	return "product"
}

// EntityOption option model
type EntityOption struct {
	ID        int64          `json:"id" gorm:"type:uint;primaryKey;autoIncrement"`
	ProductID int64          `json:"product_id" gorm:"type:uint;not null"`
	Name      string         `json:"name" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	Product   *EntityProduct `json:"product" gorm:"foreignKey:ProductID"`
}

func (entity *EntityOption) BeforeCreate(db *gorm.DB) error {
	entity.CreatedAt = time.Now().Local()
	return nil
}

func (entity *EntityOption) BeforeUpdate(db *gorm.DB) error {
	entity.UpdatedAt = time.Now().Local()
	return nil
}

func (EntityOption) TableName() string {
	return "option"
}

// EntityOptionValue OptionValue model
type EntityOptionValue struct {
	ID        int64          `json:"id" gorm:"type:uint;primaryKey;autoIncrement"`
	ProductID int64          `json:"product_id" gorm:"type:uint;not null"`
	OptionID  int64          `json:"option_id" gorm:"type:uint;not null"`
	ValueName string         `json:"value_name" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	Product   *EntityProduct `json:"product" gorm:"foreignKey:ProductID"`
	Option    *EntityOption  `json:"option" gorm:"foreignKey:OptionID"`
}

func (entity *EntityOptionValue) BeforeCreate(db *gorm.DB) error {
	entity.CreatedAt = time.Now().Local()
	return nil
}

func (entity *EntityOptionValue) BeforeUpdate(db *gorm.DB) error {
	entity.UpdatedAt = time.Now().Local()
	return nil
}

func (EntityOptionValue) TableName() string {
	return "option_value"
}

// EntityProductSku product_sku models
type EntityProductSku struct {
	SkuID     int64             `json:"sku_id" gorm:"type:uint;primaryKey;autoIncrement"`
	ProductID int64             `json:"product_id" gorm:"type:uint;not null"`
	Sku       string            `json:"sku" gorm:"type:varchar(255)"`
	Price     float64           `json:"price"`
	CreatedAt time.Time         `json:"-"`
	UpdatedAt time.Time         `json:"-"`
	Product   *EntityProduct    `json:"product" gorm:"foreignKey:ProductID"`
	SkuValues []*EntitySkuValue `json:"sku_values" gorm:"foreignKey:ProductID,SkuID;references:ProductID,SkuID"`
}

func (entity *EntityProductSku) BeforeCreate(db *gorm.DB) error {
	entity.CreatedAt = time.Now().Local()
	return nil
}

func (entity *EntityProductSku) BeforeUpdate(db *gorm.DB) error {
	entity.UpdatedAt = time.Now().Local()
	return nil
}

func (EntityProductSku) TableName() string {
	return "product_sku"
}

// EntitySkuValue product_sku models
type EntitySkuValue struct {
	ProductID   int64              `json:"product_id" gorm:"primaryKey;type:uint;not null"`
	SkuID       int64              `json:"sku_id" gorm:"primaryKey;type:uint;not null"`
	OptionID    int64              `json:"option_id" gorm:"primaryKey;type:uint;not null"`
	ValueID     int64              `json:"value_id" gorm:"primaryKey;type:uint;not null"`
	CreatedAt   time.Time          `json:"-"`
	UpdatedAt   time.Time          `json:"-"`
	Product     *EntityProduct     `json:"product" gorm:"foreignKey:ProductID"`
	ProductSku  *EntityProductSku  `json:"product_sku" gorm:"foreignKey:SkuID"`
	Option      *EntityOption      `json:"option" gorm:"foreignKey:OptionID"`
	OptionValue *EntityOptionValue `json:"option_value" gorm:"foreignKey:ValueID"`
}

func (entity *EntitySkuValue) BeforeCreate(db *gorm.DB) error {
	entity.CreatedAt = time.Now().Local()
	return nil
}

func (entity *EntitySkuValue) BeforeUpdate(db *gorm.DB) error {
	entity.UpdatedAt = time.Now().Local()
	return nil
}

// TableName return table name type string to define sql table
func (EntitySkuValue) TableName() string {
	return "sku_values"
}
