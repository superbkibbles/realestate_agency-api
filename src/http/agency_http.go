package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_agency-api/src/domain/agency"
	"github.com/superbkibbles/realestate_agency-api/src/services/agencyService"
)

type AgencyHandler interface {
	Get(*gin.Context)
	Create(*gin.Context)
	GetByID(*gin.Context)
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
