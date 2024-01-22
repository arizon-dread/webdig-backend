package main

import (
	"fmt"
	"net/http"

	businesslayer "github.com/arizon-dread/webdig-backend/businessLayer"
	"github.com/arizon-dread/webdig-backend/models"
	"github.com/gin-gonic/gin"
)

func lookup(c *gin.Context) {
	var req models.Req
	c.BindJSON(&req)

	var status int = http.StatusOK
	resp, err := businesslayer.Lookup(req)
	if err != nil {
		fmt.Printf("Unable to lookup IP, %v", err)
		status = http.StatusBadRequest
	}
	c.JSON(status, resp)

}
