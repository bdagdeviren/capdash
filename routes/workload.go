package routes

import (
	"capdash/capi"
	"capdash/db"
	"context"
	b64 "encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddWorkloadKubeconfigRoute(database *db.Database,rg *gin.RouterGroup) {
	kubeconfig := rg.Group("/kubeconfig")

	kubeconfig.GET("/:cluster", func (c *gin.Context) {
		cluster := c.Param("cluster")
		kubeConfig := ""
		result,err := database.Client.Get(context.Background(),cluster).Result()
		if err.Error() == "redis: nil" {
			kubeConfig, err = capi.RunGetWorkloadClusterKubeconfigWithParameters(database, cluster)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}else {
			encBase64, err := b64.StdEncoding.DecodeString(result)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			kubeConfig = string(encBase64)
		}
		c.JSON(http.StatusOK, gin.H{"kubeconfig":kubeConfig })
	})

}

func AddWorkloadClusterRoute(database *db.Database,rg *gin.RouterGroup) {
	cluster := rg.Group("/cluster")

	cluster.GET("/yaml/:provider/:cluster", func (c *gin.Context) {
		cluster := c.Param("cluster")
		//provider := c.Param("provider")
		template,err := database.Client.Get(context.Background(),cluster).Result()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot find cluster template!"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"template":template })
	})

	cluster.GET("/create/:provider/:cluster", func (c *gin.Context) {
		cluster := c.Param("cluster")
		//provider := c.Param("provider")

		err := capi.CreateCluster(database, cluster, "", "", "v1.20.0", 1, 2, "./hack/cluster-template.yaml")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message":"Successfully Create Cluster"})
	})

	cluster.GET("/delete/:provider/:cluster", func (c *gin.Context) {
		cluster := c.Param("cluster")
		//provider := c.Param("provider")

		err := capi.DeleteCluster(database, cluster)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message":"Successfully Delete Cluster"})
	})

}
