package grpc

import (
	api "Mongo/api/protoc"
	"Mongo/internal/pkg/constErr"
	"Mongo/internal/pkg/data"
	"Mongo/internal/pkg/storage"
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)

type StorageGrpcServer struct {
	api.UnimplementedServiceProtobufServer
	Storage storage.InterfaceStorage
}

func NewStorageGrpcServer(s storage.InterfaceStorage) *StorageGrpcServer {
	return &StorageGrpcServer{
		Storage: s,
	}
}

func (s *StorageGrpcServer) Get(ctx context.Context, r *api.RequestID) (*api.ResponseGetAds, error) {
	ad, err := s.Storage.Get(uint(r.GetId()))
	if err != nil {
		return nil, err
	}

	return &api.ResponseGetAds{
		Ads: []*api.Ad{
			{
				Id:    uint32(ad.GetID()),
				Brand: ad.GetBrand(),
				Model: ad.GetModel(),
				Color: ad.GetColor(),
				Price: int32(ad.GetPrice()),
			},
		},
	}, nil
}

func (s *StorageGrpcServer) GetAll(ctx context.Context, e *empty.Empty) (*api.ResponseGetAds, error) {
	ads, err := s.Storage.GetAll()
	if err != nil {
		return nil, err
	}

	respAds := make([]*api.Ad, 0, len(ads))

	for _, v := range ads {
		respAds = append(respAds, &api.Ad{
			Id:    uint32(v.GetID()),
			Brand: v.GetBrand(),
			Model: v.GetModel(),
			Color: v.GetColor(),
			Price: int32(v.GetPrice()),
		})
	}

	return &api.ResponseGetAds{
		Ads: respAds,
	}, nil
}

func (s *StorageGrpcServer) Update(ctx context.Context, r *api.RequestUpdateAd) (*api.ResponseStatus, error) {
	ad := data.Ad{
		ID:    uint(r.GetId()),
		Brand: r.GetAd().GetBrand(),
		Model: r.GetAd().GetModel(),
		Color: r.GetAd().GetColor(),
		Price: int(r.GetAd().GetPrice()),
	}

	err := s.Storage.Update(&ad, uint(r.GetId()))
	if err != nil {
		return &api.ResponseStatus{Status: "False"}, err
	}

	return &api.ResponseStatus{Status: "True"}, nil
}

func (s *StorageGrpcServer) Add(ctx context.Context, r *api.RequestAdd) (*api.ResponseStatus, error) {
	ad := data.Ad{
		Brand: r.GetAd().GetBrand(),
		Model: r.GetAd().GetModel(),
		Color: r.GetAd().GetColor(),
		Price: int(r.GetAd().GetPrice()),
	}
	err := s.Storage.Add(&ad)
	if err != nil {
		return &api.ResponseStatus{Status: "False"}, err
	}

	return &api.ResponseStatus{Status: "True"}, nil
}

func (s *StorageGrpcServer) Delete(ctx context.Context, r *api.RequestID) (*api.ResponseStatus, error) {
	err := s.Storage.Delete(uint(r.GetId()))
	if err != nil {
		return &api.ResponseStatus{Status: "False"}, err
	}

	return &api.ResponseStatus{Status: "True"}, nil
}

func (s *StorageGrpcServer) Size(ctx context.Context, e *empty.Empty) (*api.ResponseSize, error) {
	size, err := s.Storage.Size()
	if err != nil {
		return nil, err
	}

	return &api.ResponseSize{Size: uint32(size)}, nil
}

func (s *StorageGrpcServer) AddAccount(ctx context.Context, r *api.Account) (*api.ResponseStatus, error) {
	acc := data.Account{
		Username: r.GetUsername(),
		Password: r.GetPassword(),
		Token:    r.GetToken(),
	}

	err := s.Storage.AddAccount(&acc)
	if err != nil {
		return &api.ResponseStatus{Status: "False"}, err
	}

	return &api.ResponseStatus{Status: "True"}, nil
}

func (s *StorageGrpcServer) GetAccounts(ctx context.Context, e *empty.Empty) (*api.ResponseAccounts, error) {
	baseAcc, err := s.Storage.GetAccounts()
	if err != nil {
		return nil, err
	}

	if baseAcc == nil || len(baseAcc) == 0 {
		return nil, constErr.AccountBaseIsEmpty
	}

	respAccs := make([]*api.Account, 0, len(baseAcc))

	for _, v := range baseAcc {
		respAccs = append(respAccs, &api.Account{
			Username: v.GetUserName(),
			Password: v.GetPassword(),
			Token:    v.GetToken(),
		})
	}

	return &api.ResponseAccounts{
		Acc: respAccs,
	}, nil
}

//ADD Update token
func (s *StorageGrpcServer) UpdateTokenCurrentAcc(ctx context.Context, r *api.RequestUpdateAcc) (*api.ResponseStatus, error) {
	acc := data.Account{
		Username: r.GetAcc().GetUsername(),
		Password: r.GetAcc().GetPassword(),
		Token:    r.GetToken(),
	}

	err := s.Storage.UpdateTokenCurrentAcc(&acc, r.GetToken())
	if err != nil {
		return &api.ResponseStatus{Status: "False"}, err
	}

	return &api.ResponseStatus{Status: "True"}, err

}
