package orders

import "fmt"

type Order struct {
	OrderId int     `json:"order_id"`
	Amount  float64 `json:"amount"`
	Item    string  `json:"item"`
}

func (o Order) Record() string {
	return fmt.Sprintf("Order ID: %d\nAmount: $%.2f\nItem: %s", o.OrderId, o.Amount, o.Item)
}
