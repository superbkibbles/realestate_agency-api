package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_agency-api/src/domain/agency"
	"github.com/superbkibbles/realestate_agency-api/src/domain/agency/esUpdate"
	"github.com/superbkibbles/realestate_agency-api/src/domain/query"
	"github.com/superbkibbles/realestate_agency-api/src/services/agencyService"
)

type AgencyHandler interface {
	Get(*gin.Context)
	Create(*gin.Context)
	GetByID(*gin.Context)
	UploadIcon(*gin.Context)
	Update(*gin.Context)
	Search(*gin.Context)
	DeleteIcon(*gin.Context)
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
	agencies, err := ah.service.GetAllAgencies()
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

func (ah *agencyHandler) GetByID(c *gin.Context) {
	agencyID := strings.TrimSpace(c.Param("agency_id"))

	agency, err := ah.service.GetByID(agencyID)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusFound, agency)
}

func (ah *agencyHandler) UploadIcon(c *gin.Context) {
	agencyID := strings.TrimSpace(c.Param("agency_id"))

	file, err := c.FormFile("icon")
	if err != nil {
		restErr := rest_errors.NewBadRequestErr("Bad Request")
		c.JSON(restErr.Status(), restErr)
		return
	}

	agency, uploadErr := ah.service.UploadIcon(agencyID, file)
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

	if err := c.ShouldBindJSON(&q); err != nil {
		restErr := rest_errors.NewBadRequestErr("Invalid Body JSON")
		c.JSON(restErr.Status(), restErr)
		return
	}

	properties, err := ah.service.Search(q)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusFound, properties)
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
