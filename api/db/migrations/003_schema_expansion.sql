CREATE EXTENSION IF NOT EXISTS btree_gist;

BEGIN;

CREATE TABLE IF NOT EXISTS food_catalog_quantity_backup (
    food_id UUID PRIMARY KEY REFERENCES foods(id) ON DELETE CASCADE,
    quantity NUMERIC(12, 3),
    quantity_unit TEXT,
    migrated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO food_catalog_quantity_backup (food_id, quantity, quantity_unit)
SELECT id, quantity, quantity_unit
FROM foods
WHERE quantity IS NOT NULL
ON CONFLICT (food_id) DO NOTHING;

ALTER TABLE users ADD COLUMN IF NOT EXISTS display_name TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS date_of_birth DATE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS age INTEGER;
ALTER TABLE users ADD COLUMN IF NOT EXISTS height_cm NUMERIC(8, 2);
ALTER TABLE users ADD COLUMN IF NOT EXISTS sex TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS activity_level TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS timezone TEXT NOT NULL DEFAULT 'Asia/Kolkata';
ALTER TABLE users ADD COLUMN IF NOT EXISTS maintenance_calories NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_calories NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_protein_g_day NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_carbohydrates_g_day NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_fat_g_day NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_fiber_g_day NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_sugars_g_day NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_saturated_fat_g_day NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_calcium_mg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_iron_mg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_magnesium_mg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_phosphorus_mg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_potassium_mg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_zinc_mg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_copper_mg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_manganese_mg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_selenium_mcg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_vitamin_a_mcg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_vitamin_c_mg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_vitamin_d_mcg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_vitamin_e_mg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_vitamin_k_mcg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_thiamin_mg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_riboflavin_mg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_niacin_mg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_pantothenic_acid_mg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_vitamin_b6_mg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_biotin_mcg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_folate_mcg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_vitamin_b12_mcg_week NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_sodium_mg_day NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_cholesterol_mg_day NUMERIC(12, 3);
ALTER TABLE users ADD COLUMN IF NOT EXISTS recommended_caffeine_mg_day NUMERIC(12, 3);

ALTER TABLE foods DROP COLUMN IF EXISTS quantity;
ALTER TABLE foods DROP COLUMN IF EXISTS quantity_unit;

ALTER TABLE foods
    ADD COLUMN IF NOT EXISTS nutrition_basis TEXT;
UPDATE foods SET nutrition_basis = 'per_100g' WHERE nutrition_basis IS NULL;
ALTER TABLE foods ALTER COLUMN nutrition_basis SET DEFAULT 'per_100g';
ALTER TABLE foods ALTER COLUMN nutrition_basis SET NOT NULL;

CREATE TABLE IF NOT EXISTS user_weight_history (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    weight_kg NUMERIC(12, 3) NOT NULL CHECK (weight_kg > 0),
    effective_from DATE NOT NULL,
    effective_to DATE,
    effective_range DATERANGE GENERATED ALWAYS AS (daterange(effective_from, COALESCE(effective_to, 'infinity'::date), '[)')) STORED,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (effective_to IS NULL OR effective_from < effective_to)
);
CREATE UNIQUE INDEX IF NOT EXISTS user_weight_history_current_idx ON user_weight_history (user_id) WHERE effective_to IS NULL;
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conrelid = 'user_weight_history'::regclass AND conname = 'user_weight_history_user_id_effective_range_excl'
    ) THEN
        ALTER TABLE user_weight_history ADD EXCLUDE USING gist (user_id WITH =, effective_range WITH &&);
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS user_diet_plan_history (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    diet_plan SMALLINT NOT NULL CHECK (diet_plan BETWEEN 1 AND 5),
    calorie_adjustment_percent NUMERIC(12, 3) NOT NULL,
    effective_from DATE NOT NULL,
    effective_to DATE,
    effective_range DATERANGE GENERATED ALWAYS AS (daterange(effective_from, COALESCE(effective_to, 'infinity'::date), '[)')) STORED,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (effective_to IS NULL OR effective_from < effective_to)
);
CREATE UNIQUE INDEX IF NOT EXISTS user_diet_plan_history_current_idx ON user_diet_plan_history (user_id) WHERE effective_to IS NULL;
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conrelid = 'user_diet_plan_history'::regclass AND conname = 'user_diet_plan_history_user_id_effective_range_excl'
    ) THEN
        ALTER TABLE user_diet_plan_history ADD EXCLUDE USING gist (user_id WITH =, effective_range WITH &&);
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS meals (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    meal_type TEXT NOT NULL CHECK (meal_type IN ('breakfast', 'brunch', 'lunch', 'dinner', 'snack')),
    log_date DATE NOT NULL,
    logged_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS meals_user_date_idx ON meals (user_id, log_date);

CREATE TABLE IF NOT EXISTS meal_items (
    id UUID PRIMARY KEY,
    meal_id UUID NOT NULL REFERENCES meals(id) ON DELETE CASCADE,
    food_id UUID REFERENCES foods(id) ON DELETE SET NULL,
    display_name TEXT NOT NULL,
    quantity NUMERIC(12, 3) NOT NULL CHECK (quantity > 0),
    quantity_unit TEXT NOT NULL CHECK (quantity_unit IN ('g', 'ml', 'unit')),
    weight_grams NUMERIC(12, 3) CHECK (weight_grams IS NULL OR weight_grams > 0),
    nutrition_snapshot JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS daily_summaries (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    log_date DATE NOT NULL,
    maintenance_calories NUMERIC(12, 3) NOT NULL DEFAULT 0 CHECK (maintenance_calories >= 0),
    recommended_calories NUMERIC(12, 3) NOT NULL DEFAULT 0 CHECK (recommended_calories >= 0),
    cheat_day BOOLEAN NOT NULL DEFAULT FALSE,
    calories_consumed NUMERIC(12, 3) NOT NULL DEFAULT 0 CHECK (calories_consumed >= 0),
    protein_g NUMERIC(12, 3) NOT NULL DEFAULT 0 CHECK (protein_g >= 0),
    carbohydrates_g NUMERIC(12, 3) NOT NULL DEFAULT 0 CHECK (carbohydrates_g >= 0),
    fat_g NUMERIC(12, 3) NOT NULL DEFAULT 0 CHECK (fat_g >= 0),
    energy_balance_kcal NUMERIC(12, 3) NOT NULL DEFAULT 0,
    target_variance_kcal NUMERIC(12, 3) NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, log_date)
);

COMMIT;
