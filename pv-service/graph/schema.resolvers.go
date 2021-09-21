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
	data, err := r.Resolver.Processor.GetDailyDataBetweenDates(begin, end, energyInterval, startupInterval)
	if err != nil {
		return nil, err
	}
	return data, nil
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
