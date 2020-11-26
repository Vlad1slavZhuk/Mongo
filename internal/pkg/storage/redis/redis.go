package redis

import (
	"Mongo/internal/pkg/constErr"
	"Mongo/internal/pkg/data"
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/xuyu/goredis"
)

type StorageRedis struct {
	Conn       *goredis.Redis
	StorageAcc string
	StorageAd  string
}

func NewRedisStr() *StorageRedis {
	conn, err := goredis.Dial(&goredis.DialConfig{
		Network:  "tcp",
		Address:  "redis:8002",
		Database: 0,
		Password: "vlad",
	})
	if err != nil {
		log.Fatal(err)
	}
	return &StorageRedis{
		Conn:       conn,
		StorageAcc: "accounts",
		StorageAd:  "ads",
	}
}

func (s *StorageRedis) Add(ad *data.Ad) error {
	size, _ := s.Size()
	lastID := int(size) + 1
	_, err := s.Conn.ExecuteCommand(
		"HSET", fmt.Sprintf("%v:%v", s.StorageAd, lastID),
		"id", lastID,
		"brand", ad.GetBrand(),
		"model", ad.GetModel(),
		"color", ad.GetColor(),
		"price", ad.GetPrice(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageRedis) Get(id uint) (*data.Ad, error) {
	reply, err := s.Conn.HGetAll(fmt.Sprintf("%v:%v", s.StorageAd, id))
	if err != nil {
		return nil, err
	}
	ad := new(data.Ad)

	ID, err := strconv.Atoi(reply["id"]) // converter string to int
	if err != nil {
		return nil, err
	}
	price, err := strconv.Atoi(reply["price"]) // converter string to int
	if err != nil {
		return nil, err
	}

	ad.ID = uint(ID)
	ad.Brand = reply["brand"]
	ad.Model = reply["model"]
	ad.Color = reply["color"]
	ad.Price = price

	return ad, nil
}

func (s *StorageRedis) GetAll() ([]*data.Ad, error) {
	if size, err := s.Size(); err != nil || size == 0 {
		return nil, err
	}
	keys, err := s.Conn.Keys(s.StorageAd + ":*")
	if err != nil {
		return nil, err
	}

	ads := make([]*data.Ad, 0, len(keys))

	for _, otherAd := range keys {
		reply, err := s.Conn.HGetAll(otherAd)
		if err != nil {
			return nil, err
		}

		ad := new(data.Ad)

		ID, err := strconv.Atoi(reply["id"]) // converter string to int
		if err != nil {
			return nil, err
		}
		price, err := strconv.Atoi(reply["price"]) // converter string to int
		if err != nil {
			return nil, err
		}

		ad.ID = uint(ID)
		ad.Brand = reply["brand"]
		ad.Model = reply["model"]
		ad.Color = reply["color"]
		ad.Price = price

		ads = append(ads, ad)
	}
	log.Println(ads)
	sort.Slice(ads, func(i, j int) bool {
		return ads[i].ID < ads[j].ID
	})
	log.Println(ads)
	return ads, nil
}

func (s *StorageRedis) Update(temp *data.Ad, id uint) error {
	if ok, _ := s.Conn.Exists(fmt.Sprintf("%v:%v", s.StorageAd, id)); !ok {
		return constErr.NotFoundAd
	}

	reply, err := s.Conn.ExecuteCommand(
		"HMSET", fmt.Sprintf("%v:%v", s.StorageAd, id),
		"id", fmt.Sprintf("%v", temp.GetID()),
		"brand", temp.GetBrand(),
		"model", temp.GetModel(),
		"color", temp.GetColor(),
		"price", fmt.Sprintf("%v", temp.GetPrice()),
	)
	if err != nil {
		return err
	}
	log.Println(reply)

	return nil
}

func (s *StorageRedis) Delete(id uint) error {
	if ok, _ := s.Conn.Exists(fmt.Sprintf("%v:%v", s.StorageAd, id)); !ok {
		return constErr.NotFoundAd
	}

	_, err := s.Conn.Del(fmt.Sprintf("%v:%v", s.StorageAd, id))
	if err != nil {
		return err
	}

	count := 1
	keys, err := s.Conn.Keys(s.StorageAd + ":*")
	if err != nil {
		return err
	}
	sort.Strings(keys)
	for _, str := range keys {
		_, err := s.Conn.HSet(str, "id", fmt.Sprintf("%v", count))
		if err != nil {
			return err
		}
		err = s.Conn.Rename(str, fmt.Sprintf("%v:%v", s.StorageAd, count))
		if err != nil {
			return err
		}
		count++
	}

	return nil
}

func (s *StorageRedis) AddAccount(acc *data.Account) error {
	sacc, _ := s.Conn.Keys(s.StorageAcc + ":*")
	log.Println(sacc)
	id := len(sacc) + 1
	acc.ID = id
	reply, err := s.Conn.ExecuteCommand(
		"HSET", fmt.Sprintf("%v:%v", s.StorageAcc, acc.GetID()), // Key
		"id", fmt.Sprintf("%v", acc.GetID()), // Field Value
		"username", acc.GetUserName(), // Field Value
		"password", acc.GetPassword(), // Field Value
		"token", acc.GetToken(), // Field Value
	)
	if err != nil {
		return err
	}
	test, _ := s.Conn.HGetAll(fmt.Sprintf("%v:%v", s.StorageAcc, acc.GetID()))
	log.Println("Status Add:", test, reply.OKValue())
	return nil
}

func (s *StorageRedis) GetAccounts() ([]*data.Account, error) {
	keys, err := s.Conn.Keys(s.StorageAcc + ":*")
	if err != nil {
		return nil, err
	}
	log.Println(keys)
	accs := make([]*data.Account, 0, len(keys))

	for _, a := range keys {
		tempAcc, _ := s.Conn.HGetAll(a)
		acc := new(data.Account)
		for f, v := range tempAcc { // f - field, v - value
			if f == "id" {
				id, _ := strconv.Atoi(v)
				acc.ID = id
			}
			if f == "username" {
				acc.Username = v
			}
			if f == "password" {
				acc.Password = v
			}
			if f == "token" {
				acc.Token = v
			}
		}
		accs = append(accs, acc)
		log.Println("in for", accs)
	}
	log.Println("result", accs)
	return accs, nil
}

func (s *StorageRedis) UpdateTokenCurrentAcc(acc *data.Account, token string) error {
	keys, err := s.Conn.Keys(s.StorageAcc + ":*")
	if err != nil || len(keys) == 0 {
		return constErr.AccountBaseIsEmpty
	}

	for _, account := range keys {
		otherAcc, err := s.Conn.HGetAll(account)
		if err != nil {
			return err
		}
		if otherAcc["username"] == acc.GetUserName() && otherAcc["password"] == acc.GetPassword() {
			_, err := s.Conn.HSet(s.StorageAcc+":"+fmt.Sprintf("%v", acc.GetID()), "token", token)
			if err != nil {
				return err
			}
		}
	}
	test, _ := s.Conn.HGetAll("accounts:1")
	log.Println(test)
	return nil
}

func (s *StorageRedis) Size() (int, error) {
	reply, err := s.Conn.Keys(fmt.Sprintf("%v:*", s.StorageAd))
	if err != nil {
		return 0, err
	}
	return len(reply), err
}
