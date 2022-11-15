package odoo

import (
	errortools "github.com/leapforce-libraries/go_errortools"
)

type HrJob struct {
	Id               int64  `xmlrpc:"id"`
	Name             string `xmlrpc:"name"`
	WebsitePublished bool   `xmlrpc:"website_published"`
}

func (service *Service) SearchReadHrJobs(criteria *[]Criterion, attributes []string) (*[]HrJob, *errortools.Error) {
	var hrJobs []HrJob

	e := service.searchRead("hr.job", criteria, attributes, &hrJobs)
	if e != nil {
		return nil, e
	}

	return &hrJobs, nil
}
