package product_controllers

type InputOption struct {
	OptionId    int64
	OptionValue string
}

type InputSearchProducts struct {
	Name          string  `json:"name" form:"name"`
	MinPrice      float64 `json:"min_price" form:"min_price" validate:"min=0"`
	MaxPrice      float64 `json:"max_price" form:"max_price" validate:"min=0"`
	Branch        string  `json:"branch" form:"branch"`
	Page          int     `json:"page" form:"page"`
	PageSize      int     `json:"page_size" form:"page_size"`
	SortBy        string  `json:"sort_by" form:"sort_by"`
	SortDirection string  `json:"sort_direction" form:"sort_direction"`
}

type InputGetProduct struct {
	ProductID int64 `json:"product_id" form:"product_id"`
}
