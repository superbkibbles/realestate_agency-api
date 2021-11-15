package agencyService

import (
	"mime/multipart"
	"strings"

	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_agency-api/src/domain/agency"
	"github.com/superbkibbles/realestate_agency-api/src/domain/agency/esUpdate"
	"github.com/superbkibbles/realestate_agency-api/src/domain/query"
	"github.com/superbkibbles/realestate_agency-api/src/repository/db"
	"github.com/superbkibbles/realestate_agency-api/src/utils/date_utils"
	"github.com/superbkibbles/realestate_agency-api/src/utils/file_utils"
)

type AgencyService interface {
	SaveAgency(agency *agency.Agency) rest_errors.RestErr
	GetAllAgencies() (agency.Agencies, rest_errors.RestErr)
	GetByID(id string) (*agency.Agency, rest_errors.RestErr)
	UploadIcon(id string, fileHeader *multipart.FileHeader) (*agency.Agency, rest_errors.RestErr)
	Update(id string, updateRequest esUpdate.EsUpdate) (*agency.Agency, rest_errors.RestErr)
	Search(updateRequest query.EsQuery) (agency.Agencies, rest_errors.RestErr)
	DeleteIcon(agencyID string) rest_errors.RestErr
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

func (srv *agencyservice) UploadIcon(id string, fileHeader *multipart.FileHeader) (*agency.Agency, rest_errors.RestErr) {
	agency, err := srv.GetByID(id)
	if err != nil {
		return nil, err
	}
	file, fErr := fileHeader.Open()
	if fErr != nil {
		return nil, rest_errors.NewInternalServerErr("Error while trying to open the file", nil)
	}
	filePath, err := file_utils.SaveFile(fileHeader, file)
	if err != nil {
		return nil, err
	}
	agency.Icon = "http://localhost:3031/assets/" + filePath

	srv.db.UploadIcon(agency, id)
	return agency, nil
}

func (srv *agencyservice) Update(id string, updateRequest esUpdate.EsUpdate) (*agency.Agency, rest_errors.RestErr) {
	return srv.db.Update(id, updateRequest)
}

func (srv *agencyservice) Search(updateRequest query.EsQuery) (agency.Agencies, rest_errors.RestErr) {
	return srv.db.Search(updateRequest)
}

func (srv *agencyservice) DeleteIcon(agencyID string) rest_errors.RestErr {
	agency, err := srv.GetByID(agencyID)
	if err != nil {
		return err
	}

	splittedPath := strings.Split(agency.Icon, "/")
	fileName := splittedPath[len(splittedPath)-1]

	file_utils.DeleteFile(fileName)

	agency.Icon = ""
	srv.db.UploadIcon(agency, agencyID)
	return nil
}
