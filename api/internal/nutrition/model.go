package nutrition

type Values struct {
	Calories        *float64 `json:"calories"`
	Protein         *float64 `json:"protein"`
	Carbohydrates   *float64 `json:"carbohydrates"`
	Fat             *float64 `json:"fat"`
	Fiber           *float64 `json:"fiber"`
	Sugars          *float64 `json:"sugars"`
	SaturatedFat    *float64 `json:"saturated_fat"`
	Sodium          *float64 `json:"sodium"`
	Calcium         *float64 `json:"calcium"`
	Iron            *float64 `json:"iron"`
	Magnesium       *float64 `json:"magnesium"`
	Phosphorus      *float64 `json:"phosphorus"`
	Potassium       *float64 `json:"potassium"`
	Zinc            *float64 `json:"zinc"`
	Copper          *float64 `json:"copper"`
	Manganese       *float64 `json:"manganese"`
	Selenium        *float64 `json:"selenium"`
	VitaminA        *float64 `json:"vitamin_a"`
	VitaminC        *float64 `json:"vitamin_c"`
	VitaminD        *float64 `json:"vitamin_d"`
	VitaminE        *float64 `json:"vitamin_e"`
	VitaminK        *float64 `json:"vitamin_k"`
	Thiamin         *float64 `json:"thiamin"`
	Riboflavin      *float64 `json:"riboflavin"`
	Niacin          *float64 `json:"niacin"`
	PantothenicAcid *float64 `json:"pantothenic_acid"`
	VitaminB6       *float64 `json:"vitamin_b6"`
	Biotin          *float64 `json:"biotin"`
	Folate          *float64 `json:"folate"`
	VitaminB12      *float64 `json:"vitamin_b12"`
	Caffeine        *float64 `json:"caffeine"`
	Cholesterol     *float64 `json:"cholesterol"`
}

type Request struct {
	Query    string
	Quantity float64
	Unit     string
}
type GeneratedFood struct {
	CanonicalName    string  `json:"canonical_name"`
	DisplayName      string  `json:"display_name"`
	Quantity         float64 `json:"quantity"`
	QuantityUnit     string  `json:"quantity_unit"`
	PreparationState *string `json:"preparation_state"`
	NutritionBasis   string  `json:"nutrition_basis"`
	Nutrition        Values  `json:"nutrition"`
}
