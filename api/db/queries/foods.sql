-- name: GetFoodByExactName :one
SELECT f.*, n.food_id, n.calories::float8, n.protein::float8, n.carbohydrates::float8, n.fat::float8, n.fiber::float8, n.sugars::float8, n.saturated_fat::float8, n.sodium::float8, n.calcium::float8, n.iron::float8, n.magnesium::float8, n.phosphorus::float8, n.potassium::float8, n.zinc::float8, n.copper::float8, n.manganese::float8, n.selenium::float8, n.vitamin_a::float8, n.vitamin_c::float8, n.vitamin_d::float8, n.vitamin_e::float8, n.vitamin_k::float8, n.thiamin::float8, n.riboflavin::float8, n.niacin::float8, n.pantothenic_acid::float8, n.vitamin_b6::float8, n.biotin::float8, n.folate::float8, n.vitamin_b12::float8, n.caffeine::float8, n.cholesterol::float8, n.created_at, n.updated_at FROM foods f JOIN food_nutrition n ON n.food_id = f.id
WHERE f.normalized_name = $1 AND f.preparation_state IS NOT DISTINCT FROM $2 LIMIT 1;

-- name: GetFoodBySimilarName :one
SELECT f.*, n.food_id, n.calories::float8, n.protein::float8, n.carbohydrates::float8, n.fat::float8, n.fiber::float8, n.sugars::float8, n.saturated_fat::float8, n.sodium::float8, n.calcium::float8, n.iron::float8, n.magnesium::float8, n.phosphorus::float8, n.potassium::float8, n.zinc::float8, n.copper::float8, n.manganese::float8, n.selenium::float8, n.vitamin_a::float8, n.vitamin_c::float8, n.vitamin_d::float8, n.vitamin_e::float8, n.vitamin_k::float8, n.thiamin::float8, n.riboflavin::float8, n.niacin::float8, n.pantothenic_acid::float8, n.vitamin_b6::float8, n.biotin::float8, n.folate::float8, n.vitamin_b12::float8, n.caffeine::float8, n.cholesterol::float8, n.created_at, n.updated_at, similarity(f.normalized_name, $1) AS match_score
FROM foods f JOIN food_nutrition n ON n.food_id = f.id
WHERE similarity(f.normalized_name, $1) >= $2
ORDER BY similarity(f.normalized_name, $1) DESC LIMIT 1;

-- name: CreateFood :one
INSERT INTO foods (id, canonical_name, normalized_name, display_name, preparation_state, quantity, quantity_unit, nutrition_basis, provider_metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, 'per_100g', $8) RETURNING *;

-- name: CreateFoodNutrition :exec
INSERT INTO food_nutrition (food_id, calories, protein, carbohydrates, fat, fiber, sugars, saturated_fat, sodium, calcium, iron, magnesium, phosphorus, potassium, zinc, copper, manganese, selenium, vitamin_a, vitamin_c, vitamin_d, vitamin_e, vitamin_k, thiamin, riboflavin, niacin, pantothenic_acid, vitamin_b6, biotin, folate, vitamin_b12, caffeine, cholesterol)
SELECT $1, n.* FROM jsonb_to_record($2::jsonb) AS n(
 calories numeric, protein numeric, carbohydrates numeric, fat numeric, fiber numeric, sugars numeric, saturated_fat numeric,
 sodium numeric, calcium numeric, iron numeric, magnesium numeric, phosphorus numeric, potassium numeric, zinc numeric, copper numeric,
 manganese numeric, selenium numeric, vitamin_a numeric, vitamin_c numeric, vitamin_d numeric, vitamin_e numeric, vitamin_k numeric,
 thiamin numeric, riboflavin numeric, niacin numeric, pantothenic_acid numeric, vitamin_b6 numeric, biotin numeric, folate numeric,
 vitamin_b12 numeric, caffeine numeric, cholesterol numeric
);