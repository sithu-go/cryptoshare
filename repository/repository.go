package repository

import "cryptoshare/ds"

type Repository struct {
	DS   *ds.DataSource
	Bank *bankRepository
}

func NewRepository(ds *ds.DataSource) *Repository {
	bankRepo := newBankRepository(ds)
	return &Repository{
		DS:   ds,
		Bank: bankRepo,
	}
}
