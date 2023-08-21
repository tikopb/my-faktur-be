package partner

import "bemyfaktur/internal/model"

type Usecase interface {
	IndexPartner() ([]model.Partner, error)
	GetPartner(id int) (model.Partner, error)
	CreatePartner(request model.Partner) (model.PartnerRespon, error)
	UpdatedPartner(id int, request model.Partner) (model.Partner, error)
	Deletepartner(id int) (string, error)
}
