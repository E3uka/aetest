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

	"github.com/stretchr/testify/require"
)

var (
	store, discount  = NewStore()
	service          Service
	router           http.Handler
	goodOrderRequest OrderRequest
)

// Scaffold required globals items for use in the test cases.
func TestMain(m *testing.M) {
	service = New(store, discount)
	router = NewOrdersRouter(service)

	// Set up a base struct of a good order request to use and manipulate in the
	// below test cases.
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

func TestWellFormedRequestHTTP(t *testing.T) {
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

func TestMalformedRequestHTTP(t *testing.T) {
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
