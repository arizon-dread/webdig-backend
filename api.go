package main

import (
	"net/http"
	"strings"

	businesslayer "github.com/arizon-dread/webdig-backend/businessLayer"
	"github.com/arizon-dread/webdig-backend/models"
	"github.com/gin-gonic/gin"
)

func lookup(c *gin.Context) {
	var req models.Req
	c.BindJSON(&req)

	var status int = http.StatusOK
	resp, err := businesslayer.Lookup(c.Request.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "could not find dns record") {
			status = http.StatusNotFound
		} else {
			status = http.StatusBadRequest
		}

	}
	c.JSON(status, resp)

}
func healthz(c *gin.Context) {
	c.JSON(http.StatusOK, "Healthy")
}
