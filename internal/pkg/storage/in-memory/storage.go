package memory

import (
	"Mongo/internal/pkg/constErr"
	"Mongo/internal/pkg/data"
	"sync"
)

type Storage struct {
	accounts []*data.Account
	ads      []*data.Ad
	rwm      *sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		accounts: make([]*data.Account, 0), // 10 acc TODO
		ads:      make([]*data.Ad, 0),
		rwm:      new(sync.RWMutex),
	}
}

//----------------------------------------------------------
// AD
func (s *Storage) Get(id uint) (*data.Ad, error) {
	size, err := s.Size()
	if err != nil {
		return nil, err
	} else if int(id) > size {
		return nil, constErr.InvalidID
	}

	//TODO
	s.rwm.RLock()
	defer s.rwm.RUnlock()
	for _, ad := range s.ads {
		if ad.ID == id {
			return ad, nil
		}
	}
	return nil, constErr.NotFoundAd
}

func (s *Storage) GetAll() ([]*data.Ad, error) {
	if _, err := s.Size(); err != nil {
		return nil, err
	}
	s.rwm.RLock()
	defer s.rwm.RUnlock()
	return s.ads, nil
}

//TODO
func (s *Storage) Add(ad *data.Ad) error {
	if ad == nil {
		return constErr.AdIsNil
	}
	s.rwm.RLock()
	defer s.rwm.RUnlock()
	// create last ID
	lastID := uint(len(s.ads) + 1)
	// set ID
	ad.SetID(lastID)
	s.ads = append(s.ads, ad)
	return nil
}

func (s *Storage) Delete(id uint) error {
	size, err := s.Size()
	if err != nil {
		return err
	} else if int(id) > size {
		return constErr.InvalidID
	}

	s.rwm.RLock()
	defer s.rwm.RUnlock()
	isFind := false
	for i, ad := range s.ads {
		if ad.ID == id {
			if i == len(s.ads)-1 {
				s.ads[i] = nil
				s.ads = s.ads[:i]
			} else {
				s.ads = append(s.ads[:i], s.ads[i+1:]...)
			}
			isFind = true
			break
		}
	}

	if !isFind {
		return constErr.NotFoundAd
	}

	var index uint = 1
	for _, car := range s.ads {
		if car.ID != index {
			car.ID = index
			index++
		} else {
			index++
		}
	}

	return nil
}

func (s *Storage) Update(tempAd *data.Ad, id uint) error {
	if tempAd == nil {
		return constErr.AdIsNil
	}
	size, err := s.Size()
	if err != nil {
		return err
	} else if int(id) > size {
		return constErr.InvalidID
	}
	s.rwm.RLock()
	defer s.rwm.RUnlock()

	for _, ad := range s.ads {
		if ad.ID == id {
			ad.Brand = tempAd.GetBrand()
			ad.Model = tempAd.GetModel()
			ad.Color = tempAd.GetColor()
			ad.Price = tempAd.GetPrice()
			return nil
		}
	}

	return constErr.NotFoundAd
}

func (s *Storage) Size() (int, error) {
	s.rwm.RLock()
	defer s.rwm.RUnlock()
	size := len(s.ads)
	if size == 0 {
		return 0, constErr.AdBaseIsEmpty
	}
	return size, nil
}

//----------------------------------------------------------
// ACCOUNT
func (s *Storage) AddAccount(acc *data.Account) error {
	if acc == nil {
		return constErr.AccIsNil
	}
	s.rwm.RLock()
	defer s.rwm.RUnlock()
	s.accounts = append(s.accounts, acc)
	return nil
}

func (s *Storage) GetAccounts() ([]*data.Account, error) {
	s.rwm.RLock()
	defer s.rwm.RUnlock()
	if len(s.accounts) == 0 {
		return nil, constErr.AccountBaseIsEmpty
	}
	return s.accounts, nil
}

func (s *Storage) UpdateTokenCurrentAcc(acc *data.Account, token string) error {
	baseAcc, err := s.GetAccounts()
	if err != nil {
		return err
	}

	for _, account := range baseAcc {
		if account.GetUserName() == acc.GetUserName() && account.GetPassword() == acc.GetPassword() {
			account.SetToken(token)
			return nil
		}
	}
	return constErr.NotFoundAcc
}
