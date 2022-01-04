package app

func mapUrls() {
	router.GET("/api/agency", handler.Get)                    // Get all agencies
	router.POST("/api/agency", handler.Create)                // Create new Agency
	router.GET("/api/agency/:agency_id", handler.GetByID)     // Get Agency By ID
	router.POST("/api/agency/:agency_id", handler.UploadIcon) // Upload the Icon
	router.PATCH("/api/agency/:agency_id", handler.Update)    // Update Agency
	router.POST("/api/agency/search/s", handler.Search)       // Search for properties
	router.DELETE("/api/agency/:agency_id", handler.DeleteIcon)
}
