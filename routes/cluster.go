package routes

import (
	"capdash/capi"
	"capdash/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddClusterRoute(database *db.Database,rg *gin.RouterGroup) {
	kubeconfig := rg.Group("/cluster")

	kubeconfig.GET("/list", func (c *gin.Context) {
		clusterList, err := capi.GetClusters("management-kubeconfig")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK,clusterList)
	})

	//kubeconfig.POST("/set", func (c *gin.Context) {
	//	file, err := c.FormFile("kubeconfig")
	//	if err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//		return
	//	}
	//	err = database.SetManagementKubeconfig("management-kubeconfig",file)
	//	if err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//		return
	//	}
	//	c.JSON(http.StatusOK, gin.H{"message":"Successfully Write Kubeconfig To Redis"})
	//})

}
