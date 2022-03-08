package app

const (
	prefix = "/api/agency"
)

func mapUrls() {
	router.GET(prefix, handler.Get)                                 // Get all agencies
	router.POST(prefix, handler.Create)                             // Create new Agency
	router.GET(prefix+"/:agency_id", handler.GetByID)               // Get Agency By ID
	router.POST(prefix+"/:agency_id", handler.UploadIcon)           // Upload the Icon
	router.PATCH(prefix+"/:agency_id", handler.Update)              // Update Agency
	router.POST(prefix+"/search/s", handler.Search)                 // Search for properties
	router.DELETE(prefix+"/:agency_id", handler.DeleteIcon)         // Delete ICON
	router.PATCH(prefix+"/:agency_id/translate", handler.Translate) // Translate
}
