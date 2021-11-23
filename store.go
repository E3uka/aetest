package aetest

// ItemStore is a `map[string]int` that stores as a key the item name with a
// value of the cost of the item.
type ItemStore map[string]int

// Discount is a function that takes as input the item cost and the quantity
// that is submitted for order and returns the discount to subtract from the
// order total.
type DiscountFunction func(cost int, quantity int) int

// ItemDiscount is a `map[string]DiscountFunction that stores as key the item
// name with a value of a function that calculates the discount of the item.
// This is used to lookup the item and apply a relevant discount to it.
type ItemDiscount map[string]DiscountFunction

// NewStore creates an `ItemStore` and populates the key and values of the
// store with required item names and costs respectively.
func NewStore() (ItemStore, ItemDiscount) {
	store := make(ItemStore)
	store["Apples"] = 60
	store["Oranges"] = 25

	// Apples are buy one get one free
	var applesDiscount DiscountFunction = func(cost int, quantity int) int {
		if quantity == 1 {
			// handle base case: only 1 item, no discount applied
			return 0
		} else if quantity%2 == 0 {
			// half the items are now discounted
			return cost * (quantity / 2)
		} else {
			// 1 less than half the items are now discounted
			return cost * ((quantity - 1) / 2)
		}
	}

	// Oranges are 3 for the price of two
	var orangesDiscount DiscountFunction = func(cost int, quantity int) int {
		if quantity < 3 {
			// handle base case: less than 3 items, no discount applied
			return 0
		} else if quantity%3 == 0 {
			return cost * (quantity / 3)
		} else if quantity%3 == 1 {
			return cost * ((quantity - 1) / 3)
		} else {
			return cost * ((quantity - 2) / 3)
		}
	}

	// Add discount to discount lookup
	discount := make(ItemDiscount)
	discount["Apples"] = applesDiscount
	discount["Oranges"] = orangesDiscount

	return store, discount
}
