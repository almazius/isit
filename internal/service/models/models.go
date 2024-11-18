package models

type Material struct {
	Id            int     `json:"id" db:"id"`
	Name          string  `json:"name" db:"name" validate:"required"`
	Price         float64 `json:"price" db:"price"`
	Description   string  `json:"description" db:"description"`
	Address       string  `json:"address" db:"address"`
	RejectPercent float64 `json:"reject_percent" db:"reject_percent"`
	Count         float64 `json:"count" db:"count" validate:"required"`
}

type MaterialSmallInfo struct {
	Id    int     `json:"id"`
	Count float64 `json:"count"`
}

type Product struct {
	Id            int                 `json:"id" db:"id"`
	Name          string              `json:"name" db:"name"`
	Price         float64             `json:"price" db:"price"`
	Description   string              `json:"description" db:"description"`
	RejectPercent float64             `json:"reject_percent" db:"reject_percent"`
	Materials     []MaterialSmallInfo `json:"materials" db:"materials"`
}

type Order struct {
	ProductId int     `json:"product_id"`
	Count     float64 `json:"count"`
}

type UpdateOrder struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
}

type Orders struct {
	Id           int     `json:"id" db:"id"`
	Name         string  `json:"name" db:"name"`
	Description  string  `json:"description" db:"description"`
	Price        int     `json:"price " db:"price"`
	RejectPerson float64 `json:"reject_person" db:"reject_percent"`
	Materials    []int   `json:"materials" db:"materials"`
}

type GetOrder struct {
	Id        int     `json:"id" db:"id"`
	ProductId int     `json:"product_id" db:"product_id"`
	Count     float64 `json:"count" db:"count"`
	Status    string  `json:"status" db:"status"`
}
