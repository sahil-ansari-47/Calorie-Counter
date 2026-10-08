package food

import (
	"regexp"
	"strings"
	"time"

	"caloriecounter/api/internal/nutrition"

	"github.com/google/uuid"
)

type Food struct {
	ID               uuid.UUID        `json:"id"`
	CanonicalName    string           `json:"canonical_name"`
	DisplayName      string           `json:"display_name"`
	Quantity         float64          `json:"quantity"`
	QuantityUnit     string           `json:"quantity_unit"`
	PreparationState *string          `json:"preparation_state,omitempty"`
	NutritionBasis   string           `json:"nutrition_basis"`
	Nutrition        nutrition.Values `json:"nutrition"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

var punctuation = regexp.MustCompile(`[,:;()\[\]{}]+`)

func NormalizeName(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = punctuation.ReplaceAllString(value, " ")
	return strings.Join(strings.Fields(value), " ")
}
