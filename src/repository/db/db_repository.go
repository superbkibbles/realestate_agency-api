package db

import (
	"encoding/json"
	"errors"

	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_agency-api/src/clients/elasticsearch"
	"github.com/superbkibbles/realestate_agency-api/src/domain/agency"
)

const (
	indexAgency = "agency"
	docType     = "_doc"
)

type DbRepository interface {
	SaveAgency(agency *agency.Agency) rest_errors.RestErr
	GetAllAgencies() (agency.Agencies, rest_errors.RestErr)
	GetByID(id string) (*agency.Agency, rest_errors.RestErr)
}

type dbRepository struct {
}

func NewDbRepository() DbRepository {
	return &dbRepository{}
}

func (db *dbRepository) SaveAgency(agency *agency.Agency) rest_errors.RestErr {
	result, err := elasticsearch.Client.Save(indexAgency, docType, agency)
	if err != nil {
		return rest_errors.NewInternalServerErr("error when trying to save Property", errors.New("databse error"))
	}

	agency.ID = result.Id

	return nil
}

func (db *dbRepository) GetAllAgencies() (agency.Agencies, rest_errors.RestErr) {
	result, err := elasticsearch.Client.GetAllDoc(indexAgency)
	if err != nil {
		return nil, rest_errors.NewInternalServerErr("error when trying to Get All Agencies Property", errors.New("databse error"))
	}

	agencies := make(agency.Agencies, result.TotalHits())
	for i, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()
		var agency agency.Agency
		if err := json.Unmarshal(bytes, &agency); err != nil {
			return nil, rest_errors.NewInternalServerErr("error when trying to parse response", errors.New("database error"))
		}
		agency.ID = hit.Id
		agencies[i] = agency
	}
	return agencies, nil
}

func (db *dbRepository) GetByID(id string) (*agency.Agency, rest_errors.RestErr) {
	result, err := elasticsearch.Client.GetByID(indexAgency, docType, id)
	if err != nil {
		return nil, rest_errors.NewInternalServerErr("error when trying to Get Agency By id", errors.New("database error"))
	}

	var agency agency.Agency

	bytes, err := result.Source.MarshalJSON()
	if err != nil {
		return nil, rest_errors.NewInternalServerErr("error when trying to parse database response", errors.New("database error"))
	}

	if err := json.Unmarshal(bytes, &agency); err != nil {
		return nil, rest_errors.NewInternalServerErr("error when trying to parse response", errors.New("database error"))
	}

	agency.ID = result.Id
	return &agency, nil
}
