package partner

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/repository/partner"
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
func (m *partnerUsecase) IndexPartner() ([]model.Partner, error) {
	return m.partnerRepo.Index()
}

// GetPartner implements Usecase.
func (m *partnerUsecase) GetPartner(id int) (model.Partner, error) {
	return m.partnerRepo.Show(id)
}

// CreatePartner implements Partner_Usecase.
func (m *partnerUsecase) CreatePartner(request model.Partner) (model.PartnerRespon, error) {
	data, err := m.partnerRepo.Create(request)
	return data, err
}

// UpdatedPartner implements Partner_Usecase.
func (m *partnerUsecase) UpdatedPartner(id int, request model.Partner) (model.Partner, error) {
	data, err := m.partnerRepo.Update(id, request)
	return data, err
}

// Deletepartner implements Partner_Usecase.
func (m *partnerUsecase) Deletepartner(id int) (string, error) {
	data, err := m.partnerRepo.Delete(id)
	return data, err
}
