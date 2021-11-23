package aetest

// OrderRequest are required values for an order submission.
type OrderRequest struct {
	Cart []Item `json:"cart"`
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
	Summary   []ItemWithCost `json:"summary"`
	TotalCost int            `json:"total_cost"`
}

// GenericErrResponse is a generic error result return to the caller after an
// error is raised from an endpoint. The appropriate error reason should be
// returned to the caller.
type GenericErrResponse struct {
	Err string `json:"error,omitempty"`
}
