package aetest

// ItemStore is a `map[string]int` that stores as a key the item name with a
// value of the cost of the item.
type ItemStore map[string]int

// NewStore creates an `ItemStore` and populates the key and values of the
// store with required item names and costs respectively.
func NewStore() ItemStore {
	store := make(ItemStore)
	store["Apples"] = 60
	store["Oranges"] = 25

	return store
}
