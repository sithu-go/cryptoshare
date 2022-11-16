package repository

import "cryptoshare/ds"

type Repository struct {
	DS    *ds.DataSource
	Bank  *bankRepository
	Admin *adminRepository
}

func NewRepository(ds *ds.DataSource) *Repository {
	bankRepo := newBankRepository(ds)
	adminRepo := newAdminRepository(ds)
	return &Repository{
		DS:    ds,
		Bank:  bankRepo,
		Admin: adminRepo,
	}
}
