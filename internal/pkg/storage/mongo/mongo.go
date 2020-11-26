package mongo

import (
	"Mongo/internal/pkg/constErr"
	"Mongo/internal/pkg/data"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoStorage struct {
	CollectionAccs *mongo.Collection // Accounts
	CollectionAds  *mongo.Collection // Ads
}

func NewMongoStorage() *MongoStorage {
	uri := "mongodb://mongo:pass@mongo:8004"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("check ping")
	for err := client.Ping(ctx, readpref.Primary()); err != nil; err = client.Ping(ctx, readpref.Primary()) {
		time.Sleep(2 * time.Second)
	}
	log.Println("ping success.")

	quickstartDatabase := client.Database("mongo_vlad")

	quickstartDatabase.CreateCollection(ctx, "accounts")
	quickstartDatabase.CreateCollection(ctx, "ads")

	colAccs := quickstartDatabase.Collection("accounts")
	colAds := quickstartDatabase.Collection("ads")

	storage := &MongoStorage{
		CollectionAccs: colAccs,
		CollectionAds:  colAds,
	}

	return storage
}

func (s *MongoStorage) Add(ad *data.Ad) error {
	size, _ := s.Size()
	ad.ID = uint(size) + 1
	_, err := s.CollectionAds.InsertOne(context.Background(), ad)
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoStorage) Get(id uint) (*data.Ad, error) {
	size, err := s.Size()
	if err != nil || size == 0 {
		return nil, err
	}
	filter := bson.M{"id": id}
	result := s.CollectionAds.FindOne(context.Background(), filter)

	ad := new(data.Ad)
	if err := result.Decode(&ad); err != nil {
		return nil, err
	}

	return ad, nil
}

func (s *MongoStorage) GetAll() ([]*data.Ad, error) {
	size, err := s.Size()
	if err != nil || size == 0 {
		return nil, err
	}

	ads := make([]*data.Ad, 0, int(size))
	option := options.Find().SetSort(bson.M{"id": 1})
	cursor, err := s.CollectionAds.Find(context.Background(), bson.D{}, option)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		ad := new(data.Ad)
		if err := cursor.Decode(&ad); err != nil {
			return nil, err
		}
		ads = append(ads, ad)
	}

	return ads, nil
}

func (s *MongoStorage) Update(temp *data.Ad, id uint) error {
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{
		"brand": temp.GetBrand(),
		"model": temp.GetModel(),
		"color": temp.GetColor(),
		"price": temp.GetPrice()}}
	result := s.CollectionAds.FindOneAndUpdate(context.Background(), filter, update)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func (s *MongoStorage) Delete(id uint) error {
	size, err := s.Size()
	if err != nil || id > uint(size) {
		return err
	}
	filter := bson.M{"id": id}

	if _, err := s.CollectionAds.DeleteOne(context.Background(), filter); err != nil {
		return err
	}

	if size-1 == 0 {
		return nil
	} else {
		count := int(size - 1)
		newID := 1
		for i := 0; true; i++ {
			filter := bson.M{"id": i}
			update := bson.M{"$set": bson.M{"id": newID}}
			result := s.CollectionAds.FindOneAndUpdate(context.Background(), filter, update)
			if result.Err() != nil {
				continue
			}
			if newID == count {
				break
			}
			newID++
			log.Println(result.Err())
		}
	}

	return nil
}

func (s *MongoStorage) AddAccount(acc *data.Account) error {
	size, _ := s.CollectionAccs.CountDocuments(context.Background(), bson.D{})
	acc.ID = int(size) + 1
	_, err := s.CollectionAccs.InsertOne(context.Background(), acc)
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoStorage) GetAccounts() ([]*data.Account, error) {
	size, err := s.CollectionAccs.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		return nil, constErr.AdBaseIsEmpty
	}

	accs := make([]*data.Account, 0, int(size))

	cursor, err := s.CollectionAccs.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		acc := new(data.Account)
		if err := cursor.Decode(&acc); err != nil {
			return nil, err
		}

		accs = append(accs, acc)
	}

	return accs, nil
}

func (s *MongoStorage) UpdateTokenCurrentAcc(acc *data.Account, token string) error {
	filter := bson.M{"username": acc.GetUserName(), "password": acc.GetPassword()}
	update := bson.M{"$set": bson.M{"token": token}}
	result := s.CollectionAccs.FindOneAndUpdate(context.Background(), filter, update)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (s *MongoStorage) Size() (int, error) {
	size, err := s.CollectionAds.CountDocuments(context.Background(), bson.D{})
	if err != nil || size == 0 {
		return 0, constErr.AdBaseIsEmpty
	}
	return int(size), nil
}
