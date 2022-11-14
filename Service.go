package odoo

import (
	"fmt"
	"github.com/kolo/xmlrpc"
	errortools "github.com/leapforce-libraries/go_errortools"
)

const (
	apiName string = "Odoo"
	baseUrl string = "https://%s.odoo.com/xmlrpc/2"
)

type Service struct {
	domain       string
	database     string
	username     string
	userId       int
	password     string
	clientCommon *xmlrpc.Client
	clientObject *xmlrpc.Client
	apiCallCount int64
}

type ServiceConfig struct {
	Domain   string
	Database string
	Username string
	Password string
}

func NewService(serviceConfig *ServiceConfig) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	clientCommon, err := xmlrpc.NewClient(fmt.Sprintf("%s/common", fmt.Sprintf(baseUrl, serviceConfig.Domain)), nil)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	clientObject, err := xmlrpc.NewClient(fmt.Sprintf("%s/object", fmt.Sprintf(baseUrl, serviceConfig.Domain)), nil)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	// authenticate
	var userId int
	err = clientCommon.Call("authenticate", []interface{}{serviceConfig.Database, serviceConfig.Username, serviceConfig.Password, ""}, &userId)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	return &Service{
		domain:       serviceConfig.Domain,
		username:     serviceConfig.Username,
		userId:       userId,
		password:     serviceConfig.Password,
		clientCommon: clientCommon,
		clientObject: clientObject,
	}, nil
}

func (service *Service) ApiName() string {
	return apiName
}

func (service *Service) ApiKey() string {
	return service.username
}

func (service *Service) ApiCallCount() int64 {
	return service.apiCallCount
}

func (service *Service) ApiReset() {
	service.apiCallCount = 0
}

type Criterion struct {
	Field    string
	Operator string
	Value    interface{}
}

func (service *Service) execute(method string, model string, criteria *[]Criterion, attributes []string, responseModel interface{}) *errortools.Error {
	service.apiCallCount++

	var options = make(map[string]interface{})
	options["attributes"] = attributes

	var criteria_ = []interface{}{[]interface{}{}}
	if criteria != nil {
		for _, criterion := range *criteria {
			criteria_ = append(criteria_, []interface{}{criterion.Field, criterion.Operator, criterion.Value})
		}
	}

	err := service.clientObject.Call("execute_kw", []interface{}{service.database, service.userId, service.password, model, method, []interface{}{criteria_}, options}, responseModel)
	if err != nil {
		return errortools.ErrorMessage(err)
	}

	return nil
}

func (service *Service) getFields(model string, fields []string, responseModel interface{}) *errortools.Error {
	return service.execute("fields_get", model, nil, fields, responseModel)
}

func (service *Service) searchRead(model string, criteria *[]Criterion, attributes []string, responseModel interface{}) *errortools.Error {
	return service.execute("search_read", model, criteria, attributes, responseModel)
}
