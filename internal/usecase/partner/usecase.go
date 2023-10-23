package partner

import "bemyfaktur/internal/model"

type Usecase interface {
	IndexPartner(limit int, offset int, q string) ([]model.Partner, error)
	GetPartner(id int) (model.Partner, error)
	CreatePartner(request model.Partner, userID string) (model.PartnerRespon, error)
	UpdatedPartner(id int, request model.Partner) (model.PartnerRespon, error)
	Deletepartner(id int) (string, error)
}
