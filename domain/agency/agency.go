package agency

const (
	STATUS_ACTIVE   = "active"
	STATUS_DEACTIVE = "deactive"
)

type Agency struct {
	ID   string  `json:"id"`
	Ar   Arabic  `json:"ar"`
	Kur  Kurdish `json:"kur"`
	Name string  `json:"name"`

	Icon          string `json:"icon"`
	PublicID      string `json:"public_id"`
	BackgroundPic string `json:"background_pic"`

	Email           string `json:"email"`
	Country         string `json:"country"`
	Description     string `json:"description"`
	Address         string `json:"address"`
	Promoted        bool   `json:"promoted"`
	IsSponsored     bool   `json:"is_sponsored"`
	PhoneNumber     string `json:"phone_number"`
	WhatsappNumber  string `json:"whatsapp_number"`
	ViberNumber     string `json:"viber_number"`
	Status          string `json:"status"`
	City            string `json:"city"`
	Gps             gps    `json:"gps"`
	DateCreated     string `json:"date_created"`
	BackgroundColor string `json:"background_color"`
}

// type icon struct {
// 	Url      string `json:"url"`
// 	PublicID string `json:"public_id"`
// }

type TranslateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
}

type Arabic struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
}

type Kurdish struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
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
