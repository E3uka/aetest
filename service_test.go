package aetest

import (
	"bytes"
	"context"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

var (
	item_store, discount, order_store = NewStore()
	service                           Service
	router                            http.Handler
	goodOrderRequest                  OrderRequest
)

// Scaffold required globals items for use in the test cases.
func TestMain(m *testing.M) {
	service = New(item_store, discount, order_store)
	router = NewOrdersRouter(service)

	// Set up a base struct of a good order request to use and manipulate in
	// the below test cases.
	goodOrderRequest = OrderRequest{
		Cart: []Item{
			{ItemName: "Apples", Quantity: 2},
			{ItemName: "Oranges", Quantity: 3},
		},
	}

	// Run all tests and pass the exit code to os.Exit
	code := m.Run()
	os.Exit(code)
}

func TestWellFormedOrderRequestHTTP(t *testing.T) {
	ctx := context.Background()
	JSON, err := json.Marshal(goodOrderRequest)
	require.NoError(t, err)
	reader := bytes.NewReader(JSON)

	request := httptest.NewRequest("POST", "/submit-order", reader)
	request = request.WithContext(ctx)
	request.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, request)

	response := rec.Result()

	// The successful request will return a valid OrderSummary. There must be
	// no issues deserializing the response.
	var success OrderSummary

	err = json.NewDecoder(response.Body).Decode(&success)
	require.Nil(t, err, "error deserializing response")

	// Check that the total for this order equals (2 x Apples) + (3 Oranges)
	// subtracting the discount applied for these items i.e buy one get on free
	// for apples and 3 for the price of two on oranges.
	require.Equal(t, 110, success.TotalCost, "incorrect order total")
}

func TestMalformedOrderRequestHTTP(t *testing.T) {
	badRequestEmptyCart := []Item{}
	badRequestEmptyItemName := []Item{
		{ItemName: "Apples", Quantity: 1},
		{ItemName: "", Quantity: 27},
	}
	badRequestItemNotFound := []Item{
		{ItemName: "Magazine", Quantity: 2},
		{ItemName: "Apples", Quantity: 45},
	}
	badRequestNegativeQuantitiy := []Item{
		{ItemName: "Apples", Quantity: -2},
		{ItemName: "Apples", Quantity: 4},
	}
	badRequestCannotProcessPrice := []Item{
		{ItemName: "Apples", Quantity: math.MaxInt},
	}

	// Table driven test with created malformed requests. These must all return
	// an error.
	testCases := []struct {
		name     string
		testCase OrderRequest
	}{
		{"empty cart", OrderRequest{badRequestEmptyCart}},
		{"empty item name", OrderRequest{badRequestEmptyItemName}},
		{"item not found", OrderRequest{badRequestItemNotFound}},
		{"negative quantity requested", OrderRequest{badRequestNegativeQuantitiy}},
		{"cannot process price", OrderRequest{badRequestCannotProcessPrice}},
	}

	// Iterate through testcases and perform the request
	for _, tc := range testCases {
		ctx := context.Background()
		JSON, err := json.Marshal(tc.testCase)
		require.NoError(t, err)
		reader := bytes.NewReader(JSON)

		request := httptest.NewRequest("POST", "/submit-order", reader)
		request = request.WithContext(ctx)
		request.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		response := rec.Result()

		// The malformed requests will return a GenericErrResponse. There must
		// be no issues deserializing the response.
		var failure GenericErrResponse

		err = json.NewDecoder(response.Body).Decode(&failure)
		require.Nilf(t, err, "error in case: %v", tc.name)
	}
}

