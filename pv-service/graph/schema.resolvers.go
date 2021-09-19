package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"pv-service/graph/generated"
	"pv-service/graph/model"
	"time"
)

func (r *queryResolver) DailyDataSets(ctx context.Context, begin *time.Time, end *time.Time) ([]*model.DailyData, error) {
	data, err := r.Resolver.Processor.GetDailyDataBetweenDates(begin, end)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *queryResolver) RawDataSets(ctx context.Context, begin *time.Time, end *time.Time) ([]*model.RawData, error) {
	data, err := r.Resolver.Processor.GetRawDataBetweenDates(begin, end)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
