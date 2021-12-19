package agency

const (
	STATUS_ACTIVE   = "active"
	STATUS_DEACTIVE = "deactive"
)

type Agency struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Icon           string `json:"icon"`
	Email          string `json:"email"`
	Country        string `json:"country"`
	Address        string `json:"address"`
	Promoted       bool   `json:"promoted"`
	PhoneNumber    string `json:"phone_number"`
	WhatsappNumber string `json:"whatsapp_number"`
	ViberNumber    string `json:"viber_number"`
	Status         string `json:"status"`
	City           string `json:"city"`
	Gps            gps    `json:"gps"`
	DateCreated    string `json:"date_created"`
}

type Agencies []Agency

type gps struct {
	Long string `json:"long"`
	Lat  string `json:"lat"`
}

type UpdateAgencyRequest struct {
	Field string      `json:"field"`
	Value interface{} `json:"Value"`
}