func TestGetSingleOrderRequest(t *testing.T) {
	// Strategy:
	// 1.	Make an order request and get as response an OrderSummary.
	// 2.	Use the generated OrderID from the OrderSummary to make a GET
	//		request to the `/get-order` endpoint.
	// 3.	Check whether inital order summary is equal to the one received
	//		from get single order request this signifies that the order was
	//		stored correctly.

	// 1.
	ctx := context.Background()
	JSON, err := json.Marshal(goodOrderRequest)
	require.NoError(t, err)
	reader := bytes.NewReader(JSON)
	request := httptest.NewRequest("POST", "/submit-order", reader)
	request = request.WithContext(ctx)
	request.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, request)
	response := rec.Result()
	var success OrderSummary
	err = json.NewDecoder(response.Body).Decode(&success)
	require.Nil(t, err, "error deserializing response")

	// 2.
	//ctx = context.Background() // new context used for
	single_order := GetSingleOrderRequest{success.OrderID}
	JSON, err = json.Marshal(single_order)
	require.NoError(t, err)
	reader = bytes.NewReader(JSON)
	request = httptest.NewRequest("POST", "/get-order", reader)
	request = request.WithContext(ctx)
	request.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, request)
	response = rec.Result()
	var newResponse OrderSummary
	err = json.NewDecoder(response.Body).Decode(&newResponse)
	require.Nil(t, err, "error deserializing response")

	// 3.
	require.Equal(t, newResponse, success, "order summaries do not match")
}

func TestOrderNotFoundRequest(t *testing.T) {
	// Create an order id that does not exist in the store.
	random_id := uuid.NewV4().String()
	strange_order := GetSingleOrderRequest{random_id}

	ctx := context.Background()
	JSON, err := json.Marshal(strange_order)
	require.NoError(t, err)
	reader := bytes.NewReader(JSON)
	request := httptest.NewRequest("POST", "/get-order", reader)
	request = request.WithContext(ctx)
	request.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, request)
	response := rec.Result()

	// The alien get order requests will return fail and return a
	// GenericErrResponse. There must be no issues deserializing the response.
	var newResponse GenericErrResponse

	err = json.NewDecoder(response.Body).Decode(&newResponse)
	require.Nil(t, err, "error deserializing the response")
}

func TestInvalidGetOrderRequest(t *testing.T) {
	// Create an order id that does not exist in the store.
	invalid_id := "this is not an id"
	strange_order := GetSingleOrderRequest{invalid_id}

	ctx := context.Background()
	JSON, err := json.Marshal(strange_order)
	require.NoError(t, err)
	reader := bytes.NewReader(JSON)
	request := httptest.NewRequest("POST", "/get-order", reader)
	request = request.WithContext(ctx)
	request.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, request)
	response := rec.Result()

	// The alien get order requests will return fail and return a
	// GenericErrResponse. There must be no issues deserializing the response.
	var newResponse GenericErrResponse

	err = json.NewDecoder(response.Body).Decode(&newResponse)
	require.Nil(t, err, "error deserializing the response")
}

func TestGetAllOrdersRequest(t *testing.T) {
	ctx := context.Background()
	reader := bytes.NewReader([]byte{})
	request := httptest.NewRequest("GET", "/get-all-orders", reader)
	request = request.WithContext(ctx)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, request)
	response := rec.Result()

	var all_orders AllOrders

	err := json.NewDecoder(response.Body).Decode(&all_orders)
	require.Nil(t, err, "error deserializing the response")
	require.NotEmpty(t, all_orders.Orders, "store should contain items")
}

func TestEmptyGetAllOrdersRequest(t *testing.T) {
	// Create empty order store and use that to create new service and router.
	empty_store := make(OrderStore)
	service_with_empty_store := New(item_store, discount, empty_store)
	router_with_empty_store := NewOrdersRouter(service_with_empty_store)

	ctx := context.Background()
	reader := bytes.NewReader([]byte{})
	request := httptest.NewRequest("GET", "/get-all-orders", reader)
	request = request.WithContext(ctx)
	rec := httptest.NewRecorder()
	router_with_empty_store.ServeHTTP(rec, request)
	response := rec.Result()

	var all_orders AllOrders

	err := json.NewDecoder(response.Body).Decode(&all_orders)
	require.Nil(t, err, "error deserializing the response")

	// Verify that the store is indeed empty.
	require.Empty(t, all_orders.Orders, "store should be empty")
}
