package aetest

// OrderRequest are required values for an order submission.
type OrderRequest struct {
	Cart []Item `json:"cart"`
}

// GetSingleOrderRequest are required values for retrieving a single stored
// order. The OrderID must be of type uuid.
type GetSingleOrderRequest struct {
	OrderID string `json:"order_id"`
}

// NOTE: for Quantity, Cost and TotalCost `int` is used instead of usigned
// variant `uint`. Golang does not have a clean way of handling integer
// overflows for `uint`.

// Items are details regaring the name and quantity of items submitted for an
// order.
type Item struct {
	ItemName string `json:"item_name"`
	Quantity int    `json:"quantity"`
}

// ItemsWithCost are `Items` with the items respective cost included.
type ItemWithCost struct {
	ItemName string `json:"item_name"`
	Quantity int    `json:"quantity"`
	Cost     int    `json:"cost"`
}

// Summary is the response to the call to the orders API.
type OrderSummary struct {
	OrderID   string         `json:"order_id"`
	Summary   []ItemWithCost `json:"summary"`
	TotalCost int            `json:"total_cost"`
}

// AllOrders is the response to the call to get all stored orders.
type AllOrders struct {
	// omitempty structtag used to return an empty object if no order
	// previously exists.
	Orders []OrderSummary `json:"orders,omitempty"`
}

// GenericErrResponse is a generic error result return to the caller after an
// error is raised from an endpoint. The appropriate error reason should be
// returned to the caller.
type GenericErrResponse struct {
	Err string `json:"error,omitempty"`
}
