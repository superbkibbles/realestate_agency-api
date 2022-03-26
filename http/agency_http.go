package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_agency-api/domain/agency"
	"github.com/superbkibbles/realestate_agency-api/domain/agency/esUpdate"
	"github.com/superbkibbles/realestate_agency-api/domain/query"
	"github.com/superbkibbles/realestate_agency-api/services/agencyService"
)

type AgencyHandler interface {
	Get(*gin.Context)
	Create(*gin.Context)
	GetByID(*gin.Context)
	UploadIcon(*gin.Context)
	Update(*gin.Context)
	Search(*gin.Context)
	DeleteIcon(*gin.Context)
	Translate(*gin.Context)
}

type agencyHandler struct {
	service agencyService.AgencyService
}

func NewAgencyHandler(srv agencyService.AgencyService) AgencyHandler {
	return &agencyHandler{
		service: srv,
	}
}

func (ah *agencyHandler) Get(c *gin.Context) {
	local := c.GetHeader("local")
	agencies, err := ah.service.GetAllAgencies(local)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, agencies)
}

func (ah *agencyHandler) Create(c *gin.Context) {
	var agency agency.Agency
	if err := c.ShouldBindJSON(&agency); err != nil {
		restErr := rest_errors.NewBadRequestErr("Bad Request Body JSON")
		c.JSON(restErr.Status(), restErr)
		return
	}

	if err := ah.service.SaveAgency(&agency); err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusCreated, agency)
}

func (ah *agencyHandler) Translate(c *gin.Context) {
	agencyID := strings.TrimSpace(c.Param("agency_id"))
	local := c.GetHeader("local")
	var agencyRequest agency.TranslateRequest
	if err := c.ShouldBindJSON(&agencyRequest); err != nil {
		restErr := rest_errors.NewBadRequestErr("Bad Request Body JSON")
		c.JSON(restErr.Status(), restErr)
		return
	}
	agency, err := ah.service.Translate(agencyID, agencyRequest, local)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, agency)
}

func (ah *agencyHandler) GetByID(c *gin.Context) {
	agencyID := strings.TrimSpace(c.Param("agency_id"))
	local := c.GetHeader("local")

	agency, err := ah.service.GetByID(agencyID, local)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, agency)
}

func (ah *agencyHandler) UploadIcon(c *gin.Context) {
	agencyID := strings.TrimSpace(c.Param("agency_id"))

	file, err := c.FormFile("icon")
	if err != nil {
		restErr := rest_errors.NewBadRequestErr("Bad Request")
		c.JSON(restErr.Status(), restErr)
		return
	}

	headerPhoto, err := c.FormFile("header_photo")
	if err != nil {
		restErr := rest_errors.NewBadRequestErr("Bad Request")
		c.JSON(restErr.Status(), restErr)
		return
	}

	agency, uploadErr := ah.service.UploadIcon(agencyID, file, headerPhoto)
	if err != nil {
		c.JSON(uploadErr.Status(), uploadErr)
		return
	}

	c.JSON(http.StatusOK, agency)
}

func (ah *agencyHandler) Update(c *gin.Context) {
	id := strings.TrimSpace(c.Param("agency_id"))
	var updateRequest esUpdate.EsUpdate

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		restErr := rest_errors.NewBadRequestErr("Invalid Body JSON")
		c.JSON(restErr.Status(), restErr)
		return
	}

	agency, err := ah.service.Update(id, updateRequest)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, agency)
}

func (ah *agencyHandler) Search(c *gin.Context) {
	var q query.EsQuery
	local := c.GetHeader("local")

	if err := c.ShouldBindJSON(&q); err != nil {
		restErr := rest_errors.NewBadRequestErr("Invalid Body JSON")
		c.JSON(restErr.Status(), restErr)
		return
	}

	properties, err := ah.service.Search(q, local)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, properties)
}

func (ah *agencyHandler) DeleteIcon(c *gin.Context) {
	agencyID := strings.TrimSpace(c.Param("agency_id"))

	err := ah.service.DeleteIcon(agencyID)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.String(200, "Icon Deleted")
}
