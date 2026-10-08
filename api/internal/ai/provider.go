package ai

import (
	"caloriecounter/api/internal/nutrition"
	"context"
	"errors"
)

var ErrUnavailable = errors.New("nutrition provider unavailable")

type Provider struct{}

func (Provider) GenerateNutrition(context.Context, nutrition.Request) (*nutrition.GeneratedFood, error) {
	return nil, ErrUnavailable
}
