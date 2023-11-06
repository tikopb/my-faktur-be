package partner

import (
	"bemyfaktur/internal/model"

	"github.com/google/uuid"
)

type Usecase interface {
	IndexPartner(limit int, offset int, q string, order []string) ([]model.PartnerRespon, error)
	GetPartner(id uuid.UUID) (model.Partner, error)
	CreatePartner(request model.Partner, userID string) (model.PartnerRespon, error)
	UpdatedPartner(id uuid.UUID, request model.Partner) (model.PartnerRespon, error)
	Deletepartner(id uuid.UUID) (string, error)
	PartialGet(q string) ([]model.PartnerPartialRespon, error)
}
