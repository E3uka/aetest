package aetest

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewOrdersRouter(svc Service) http.Handler {
	router := gin.New()

	// Ignoring extra router options i.e. cors, timeouts, allowed methods etc.
	// for simplicity.
	router.Use(gin.Recovery())

	router.POST("/submit-order", func(c *gin.Context) {
		var request OrderRequest

		// Deserialize JSON POST request into the OrderRequest struct, if
		// serialization fails return a `GenericErrResponse` to the caller with
		// the appropriate status code for bad request.
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, GenericErrResponse{
				Err: err.Error(),
			})
			return
		}

		// Submit an order request to the `Service`, if successful this will
		// return an OrderSummary and a nil error. If an error has occurred,
		// this returns and empty OrderSummary and an error.
		response, err := svc.SimpleSummary(request)
		if err != nil {
			c.JSON(http.StatusBadRequest, GenericErrResponse{
				Err: err.Error(),
			})
			return
		}

		// Serialize response as JSON and return to caller with a `Ok` status.
		c.JSON(http.StatusOK, response)
	})

	router.POST("/get-order", func(c *gin.Context) {
		var request GetSingleOrderRequest

		// Deserialize JSON POST request into the GetSingleOrderRequest struct,
		// if serialization fails return a `GenericErrResponse` to the caller
		// with the appropriate status code for bad request.
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, GenericErrResponse{
				Err: err.Error(),
			})
			return
		}

		// Submit a get single order to the `Service`, if successful this will
		// return an OrderSummary and a nil error. If an error has occurred,
		// this returns and empty OrderSummary and an error.
		response, err := svc.GetSingleOrder(request)
		if err != nil {
			// Malformed request, respond with 400
			if ok := errors.Is(err, ErrInvalidRequest); ok {
				c.JSON(http.StatusBadRequest, GenericErrResponse{
					Err: err.Error(),
				})
				return
			}

			// Order not found, respond with 404
			c.JSON(http.StatusNotFound, GenericErrResponse{
				Err: err.Error(),
			})
			return
		}

		// Serialize response as JSON and return to caller with a `Ok` status.
		c.JSON(http.StatusOK, response)
	})

	router.GET("get-all-orders", func(c *gin.Context) {
		// Submit a get all orders request to the `Service`. This will return
		// a GetAllOrders object. If no orders exist in the OrderStore, this
		// returns an Okay status with an empty response.
		orders := svc.GetAllOrders()
		c.JSON(http.StatusOK, orders)
	})

	return router
}
