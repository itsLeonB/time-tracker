package service

import (
	"context"

	"github.com/itsLeonB/catfeinated-time-tracker/internal/mapper"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/repository"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/util"
)

type youtrackService struct {
	youtrackRepository *repository.YoutrackRepository
}

func NewYoutrackService(
	youtrackRepository *repository.YoutrackRepository,
) ExternalTrackerService {
	return &youtrackService{
		youtrackRepository,
	}
}

func (ys *youtrackService) FindTask(ctx context.Context, queryOptions model.ExternalQueryOptions) ([]model.ExternalTask, error) {
	youtrackQueryOpts := model.YoutrackQueryOptions{
		IdReadable: queryOptions.Number,
	}

	tasks, err := ys.youtrackRepository.FindTask(ctx, youtrackQueryOpts)
	if err != nil {
		return nil, err
	}

	return util.MapSlice(tasks, mapper.YoutrackToExternal), nil
}
