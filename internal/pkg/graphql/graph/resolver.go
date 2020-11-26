package graph

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"Mongo/internal/pkg/constErr"
	"Mongo/internal/pkg/data"
	"Mongo/internal/pkg/graphql/graph/generated"
	"Mongo/internal/pkg/graphql/graph/model"
	"Mongo/internal/pkg/service"
	"context"
	"fmt"
	"strconv"
)

type Resolver struct {
	IGraphqL service.InterfaceServer
}

func (r *mutationResolver) CreateAd(ctx context.Context, ad *model.AdRequest) (string, error) {
	var ans string
	adv := data.Ad{
		Brand: ad.Brand,
		Model: ad.Model,
		Color: ad.Color,
		Price: ad.Price,
	}

	err := r.IGraphqL.Add(&adv)
	if err != nil {
		ans = "No create"
	} else {
		ans = "Create"
	}
	return ans, err
}

func (r *mutationResolver) UpdateAd(ctx context.Context, ad *model.AdRequest, id string) (string, error) {
	ID, err := strconv.Atoi(id)
	if err != nil {
		return "Invalid ID", err
	}
	adv := data.Ad{
		Brand: ad.Brand,
		Model: ad.Model,
		Color: ad.Color,
		Price: ad.Price,
	}
	err = r.IGraphqL.Update(uint(ID), &adv)
	if err != nil {
		return "No Update", err
	}

	return "Update", err
}

func (r *mutationResolver) DeleteAd(ctx context.Context, id string) (string, error) {
	ID, err := strconv.Atoi(id)
	if err != nil {
		return "Invalid ID", err
	}

	err = r.IGraphqL.Delete(uint(ID))
	return "Delete", err
}

func (r *queryResolver) Getall(ctx context.Context) ([]*model.Ad, error) {
	ads, err := r.IGraphqL.GetAll()
	if err != nil {
		return nil, err
	}

	adsGQL := make([]*model.Ad, 0, len(ads))

	for _, v := range ads {
		adsGQL = append(adsGQL, &model.Ad{
			ID:    fmt.Sprintf("%v", v.GetID()),
			Brand: v.GetBrand(),
			Model: v.GetModel(),
			Color: v.GetColor(),
			Price: v.GetPrice(),
		})
	}

	return adsGQL, err
}

func (r *queryResolver) Get(ctx context.Context, id string) (*model.Ad, error) {
	ID, err := strconv.Atoi(id)
	if err != nil {
		return nil, constErr.InvalidID
	}
	ad, err := r.IGraphqL.Get(uint(ID))
	if err != nil {
		return nil, err
	}

	adGQL := model.Ad{
		ID:    fmt.Sprintf("%v", ad.GetID()),
		Brand: ad.GetBrand(),
		Model: ad.GetModel(),
		Color: ad.GetColor(),
		Price: ad.GetPrice(),
	}
	return &adGQL, err
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
