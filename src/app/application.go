package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/realestate_agency-api/src/clients/elasticsearch"
	"github.com/superbkibbles/realestate_agency-api/src/http"
	"github.com/superbkibbles/realestate_agency-api/src/repository/db"
	agencyservice "github.com/superbkibbles/realestate_agency-api/src/services/agencyService"
)

var (
	router  = gin.Default()
	handler http.AgencyHandler
)

func StartApplication() {
	elasticsearch.Client.Init()
	handler = http.NewAgencyHandler(agencyservice.NewAgencyService(db.NewDbRepository()))
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))
	mapUrls()
	// router.Use(cors.Default())
	router.Static("assets", "clients/visuals")
	router.Run(":3031")
}
