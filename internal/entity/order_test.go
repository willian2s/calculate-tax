package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIfGetAnErrorIfIdIsBlank(t *testing.T) {
	order := Order{}
	assert.Error(t, order.Validate(), "invalid id")
}

func TestIfGetAnErrorIfPriceIsBlank(t *testing.T) {
	order := Order{ID: "123"}
	assert.Error(t, order.Validate(), "invalid id")
}

func TestIfGetAnErrorIfTaxIsBlank(t *testing.T) {
	order := Order{ID: "123", Price: 5.0}
	assert.Error(t, order.Validate(), "invalid price")
}

func TestIfAllValidParams(t *testing.T) {
	order := Order{ID: "123", Price: 5.0, Tax: 0.5}
	assert.NoError(t, order.Validate())
	assert.Equal(t, 5.0, order.Price)
	assert.Equal(t, 0.5, order.Tax)
	order.CalculateFinalPrice()
	assert.Equal(t, 5.5, order.FinalPrice)

}
