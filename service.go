package aetest

import (
	"errors"

	"github.com/johncgriffin/overflow"
)

var (
	// ErrInvalidRequest is returned when the user input request is invalid.
	ErrInvalidRequest = errors.New("invalid submitted request")

	// ErrItemDoesntExist is returned when the user has submitted an invalid
	// item to the orders service.
	ErrItemDoesNotExist = errors.New("one or more items in the request does not exist")

	// ErrIntegerOverflow is returned when the user enters an integer value
	// that the hardware cannot process without overlow.
	ErrIntegerOverflow = errors.New("unable to process order request, item total too large")
)

// Service is an interface that encapsulates all the functionalities of the
// orders service.
type Service interface {
	// SimpleSummary creates an OrderSummary from a submitted order request. If
	// the order is invalid this returns an empty OrderSummary and a relevant
	// error message to the caller.
	SimpleSummary(req OrderRequest) (OrderSummary, error)
}

// orderService is a private struct that is used to satisfy the interface
// requirements of the Service. The methods of this structure is used to call
// the Service' methods. This struct holds an ItemStore that is used to provide
// a lookup of the cost of the users items.
type orderService struct {
	store    ItemStore
	discount ItemDiscount
}

// InjectCost adds the cost the user supplied Cart. This makes use of the
// orderService' internal ItemStore map to lookup the items cost. If an Item
// does not exist in the internal ItemStore this returns an empty
// `ItemsWithCost` and `false`.
func (svc orderService) InjectCost(cart []Item) ([]ItemWithCost, bool) {
	injectedItems := []ItemWithCost{}

	for _, item := range cart {
		cost, ok := svc.store[item.ItemName]
		if !ok {
			// item does not exist
			return []ItemWithCost{}, false
		}
		with_cost := ItemWithCost{item.ItemName, item.Quantity, cost}
		injectedItems = append(injectedItems, with_cost)
	}

	return injectedItems, true
}

// New returns a new Service to the caller.
func New(store ItemStore, discount ItemDiscount) Service {
	return orderService{store, discount}
}

func (svc orderService) SimpleSummary(
	req OrderRequest,
) (OrderSummary, error) {
	// Validate the user input using custom validation schema.
	if err := req.Validate(); err != nil {
		return OrderSummary{}, ErrInvalidRequest
	}

	// Inject associated costs of the items to the cart using a price lookup.
	cart_with_costs, ok := svc.InjectCost(req.Cart)
	if !ok {
		return OrderSummary{}, ErrItemDoesNotExist
	}

	var running_total int = 0

	// Iterate through the items in the cart adding the calculated amount to
	// the running total. Integer overflows need to be handled appropriately,
	// this is detected initially on the multiplication of the item Cost x
	// Quantity and finally during the sum of the result and running total.
	for _, item := range cart_with_costs {
		intermediate_result, ok := overflow.Mul(item.Cost, item.Quantity)
		if !ok {
			return OrderSummary{}, ErrIntegerOverflow
		}

		result, ok := overflow.Add(intermediate_result, running_total)
		if !ok {
			return OrderSummary{}, ErrIntegerOverflow
		}

		// Using the ItemDiscount lookup whether a discount exists for that
		// item. If a discount is not found by the above logic this does not
		// mean that the item does not exist in the store, it is fine to skip
		// the discount step. If the discount exists apply the discount to the
		// result.
		calculate_discount, ok := svc.discount[item.ItemName]
		if ok {
			// discount found, apply the discount
			discount := calculate_discount(item.Cost, item.Quantity)
			result -= discount
		}

		running_total = result
	}

	return OrderSummary{cart_with_costs, running_total}, nil
}
