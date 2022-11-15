package odoo

import (
	errortools "github.com/leapforce-libraries/go_errortools"
)

type IrAttachment struct {
	Name     string `xmlrpc:"name"`
	Type     string `xmlrpc:"type"`
	Datas    string `xmlrpc:"datas"`
	ResModel string `xmlrpc:"res_model"`
	ResId    int64  `xmlrpc:"res_id"`
	MimeType string `xmlrpc:"mimetype"`
}

func (service *Service) SearchReadIrAttachments(criteria *[]Criterion, attributes []string) (*[]IrAttachment, *errortools.Error) {
	var irAttachments []IrAttachment

	e := service.searchRead("ir.attachment", criteria, attributes, &irAttachments)
	if e != nil {
		return nil, e
	}

	return &irAttachments, nil
}

func (service *Service) CreateIrAttachment(irAttachment IrAttachment) (int64, *errortools.Error) {
	return service.create("ir.attachment", irAttachment)
}
