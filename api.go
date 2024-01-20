package main

import (
	"fmt"
	"net/http"
	"unicode"

	businesslayer "github.com/arizon-dread/webdig-backend/businessLayer"
	"github.com/arizon-dread/webdig-backend/models"
	"github.com/gin-gonic/gin"
)

func lookup(c *gin.Context) {
	var req models.Req
	c.BindJSON(&req)
	isDNS := func() bool {
		for _, r := range req.Host {
			if unicode.IsLetter(r) {
				return true
			}
		}
		return false
	}
	var status int = http.StatusOK
	if isDNS() {
		resp, err := businesslayer.LookupIP(req)
		if err != nil {
			fmt.Printf("Unable to lookup IP, %v", err)
			status = http.StatusBadRequest
		}
		c.JSON(status, resp)

	} else {
		resp, err := businesslayer.LookupDNS(req)
		if err != nil {
			fmt.Printf("Unable to lookup dns, %v", err)
			status = http.StatusBadRequest
		}
		c.JSON(status, resp)
	}
}
