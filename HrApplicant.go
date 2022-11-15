package odoo

import (
	errortools "github.com/leapforce-libraries/go_errortools"
)

type HrApplicant struct {
	Name         string `xmlrpc:"name"`
	JobId        int64  `xmlrpc:"job_id"`
	EmailFrom    string `xmlrpc:"email_from"`
	PartnerName  string `xmlrpc:"partner_name"`
	PartnerPhone string `xmlrpc:"partner_phone"`
}

func (service *Service) SearchReadHrApplicants(criteria *[]Criterion, attributes []string) (*[]HrApplicant, *errortools.Error) {
	var hrApplicants []HrApplicant

	e := service.searchRead("hr.applicant", criteria, attributes, &hrApplicants)
	if e != nil {
		return nil, e
	}

	return &hrApplicants, nil
}

func (service *Service) CreateHrApplicant(hrApplicant HrApplicant) (int64, *errortools.Error) {
	return service.create("hr.applicant", hrApplicant)
}
