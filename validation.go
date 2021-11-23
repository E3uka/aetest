package aetest

import (
	"math"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Validate the order request from user input.
func (req OrderRequest) Validate() error {
	return validation.ValidateStruct(
		&req,
		validation.Field(&req.Cart),
	)
}

// Validate the items in the order request from user input.
func (req Item) Validate() error {
	return validation.ValidateStruct(
		&req,
		// ItemName is a required field `NotNil` means it cannot be the nil
		// value for the type i.e. the empty string "".
		validation.Field(
			&req.ItemName,
			validation.NotNil,
		),
		// ItemName is a required field. Quantity cannot be 0.
		validation.Field(
			&req.Quantity,
			validation.Required,
			validation.Min(1),
			validation.Max(math.MaxInt),
		),
	)
}

// Validate the request to get a single order from user input.
func (req GetSingleOrderRequest) Validate() error {
	return validation.ValidateStruct(
		&req,
		validation.Field(
			&req.OrderID,
			validation.NotNil,
			is.UUIDv4,
		),
	)
}
