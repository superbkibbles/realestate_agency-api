package app

func mapUrls() {
	router.GET("/api/agency", handler.Get)     // Get all agencies
	router.POST("/api/agency", handler.Create) // Create new Agency
	router.GET("/api/agency/:agency_id", handler.GetByID)
}
