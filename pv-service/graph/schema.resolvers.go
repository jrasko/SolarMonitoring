package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"pv-service/entities/dto"
	"pv-service/graph/generated"
	"pv-service/graph/model"
)

// DailyDataSets is the resolver for the dailyDataSets field.
func (r *queryResolver) DailyDataSets(ctx context.Context, begin *dto.PVTime, end *dto.PVTime, energyInterval uint32, startupInterval uint32) ([]*model.DailyData, error) {
	if energyInterval == 0 || startupInterval == 0 {
		return nil, errors.New("interval cannot be 0")
	}
	dailyData, err := r.Resolver.Processor.GetDailyData(ctx, begin, end, energyInterval, startupInterval)
	if err != nil {
		return nil, err
	}
	return dailyData, nil
}

// MinuteDataSets is the resolver for the MinuteDataSets field.
func (r *queryResolver) MinuteDataSets(ctx context.Context, begin *dto.PVTime, end *dto.PVTime, currentInterval uint32) ([]*model.MinuteDataOfDay, error) {
	minuteData, err := r.Resolver.Processor.GetMinuteDataOfDay(ctx, begin, end, currentInterval)
	if err != nil {
		return nil, err
	}
	return minuteData, nil
}

// RawDataSets is the resolver for the RawDataSets field.
func (r *queryResolver) RawDataSets(ctx context.Context, begin *dto.PVTime, end *dto.PVTime) ([]*model.RawData, error) {
	data, err := r.Resolver.Processor.GetRawDataBetweenDates(ctx, begin, end)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ZappiDataSets is the resolver for the ZappiDataSets field.
func (r *queryResolver) ZappiDataSets(ctx context.Context, begin *dto.PVTime, end *dto.PVTime) ([]*model.ZappiData, error) {
	data, err := r.Resolver.Processor.GetZappiDataBetweenDates(ctx, begin, end)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
