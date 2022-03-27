package app

import (
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/realestate_agency-api/clients/elasticsearch"
	"github.com/superbkibbles/realestate_agency-api/constants"
	"github.com/superbkibbles/realestate_agency-api/http"
	cloudstorage "github.com/superbkibbles/realestate_agency-api/repository/cloudStorage"
	"github.com/superbkibbles/realestate_agency-api/repository/db"
	agencyservice "github.com/superbkibbles/realestate_agency-api/services/agencyService"
)

var (
	router  = gin.Default()
	handler http.AgencyHandler
)

func StartApplication() {
	elasticsearch.Client.Init()
	cld, err := cloudinary.NewFromParams(os.Getenv(constants.CLOUD_STORAGE_NAME), os.Getenv(constants.CLOUD_STORAGE_API_KEY), os.Getenv(constants.CLOUD_STORAGE_API_SECRET))
	if err != nil {
		panic(err)
	}
	handler = http.NewAgencyHandler(agencyservice.NewAgencyService(db.NewDbRepository(), cloudstorage.NewRepository(cld)))

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("lang")
	router.Use(cors.New(config))

	mapUrls()
	// router.Static("assets", "clients/visuals")
	router.Run(":3031")
}
