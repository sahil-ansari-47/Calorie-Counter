package nutrition

import (
	"errors"
	"math"
	"strings"
)

var ErrInvalidNutrition = errors.New("invalid nutrition data")
var validStates = map[string]bool{"raw": true, "boiled": true, "fried": true, "baked": true, "grilled": true, "roasted": true, "steamed": true, "poached": true, "sauteed": true, "dried": true, "canned": true, "frozen": true}

func Validate(food *GeneratedFood) error {
	if food == nil || strings.TrimSpace(food.CanonicalName) == "" || strings.TrimSpace(food.DisplayName) == "" || food.Quantity <= 0 || food.NutritionBasis != "per_100g" {
		return ErrInvalidNutrition
	}
	if food.QuantityUnit != "g" && food.QuantityUnit != "ml" && food.QuantityUnit != "oz" && food.QuantityUnit != "serving" {
		return ErrInvalidNutrition
	}
	if food.PreparationState != nil && !validStates[*food.PreparationState] {
		return ErrInvalidNutrition
	}
	for _, value := range nutrientValues(food.Nutrition) {
		if value != nil && (math.IsNaN(*value) || math.IsInf(*value, 0) || *value < 0 || *value > 100000) {
			return ErrInvalidNutrition
		}
	}
	return nil
}

func nutrientValues(v Values) []*float64 {
	return []*float64{v.Calories, v.Protein, v.Carbohydrates, v.Fat, v.Fiber, v.Sugars, v.SaturatedFat, v.Sodium, v.Calcium, v.Iron, v.Magnesium, v.Phosphorus, v.Potassium, v.Zinc, v.Copper, v.Manganese, v.Selenium, v.VitaminA, v.VitaminC, v.VitaminD, v.VitaminE, v.VitaminK, v.Thiamin, v.Riboflavin, v.Niacin, v.PantothenicAcid, v.VitaminB6, v.Biotin, v.Folate, v.VitaminB12, v.Caffeine, v.Cholesterol}
}
