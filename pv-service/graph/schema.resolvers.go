package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"pv-service/graph/generated"
	"pv-service/graph/model"
	"time"
)

func (r *queryResolver) DailyDataSets(_ context.Context, begin *time.Time, end *time.Time, energyInterval uint32, startupInterval uint32) ([]*model.DailyData, error) {
	if energyInterval == 0 || startupInterval == 0 {
		return nil, errors.New("interval cannot be 0")
	}
	dailyData, err := r.Resolver.Processor.GetDailyData(begin, end, energyInterval, startupInterval)
	if err != nil {
		return nil, err
	}
	return dailyData, nil
}

func (r *queryResolver) MinuteDataSets(_ context.Context, begin *time.Time, end *time.Time) ([]*model.MinuteDataOfDay, error) {
	minuteData, err := r.Resolver.Processor.GetMinuteDataOfDay(begin, end)
	if err != nil {
		return nil, err
	}
	return minuteData, nil
}

func (r *queryResolver) RawDataSets(_ context.Context, begin *time.Time, end *time.Time) ([]*model.RawData, error) {
	data, err := r.Resolver.Processor.GetRawDataBetweenDates(begin, end)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
