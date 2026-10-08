package food

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"

	"caloriecounter/api/internal/database/sqlc"
	"caloriecounter/api/internal/nutrition"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("food not found")

type Store interface {
	FindExact(context.Context, string, *string) (Food, error)
	FindSimilar(context.Context, string, float32) (Food, error)
	Create(context.Context, *nutrition.GeneratedFood) (Food, error)
}

type Repository struct{ pool *pgxpool.Pool }

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

func (r *Repository) FindExact(ctx context.Context, name string, state *string) (Food, error) {
	row, err := sqlc.New(r.pool).GetFoodByExactName(ctx, sqlc.GetFoodByExactNameParams{NormalizedName: name, PreparationState: state})
	if errors.Is(err, pgx.ErrNoRows) {
		return Food{}, ErrNotFound
	}
	if err != nil {
		return Food{}, err
	}
	return fromExact(row), nil
}
func (r *Repository) FindSimilar(ctx context.Context, name string, threshold float32) (Food, error) {
	row, err := sqlc.New(r.pool).GetFoodBySimilarName(ctx, sqlc.GetFoodBySimilarNameParams{NormalizedName: name, Similarity: formatThreshold(threshold)})
	if errors.Is(err, pgx.ErrNoRows) {
		return Food{}, ErrNotFound
	}
	if err != nil {
		return Food{}, err
	}
	return fromSimilar(row), nil
}
func (r *Repository) Create(ctx context.Context, generated *nutrition.GeneratedFood) (Food, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return Food{}, err
	}
	defer tx.Rollback(ctx)
	queries := sqlc.New(tx)
	id := uuid.New()
	normalized := NormalizeName(generated.CanonicalName)
	if existing, findErr := queries.GetFoodByExactName(ctx, sqlc.GetFoodByExactNameParams{NormalizedName: normalized, PreparationState: generated.PreparationState}); findErr == nil {
		return fromExact(existing), nil
	}
	created, err := queries.CreateFood(ctx, sqlc.CreateFoodParams{ID: uuidValue(id), CanonicalName: generated.CanonicalName, NormalizedName: normalized, DisplayName: generated.DisplayName, PreparationState: generated.PreparationState, Quantity: number(generated.Quantity), QuantityUnit: generated.QuantityUnit})
	if err != nil {
		return Food{}, err
	}
	payload, err := json.Marshal(generated.Nutrition)
	if err != nil {
		return Food{}, err
	}
	if err = queries.CreateFoodNutrition(ctx, sqlc.CreateFoodNutritionParams{FoodID: uuidValue(id), Column2: payload}); err != nil {
		return Food{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return Food{}, err
	}
	return fromCreated(created, generated.Nutrition), nil
}

func uuidValue(value uuid.UUID) pgtype.UUID { return pgtype.UUID{Bytes: value, Valid: true} }
func number(value float64) pgtype.Numeric {
	scaled := int64(value * 1000)
	return pgtype.Numeric{Int: big.NewInt(scaled), Exp: -3, Valid: true}
}
func formatThreshold(value float32) string { return formatFloat(float64(value)) }
func formatFloat(value float64) string     { return json.Number(stringify(value)).String() }
func stringify(value float64) []byte       { return []byte(big.NewFloat(value).Text('f', 4)) }

func fromCreated(row sqlc.Food, values nutrition.Values) Food {
	return Food{ID: uuid.UUID(row.ID.Bytes), CanonicalName: row.CanonicalName, DisplayName: row.DisplayName, Quantity: floatValue(row.Quantity), QuantityUnit: row.QuantityUnit, PreparationState: textPointer(row.PreparationState), NutritionBasis: row.NutritionBasis, Nutrition: values, CreatedAt: row.CreatedAt.Time, UpdatedAt: row.UpdatedAt.Time}
}
func fromExact(row sqlc.GetFoodByExactNameRow) Food {
	return Food{ID: uuid.UUID(row.ID.Bytes), CanonicalName: row.CanonicalName, DisplayName: row.DisplayName, Quantity: floatValue(row.Quantity), QuantityUnit: row.QuantityUnit, PreparationState: textPointer(row.PreparationState), NutritionBasis: row.NutritionBasis, Nutrition: nutrition.Values{Calories: &row.NCalories, Protein: &row.NProtein, Carbohydrates: &row.NCarbohydrates, Fat: &row.NFat, Fiber: &row.NFiber, Sugars: &row.NSugars, SaturatedFat: &row.NSaturatedFat, Sodium: &row.NSodium, Calcium: &row.NCalcium, Iron: &row.NIron, Magnesium: &row.NMagnesium, Phosphorus: &row.NPhosphorus, Potassium: &row.NPotassium, Zinc: &row.NZinc, Copper: &row.NCopper, Manganese: &row.NManganese, Selenium: &row.NSelenium, VitaminA: &row.NVitaminA, VitaminC: &row.NVitaminC, VitaminD: &row.NVitaminD, VitaminE: &row.NVitaminE, VitaminK: &row.NVitaminK, Thiamin: &row.NThiamin, Riboflavin: &row.NRiboflavin, Niacin: &row.NNiacin, PantothenicAcid: &row.NPantothenicAcid, VitaminB6: &row.NVitaminB6, Biotin: &row.NBiotin, Folate: &row.NFolate, VitaminB12: &row.NVitaminB12, Caffeine: &row.NCaffeine, Cholesterol: &row.NCholesterol}, CreatedAt: row.CreatedAt.Time, UpdatedAt: row.UpdatedAt.Time}
}
func fromSimilar(row sqlc.GetFoodBySimilarNameRow) Food {
	return Food{ID: uuid.UUID(row.ID.Bytes), CanonicalName: row.CanonicalName, DisplayName: row.DisplayName, Quantity: floatValue(row.Quantity), QuantityUnit: row.QuantityUnit, PreparationState: textPointer(row.PreparationState), NutritionBasis: row.NutritionBasis, CreatedAt: row.CreatedAt.Time, UpdatedAt: row.UpdatedAt.Time, Nutrition: nutrition.Values{Calories: &row.NCalories, Protein: &row.NProtein, Carbohydrates: &row.NCarbohydrates, Fat: &row.NFat, Fiber: &row.NFiber, Sugars: &row.NSugars, SaturatedFat: &row.NSaturatedFat, Sodium: &row.NSodium, Calcium: &row.NCalcium, Iron: &row.NIron, Magnesium: &row.NMagnesium, Phosphorus: &row.NPhosphorus, Potassium: &row.NPotassium, Zinc: &row.NZinc, Copper: &row.NCopper, Manganese: &row.NManganese, Selenium: &row.NSelenium, VitaminA: &row.NVitaminA, VitaminC: &row.NVitaminC, VitaminD: &row.NVitaminD, VitaminE: &row.NVitaminE, VitaminK: &row.NVitaminK, Thiamin: &row.NThiamin, Riboflavin: &row.NRiboflavin, Niacin: &row.NNiacin, PantothenicAcid: &row.NPantothenicAcid, VitaminB6: &row.NVitaminB6, Biotin: &row.NBiotin, Folate: &row.NFolate, VitaminB12: &row.NVitaminB12, Caffeine: &row.NCaffeine, Cholesterol: &row.NCholesterol}}
}
func floatValue(value pgtype.Numeric) float64 {
	result, _ := value.Float64Value()
	return result.Float64
}
func textPointer(value *string) *string { return value }
