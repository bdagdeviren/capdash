package main

import (
	"capdash/db"
	"capdash/routes"
	"github.com/gin-gonic/gin"
	"log"
)

var (
	ListenAddr = "localhost:8080"
	RedisAddr = "192.168.1.2:6379"
)
func main() {
	database, err := db.NewDatabase(RedisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %s", err.Error())
	}
	router := initRouter(database)
	err = router.Run(ListenAddr)
	if err != nil {
		return 
	}
}

func initRouter(database *db.Database) *gin.Engine {
	r := gin.Default()
	route := r.Group("api/v1")
	subRouteManagement := route.Group("/management")
	routes.AddManagementKubeconfigRoute(database,subRouteManagement)
	routes.AddProviderRoute(database,subRouteManagement)
	routes.AddClusterRoute(database,subRouteManagement)
	subRouteWorkload := route.Group("workload")
	routes.AddWorkloadKubeconfigRoute(database,subRouteWorkload)
	routes.AddWorkloadClusterRoute(database,subRouteWorkload)
	return r
}

