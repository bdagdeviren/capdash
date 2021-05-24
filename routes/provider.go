package routes

import (
	"capdash/capi"
	"capdash/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddProviderRoute(database *db.Database,rg *gin.RouterGroup) {
	provider := rg.Group("/provider")

	provider.GET("/list", func (c *gin.Context) {
		result, err := capi.RunListProvider()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	})

	provider.POST("/init/:provider", func (c *gin.Context) {
		provider := c.Param("provider")
		err := capi.RunInitWithParameters(database, []string{provider})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message":"Successfully Init Provider"})
	})

	provider.POST("/delete/:provider", func (c *gin.Context) {
		provider := c.Param("provider")
		err := capi.RunDeleteWithParameters(database, []string{provider})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message":"Successfully Delete Provider"})
	})

}
