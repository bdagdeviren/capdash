package routes

import (
	"capdash/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddManagementKubeconfigRoute(database *db.Database,rg *gin.RouterGroup) {
	kubeconfig := rg.Group("/kubeconfig")

	kubeconfig.GET("/get", func (c *gin.Context) {
		kubeConfig, err := database.GetManagementKubeconfig("management-kubeconfig")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"kubeconfig":kubeConfig })
	})

	kubeconfig.POST("/set", func (c *gin.Context) {
		file, err := c.FormFile("kubeconfig")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = database.SetManagementKubeconfig("management-kubeconfig",file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message":"Successfully Write Kubeconfig To Redis"})
	})

}

