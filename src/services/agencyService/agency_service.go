package agencyService

import (
	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_agency-api/src/domain/agency"
	"github.com/superbkibbles/realestate_agency-api/src/repository/db"
	"github.com/superbkibbles/realestate_agency-api/src/utils/date_utils"
)

type AgencyService interface {
	SaveAgency(agency *agency.Agency) rest_errors.RestErr
	GetAllAgencies() (agency.Agencies, rest_errors.RestErr)
	GetByID(id string) (*agency.Agency, rest_errors.RestErr)
}

type agencyservice struct {
	db db.DbRepository
}

func NewAgencyService(db db.DbRepository) AgencyService {
	return &agencyservice{
		db: db,
	}
}

func (srv *agencyservice) SaveAgency(a *agency.Agency) rest_errors.RestErr {
	a.Status = agency.STATUS_ACTIVE
	a.DateCreated = date_utils.GetNowDBFromat()
	return srv.db.SaveAgency(a)
}

func (srv *agencyservice) GetAllAgencies() (agency.Agencies, rest_errors.RestErr) {
	return srv.db.GetAllAgencies()
}

func (srv *agencyservice) GetByID(id string) (*agency.Agency, rest_errors.RestErr) {
	return srv.db.GetByID(id)
}
