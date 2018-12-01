package models

type Order struct {
	Id       int64       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	BuyerId  int64       `gorm:"type:" json:"buyer_id"`
	SellerId int64       `gorm:"" json:"buyer_id"`
	Price    float32     `json:"price"`
	Quantity float64     `json:"quantity"`
	Status   OrderStatus `gorm:"type:int"`
}

type OrderStatus int

const (
	NEW         OrderStatus = 0
	WAIT_ACCEPT OrderStatus = 1
	ACCEPTED    OrderStatus = 2
	PAID        OrderStatus = 3
	UNPAID      OrderStatus = 4
	ACCOMPLISH  OrderStatus = 5
)
