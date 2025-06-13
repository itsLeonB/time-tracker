package service

import (
	"context"

	"github.com/itsLeonB/time-tracker/internal/dto"
	"github.com/itsLeonB/time-tracker/internal/mapper"
	"github.com/itsLeonB/time-tracker/internal/model"
	"github.com/itsLeonB/time-tracker/internal/repository"
	"github.com/itsLeonB/time-tracker/internal/util"
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

func (ys *youtrackService) FindTask(ctx context.Context, queryOptions dto.ExternalQueryOptions) ([]model.ExternalTask, error) {
	youtrackQueryOpts := dto.YoutrackQueryOptions{
		IdReadable: queryOptions.Number,
	}

	tasks, err := ys.youtrackRepository.FindTask(ctx, youtrackQueryOpts)
	if err != nil {
		return nil, err
	}

	return util.MapSlice(tasks, mapper.YoutrackToExternal), nil
}
