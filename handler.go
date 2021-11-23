package aetest

import (
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

		// Submit a simple request to the `Service`, if successful this will
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

	return router
}
