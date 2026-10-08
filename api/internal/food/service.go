package food

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"caloriecounter/api/internal/nutrition"
)

var ErrValidation = errors.New("invalid food search request")

type Service struct {
	store               Store
	provider            nutrition.Provider
	SimilarityThreshold float32
}

func NewService(store Store, provider nutrition.Provider, threshold float32) *Service {
	return &Service{store: store, provider: provider, SimilarityThreshold: threshold}
}
func (s *Service) Search(ctx context.Context, request nutrition.Request) (Food, error) {
	if strings.TrimSpace(request.Query) == "" || request.Quantity <= 0 || request.Unit == "" {
		slog.Warn("food search rejected by validation", "query", request.Query)
		return Food{}, ErrValidation
	}
	name := NormalizeName(request.Query)
	state := preparationFromQuery(name)
	if result, err := s.store.FindExact(ctx, name, state); err == nil {
		slog.Debug("food found by exact database match", "normalized_name", name)
		result.Quantity = request.Quantity
		result.QuantityUnit = request.Unit
		return result, nil
	}
	if result, err := s.store.FindSimilar(ctx, name, s.SimilarityThreshold); err == nil {
		slog.Debug("food found by fuzzy database match", "normalized_name", name, "similarity_threshold", s.SimilarityThreshold)
		result.Quantity = request.Quantity
		result.QuantityUnit = request.Unit
		return result, nil
	}
	generated, err := s.provider.GenerateNutrition(ctx, request)
	if err != nil {
		slog.Error("nutrition provider failed", "query", request.Query, "error", err)
		return Food{}, err
	}
	if err = nutrition.Validate(generated); err != nil {
		slog.Error("generated nutrition failed validation", "query", request.Query, "error", err)
		return Food{}, err
	}
	result, err := s.store.Create(ctx, generated)
	if err != nil {
		slog.Error("generated food could not be stored", "query", request.Query, "error", err)
		return Food{}, err
	}
	result.Quantity = request.Quantity
	result.QuantityUnit = request.Unit
	slog.Info("food generated and stored", "query", request.Query, "food_id", result.ID)
	return result, nil
}
func preparationFromQuery(query string) *string {
	for _, state := range []string{"raw", "boiled", "fried", "baked", "grilled", "roasted", "steamed", "poached", "sauteed", "dried", "canned", "frozen"} {
		if strings.Contains(query, state) {
			value := state
			return &value
		}
	}
	return nil
}
