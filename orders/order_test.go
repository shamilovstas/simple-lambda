package orders

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrder_Record(t *testing.T) {
	order := Order{1, 2.5, "my_item"}
	expected := "Order ID: 1\nAmount: $2.50\nItem: my_item"
	actual := order.Record()
	assert.Equal(t, expected, actual)
}

func TestOrder_RecordAmountPrecision(t *testing.T) {
	order := Order{1, 0.6666666667, "my_item"}
	expected := "Order ID: 1\nAmount: $0.67\nItem: my_item"
	actual := order.Record()
	assert.Equal(t, expected, actual)
}
