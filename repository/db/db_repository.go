package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_agency-api/clients/elasticsearch"
	"github.com/superbkibbles/realestate_agency-api/domain/agency"
	"github.com/superbkibbles/realestate_agency-api/domain/agency/esUpdate"
	"github.com/superbkibbles/realestate_agency-api/domain/query"
)

const (
	indexAgency = "agency"
	docType     = "_doc"
)

type DbRepository interface {
	SaveAgency(agency *agency.Agency) rest_errors.RestErr
	GetAllAgencies() (agency.Agencies, rest_errors.RestErr)
	GetByID(id string) (*agency.Agency, rest_errors.RestErr)
	UploadIcon(agency *agency.Agency, id string) rest_errors.RestErr
	Update(id string, updateRequest esUpdate.EsUpdate) (*agency.Agency, rest_errors.RestErr)
	Search(query query.EsQuery) (agency.Agencies, rest_errors.RestErr)
}

type dbRepository struct{}

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
		if strings.Contains(err.Error(), "404") {
			return nil, rest_errors.NewNotFoundErr(fmt.Sprintf("no Property was found with id %s", id))
		}
		return nil, rest_errors.NewInternalServerErr(fmt.Sprintf("error when trying to id %s", id), errors.New("database error"))
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

func (db *dbRepository) UploadIcon(agency *agency.Agency, id string) rest_errors.RestErr {
	var es esUpdate.EsUpdate
	update := esUpdate.UpdatePropertyRequest{
		Field: "icon",
		Value: agency.Icon,
	}
	es.Fields = append(es.Fields, update)
	_, err := db.Update(id, es)
	if err != nil {
		return err
	}
	return nil
}

func (db *dbRepository) Update(id string, updateRequest esUpdate.EsUpdate) (*agency.Agency, rest_errors.RestErr) {
	result, err := elasticsearch.Client.Update(indexAgency, docType, id, updateRequest)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, rest_errors.NewNotFoundErr(fmt.Sprintf("no Property was found with id %s", id))
		}
		return nil, rest_errors.NewInternalServerErr("error when trying to Update Property", errors.New("databse error"))
	}

	var ag agency.Agency

	bytes, err := result.GetResult.Source.MarshalJSON()
	if err != nil {
		return nil, rest_errors.NewInternalServerErr(fmt.Sprintf("error when trying to parse database response"), errors.New("database error"))
	}
	if err := json.Unmarshal(bytes, &ag); err != nil {
		return nil, rest_errors.NewInternalServerErr(fmt.Sprintf("error when trying to parse database response"), errors.New("database error"))
	}

	ag.ID = result.Id
	return &ag, nil
}

func (db *dbRepository) Search(query query.EsQuery) (agency.Agencies, rest_errors.RestErr) {
	result, err := elasticsearch.Client.Search(indexAgency, query.Build())
	if err != nil {
		return nil, rest_errors.NewInternalServerErr("error when trying to search documents", errors.New("database error"))
	}

	properties := make(agency.Agencies, result.TotalHits())
	for i, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()
		var property agency.Agency
		if err := json.Unmarshal(bytes, &property); err != nil {
			return nil, rest_errors.NewInternalServerErr("error when trying to parse response", errors.New("database error"))
		}
		property.ID = hit.Id
		properties[i] = property
	}

	if len(properties) == 0 {
		return nil, rest_errors.NewNotFoundErr("no items found matching given critirial")
	}

	return properties, nil
}
