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
		database:     serviceConfig.Database,
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

func (service *Service) execute(method string, model string, args interface{}, options map[string]interface{}, responseModel interface{}) *errortools.Error {
	service.apiCallCount++

	err := service.clientObject.Call("execute_kw", []interface{}{service.database, service.userId, service.password, model, method, []interface{}{args}, options}, responseModel)
	if err != nil {
		return errortools.ErrorMessage(err)
	}

	return nil
}

func (service *Service) getFields(model string, fields []string, responseModel interface{}) *errortools.Error {
	var options = make(map[string]interface{})
	options["fields"] = fields

	return service.execute("fields_get", model, nil, options, responseModel)
}

func (service *Service) searchRead(model string, criteria *[]Criterion, attributes []string, responseModel interface{}) *errortools.Error {
	var criteria_ []interface{}
	if criteria != nil {
		for _, criterion := range *criteria {
			criteria_ = append(criteria_, []interface{}{criterion.Field, criterion.Operator, criterion.Value})
		}
	}

	var options = make(map[string]interface{})
	options["fields"] = attributes

	return service.execute("search_read", model, criteria_, options, responseModel)
}

func (service *Service) create(model string, data interface{}) (int64, *errortools.Error) {
	var id int64

	e := service.execute("create", model, data, nil, &id)
	if e != nil {
		return 0, e
	}

	return id, nil
}
