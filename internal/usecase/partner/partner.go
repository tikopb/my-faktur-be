package partner

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/repository/partner"

	"github.com/google/uuid"
)

type partnerUsecase struct {
	partnerRepo partner.Repository
}

func GetUsecase(partnerRepo partner.Repository) Usecase {
	return &partnerUsecase{
		partnerRepo: partnerRepo,
	}
}

// GetPartner implements Usecase.
func (m *partnerUsecase) IndexPartner(limit int, offset int, q string, order []string) ([]model.PartnerRespon, error) {
	return m.partnerRepo.Index(limit, offset, q, order)
}

// GetPartner implements Usecase.
func (m *partnerUsecase) GetPartner(id uuid.UUID) (model.PartnerRespon, error) {
	return m.partnerRepo.Show(id)
}

// CreatePartner implements Partner_Usecase.
func (m *partnerUsecase) CreatePartner(request model.Partner, userID string) (model.PartnerRespon, error) {
	request.CreatedBy = userID
	data, err := m.partnerRepo.Create(request)
	return data, err
}

// UpdatedPartner implements Partner_Usecase.
func (m *partnerUsecase) UpdatedPartner(id uuid.UUID, request model.Partner) (model.PartnerRespon, error) {
	data, err := m.partnerRepo.Update(id, request)
	return data, err
}

// Deletepartner implements Partner_Usecase.
func (m *partnerUsecase) Deletepartner(id uuid.UUID) (string, error) {
	data, err := m.partnerRepo.Delete(id)
	return data, err
}

// PartialGet implements Usecase.
func (m *partnerUsecase) PartialGet(q string) ([]model.PartnerPartialRespon, error) {
	data, err := m.partnerRepo.Partial(q)
	return data, err
}
