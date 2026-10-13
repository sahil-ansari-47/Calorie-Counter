-- name: CreateUser :one
INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: UpdateUserProfile :one
UPDATE users
SET
display_name = $2,
    date_of_birth = $3,
    age = $4,
    height_cm = $5,
    sex = $6,
    activity_level = $7,
    timezone = COALESCE(NULLIF($8, ''), 'Asia/Kolkata'),
    maintenance_calories = $9,
    recommended_calories = $10,
    recommended_protein_g_day = $11,
    recommended_carbohydrates_g_day = $12,
    recommended_fat_g_day = $13,
    recommended_fiber_g_day = $14,
    recommended_sugars_g_day = $15,
    recommended_saturated_fat_g_day = $16,
    recommended_calcium_mg_week = $17,
    recommended_iron_mg_week = $18,
    recommended_magnesium_mg_week = $19,
    recommended_phosphorus_mg_week = $20,
    recommended_potassium_mg_week = $21,
    recommended_zinc_mg_week = $22,
    recommended_copper_mg_week = $23,
    recommended_manganese_mg_week = $24,
    recommended_selenium_mcg_week = $25,
    recommended_vitamin_a_mcg_week = $26,
    recommended_vitamin_c_mg_week = $27,
    recommended_vitamin_d_mcg_week = $28,
    recommended_vitamin_e_mg_week = $29,
    recommended_vitamin_k_mcg_week = $30,
    recommended_thiamin_mg_week = $31,
    recommended_riboflavin_mg_week = $32,
    recommended_niacin_mg_week = $33,
    recommended_pantothenic_acid_mg_week = $34,
    recommended_vitamin_b6_mg_week = $35,
    recommended_biotin_mcg_week = $36,
    recommended_folate_mcg_week = $37,
    recommended_vitamin_b12_mcg_week = $38,
    recommended_sodium_mg_day = $39,
    recommended_cholesterol_mg_day = $40,
    recommended_caffeine_mg_day = $41,
    updated_at = now()
WHERE id = $1
RETURNING *;