package strategy

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/itsLeonB/time-tracker/internal/apperror"
	"github.com/itsLeonB/time-tracker/internal/dto"
)

type ErrorStrategy interface {
	HandleError(err error) *dto.ErrorResponse
}

type ErrorStrategyMap struct {
	Strategies      map[reflect.Type]ErrorStrategy
	DefaultStrategy ErrorStrategy
}

func (esm *ErrorStrategyMap) DetermineStrategy(err error) ErrorStrategy {
	if strategy, ok := esm.Strategies[reflect.TypeOf(err)]; ok {
		return strategy
	}

	return esm.DefaultStrategy
}

func NewErrorStrategyMap() *ErrorStrategyMap {
	return &ErrorStrategyMap{
		Strategies: map[reflect.Type]ErrorStrategy{
			reflect.TypeOf(&apperror.AppError{}):         &appErrorStrategy{},
			reflect.TypeOf(validator.ValidationErrors{}): &validationErrorStrategy{},
		},
		DefaultStrategy: &unhandledErrorStrategy{},
	}
}
