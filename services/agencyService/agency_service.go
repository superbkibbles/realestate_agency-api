package agencyService

import (
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_agency-api/domain/agency"
	"github.com/superbkibbles/realestate_agency-api/domain/agency/esUpdate"
	"github.com/superbkibbles/realestate_agency-api/domain/query"
	cloudstorage "github.com/superbkibbles/realestate_agency-api/repository/cloudStorage"
	"github.com/superbkibbles/realestate_agency-api/repository/db"
	"github.com/superbkibbles/realestate_agency-api/utils/crypto_utils"
	"github.com/superbkibbles/realestate_agency-api/utils/date_utils"
)

type AgencyService interface {
	SaveAgency(agency *agency.Agency) rest_errors.RestErr
	GetAllAgencies(local string) (agency.Agencies, rest_errors.RestErr)
	GetByID(id string, local string) (*agency.Agency, rest_errors.RestErr)
	UploadIcon(id string, fileHeader *multipart.FileHeader) (*agency.Agency, rest_errors.RestErr)
	Update(id string, updateRequest esUpdate.EsUpdate) (*agency.Agency, rest_errors.RestErr)
	Search(updateRequest query.EsQuery, local string) (agency.Agencies, rest_errors.RestErr)
	DeleteIcon(agencyID string) rest_errors.RestErr
	Translate(agencyID string, agencyRequest agency.TranslateRequest, local string) (*agency.Agency, rest_errors.RestErr)
}

type agencyservice struct {
	db        db.DbRepository
	cloudRepo cloudstorage.CloudStorage
}

func NewAgencyService(db db.DbRepository, cloudRepo cloudstorage.CloudStorage) AgencyService {
	return &agencyservice{
		db:        db,
		cloudRepo: cloudRepo,
	}
}

func (srv *agencyservice) SaveAgency(a *agency.Agency) rest_errors.RestErr {
	a.Status = agency.STATUS_ACTIVE
	a.DateCreated = date_utils.GetNowDBFromat()
	return srv.db.SaveAgency(a)
}

func (srv *agencyservice) GetAllAgencies(local string) (agency.Agencies, rest_errors.RestErr) {
	agencies, err := srv.db.GetAllAgencies()
	if err != nil {
		return nil, err
	}

	if local == "ar" {
		for i, agency := range agencies {
			if agencies[i].Ar.Address != "" {
				agencies[i].Address = agency.Ar.Address
			}

			if agencies[i].Ar.Description != "" {
				agencies[i].Description = agency.Ar.Description
			}

			if agencies[i].Ar.Name != "" {
				agencies[i].Name = agency.Ar.Name
			}
		}
	}

	if local == "kur" {
		for i, agency := range agencies {
			if agencies[i].Kur.Address != "" {
				agencies[i].Address = agency.Kur.Address
			}

			if agencies[i].Kur.Description != "" {
				agencies[i].Description = agency.Kur.Description
			}

			if agencies[i].Ar.Name != "" {
				agencies[i].Name = agency.Kur.Name
			}
		}
	}

	return agencies, err
}

func (srv *agencyservice) GetByID(id string, local string) (*agency.Agency, rest_errors.RestErr) {
	agency, err := srv.db.GetByID(id)
	if err != nil {
		return nil, err
	}
	if local == "ar" {
		if agency.Ar.Address != "" {
			agency.Address = agency.Ar.Address
		}

		if agency.Ar.Description != "" {
			agency.Description = agency.Ar.Description
		}

		if agency.Ar.Name != "" {
			agency.Name = agency.Ar.Name
		}
	}
	if local == "kur" {
		if agency.Kur.Address != "" {
			agency.Address = agency.Kur.Address
		}

		if agency.Kur.Description != "" {
			agency.Description = agency.Kur.Description
		}

		if agency.Kur.Name != "" {
			agency.Name = agency.Kur.Name
		}
	}
	return agency, nil
}

func (srv *agencyservice) Translate(agencyID string, agencyRequest agency.TranslateRequest, local string) (*agency.Agency, rest_errors.RestErr) {
	agency, err := srv.db.GetByID(agencyID)
	if err != nil {
		return nil, err
	}

	if local == "ar" || local == "kur" {
		agency, err = srv.db.Translate(agencyID, agencyRequest, local)
		if err != nil {
			return nil, err
		}
	}

	return agency, nil
}

func (srv *agencyservice) UploadIcon(id string, fileHeader *multipart.FileHeader) (*agency.Agency, rest_errors.RestErr) {
	agency, err := srv.GetByID(id, "")
	if err != nil {
		return nil, err
	}
	file, fErr := fileHeader.Open()
	if fErr != nil {
		return nil, rest_errors.NewInternalServerErr("Error while trying to open the file", nil)
	}

	res, cloudErr := srv.cloudRepo.Save(file, id+crypto_utils.GetMd5(uuid.New().String()), id)
	if cloudErr != nil {
		return nil, cloudErr
	}
	agency.Icon = res.Url
	agency.PublicID = res.PublicID

	srv.db.UploadIcon(agency, id)
	return agency, nil
}

func (srv *agencyservice) Update(id string, updateRequest esUpdate.EsUpdate) (*agency.Agency, rest_errors.RestErr) {
	return srv.db.Update(id, updateRequest)
}

func (srv *agencyservice) Search(updateRequest query.EsQuery, local string) (agency.Agencies, rest_errors.RestErr) {
	agencies, err := srv.db.Search(updateRequest)
	if err != nil {
		return nil, err
	}
	if local == "ar" {
		for i, agency := range agencies {
			if agencies[i].Ar.Address != "" {
				agencies[i].Address = agency.Ar.Address
			}

			if agencies[i].Ar.Description != "" {
				agencies[i].Description = agency.Ar.Description
			}

			if agencies[i].Ar.Name != "" {
				agencies[i].Name = agency.Ar.Name
			}
		}
	}

	if local == "kur" {
		for i, agency := range agencies {
			if agencies[i].Kur.Address != "" {
				agencies[i].Address = agency.Kur.Address
			}

			if agencies[i].Kur.Description != "" {
				agencies[i].Description = agency.Kur.Description
			}

			if agencies[i].Ar.Name != "" {
				agencies[i].Name = agency.Kur.Name
			}
		}
	}

	return agencies, nil
}

func (srv *agencyservice) DeleteIcon(agencyID string) rest_errors.RestErr {
	agency, err := srv.GetByID(agencyID, "")
	if err != nil {
		return err
	}
	if err := srv.cloudRepo.Delete(agency.PublicID); err != nil {
		return err
	}

	agency.Icon = ""
	agency.PublicID = ""
	srv.db.UploadIcon(agency, agencyID)
	return nil
}
