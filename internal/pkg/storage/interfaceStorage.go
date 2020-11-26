package storage

import "Mongo/internal/pkg/data"

type InterfaceStorage interface {
	Add(ad *data.Ad) error
	Get(id uint) (*data.Ad, error)
	GetAll() ([]*data.Ad, error)
	Update(temp *data.Ad, id uint) error
	Delete(id uint) error
	Size() (int, error)
	AddAccount(acc *data.Account) error
	GetAccounts() ([]*data.Account, error)
	UpdateTokenCurrentAcc(*data.Account, string) error
}
