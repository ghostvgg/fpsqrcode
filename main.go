package main

import (
	"net/http"

	"fpsqrcode/qr"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	v1 := r.Group("/api/v1/fps")
	{
		v1.POST("/qr", func(c *gin.Context) {
			var req qr.QRRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			content := qr.GenerateQRString(req)
			c.JSON(http.StatusOK, gin.H{"qr_string": content})
		})
	}

	r.Run(":8080")
}
