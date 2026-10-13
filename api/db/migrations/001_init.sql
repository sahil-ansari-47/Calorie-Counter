CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    display_name TEXT,
    date_of_birth DATE,
    age INTEGER CHECK (age IS NULL OR age > 0),
    height_cm NUMERIC(8, 2) CHECK (height_cm IS NULL OR height_cm > 0),
    sex TEXT CHECK (sex IN ('male', 'female', 'other', 'prefer_not_to_say')),
    activity_level TEXT CHECK (activity_level IN ('sedentary', 'light', 'moderate', 'active', 'very_active')),
    timezone TEXT NOT NULL DEFAULT 'Asia/Kolkata',
    maintenance_calories NUMERIC(12, 3) CHECK (maintenance_calories IS NULL OR maintenance_calories > 0),
    recommended_calories NUMERIC(12, 3) CHECK (recommended_calories IS NULL OR recommended_calories > 0),
    recommended_protein_g_day NUMERIC(12, 3) CHECK (recommended_protein_g_day IS NULL OR recommended_protein_g_day > 0),
    recommended_carbohydrates_g_day NUMERIC(12, 3) CHECK (recommended_carbohydrates_g_day IS NULL OR recommended_carbohydrates_g_day > 0),
    recommended_fat_g_day NUMERIC(12, 3) CHECK (recommended_fat_g_day IS NULL OR recommended_fat_g_day > 0),
    recommended_fiber_g_day NUMERIC(12, 3) CHECK (recommended_fiber_g_day IS NULL OR recommended_fiber_g_day > 0),
    recommended_sugars_g_day NUMERIC(12, 3) CHECK (recommended_sugars_g_day IS NULL OR recommended_sugars_g_day > 0),
    recommended_saturated_fat_g_day NUMERIC(12, 3) CHECK (recommended_saturated_fat_g_day IS NULL OR recommended_saturated_fat_g_day > 0),
    recommended_calcium_mg_week NUMERIC(12, 3) CHECK (recommended_calcium_mg_week IS NULL OR recommended_calcium_mg_week > 0),
    recommended_iron_mg_week NUMERIC(12, 3) CHECK (recommended_iron_mg_week IS NULL OR recommended_iron_mg_week > 0),
    recommended_magnesium_mg_week NUMERIC(12, 3) CHECK (recommended_magnesium_mg_week IS NULL OR recommended_magnesium_mg_week > 0),
    recommended_phosphorus_mg_week NUMERIC(12, 3) CHECK (recommended_phosphorus_mg_week IS NULL OR recommended_phosphorus_mg_week > 0),
    recommended_potassium_mg_week NUMERIC(12, 3) CHECK (recommended_potassium_mg_week IS NULL OR recommended_potassium_mg_week > 0),
    recommended_zinc_mg_week NUMERIC(12, 3) CHECK (recommended_zinc_mg_week IS NULL OR recommended_zinc_mg_week > 0),
    recommended_copper_mg_week NUMERIC(12, 3) CHECK (recommended_copper_mg_week IS NULL OR recommended_copper_mg_week > 0),
    recommended_manganese_mg_week NUMERIC(12, 3) CHECK (recommended_manganese_mg_week IS NULL OR recommended_manganese_mg_week > 0),
    recommended_selenium_mcg_week NUMERIC(12, 3) CHECK (recommended_selenium_mcg_week IS NULL OR recommended_selenium_mcg_week > 0),
    recommended_vitamin_a_mcg_week NUMERIC(12, 3) CHECK (recommended_vitamin_a_mcg_week IS NULL OR recommended_vitamin_a_mcg_week > 0),
    recommended_vitamin_c_mg_week NUMERIC(12, 3) CHECK (recommended_vitamin_c_mg_week IS NULL OR recommended_vitamin_c_mg_week > 0),
    recommended_vitamin_d_mcg_week NUMERIC(12, 3) CHECK (recommended_vitamin_d_mcg_week IS NULL OR recommended_vitamin_d_mcg_week > 0),
    recommended_vitamin_e_mg_week NUMERIC(12, 3) CHECK (recommended_vitamin_e_mg_week IS NULL OR recommended_vitamin_e_mg_week > 0),
    recommended_vitamin_k_mcg_week NUMERIC(12, 3) CHECK (recommended_vitamin_k_mcg_week IS NULL OR recommended_vitamin_k_mcg_week > 0),
    recommended_thiamin_mg_week NUMERIC(12, 3) CHECK (recommended_thiamin_mg_week IS NULL OR recommended_thiamin_mg_week > 0),
    recommended_riboflavin_mg_week NUMERIC(12, 3) CHECK (recommended_riboflavin_mg_week IS NULL OR recommended_riboflavin_mg_week > 0),
    recommended_niacin_mg_week NUMERIC(12, 3) CHECK (recommended_niacin_mg_week IS NULL OR recommended_niacin_mg_week > 0),
    recommended_pantothenic_acid_mg_week NUMERIC(12, 3) CHECK (recommended_pantothenic_acid_mg_week IS NULL OR recommended_pantothenic_acid_mg_week > 0),
    recommended_vitamin_b6_mg_week NUMERIC(12, 3) CHECK (recommended_vitamin_b6_mg_week IS NULL OR recommended_vitamin_b6_mg_week > 0),
    recommended_biotin_mcg_week NUMERIC(12, 3) CHECK (recommended_biotin_mcg_week IS NULL OR recommended_biotin_mcg_week > 0),
    recommended_folate_mcg_week NUMERIC(12, 3) CHECK (recommended_folate_mcg_week IS NULL OR recommended_folate_mcg_week > 0),
    recommended_vitamin_b12_mcg_week NUMERIC(12, 3) CHECK (recommended_vitamin_b12_mcg_week IS NULL OR recommended_vitamin_b12_mcg_week > 0),
    recommended_sodium_mg_day NUMERIC(12, 3) CHECK (recommended_sodium_mg_day IS NULL OR recommended_sodium_mg_day > 0),
    recommended_cholesterol_mg_day NUMERIC(12, 3) CHECK (recommended_cholesterol_mg_day IS NULL OR recommended_cholesterol_mg_day > 0),
    recommended_caffeine_mg_day NUMERIC(12, 3) CHECK (recommended_caffeine_mg_day IS NULL OR recommended_caffeine_mg_day > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS foods (
    id UUID PRIMARY KEY,
    canonical_name TEXT NOT NULL,
    normalized_name TEXT NOT NULL,
    display_name TEXT NOT NULL,
    preparation_state TEXT,
    nutrition_basis TEXT NOT NULL CHECK (nutrition_basis = 'per_100g'),
    provider_metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (normalized_name, preparation_state)
);

CREATE TABLE IF NOT EXISTS food_nutrition (
    food_id UUID PRIMARY KEY REFERENCES foods(id) ON DELETE CASCADE,
    calories NUMERIC(12, 3), protein NUMERIC(12, 3), carbohydrates NUMERIC(12, 3),
    fat NUMERIC(12, 3), fiber NUMERIC(12, 3), sugars NUMERIC(12, 3), saturated_fat NUMERIC(12, 3),
    sodium NUMERIC(12, 3), calcium NUMERIC(12, 3), iron NUMERIC(12, 3), magnesium NUMERIC(12, 3),
    phosphorus NUMERIC(12, 3), potassium NUMERIC(12, 3), zinc NUMERIC(12, 3), copper NUMERIC(12, 3),
    manganese NUMERIC(12, 3), selenium NUMERIC(12, 3), vitamin_a NUMERIC(12, 3), vitamin_c NUMERIC(12, 3),
    vitamin_d NUMERIC(12, 3), vitamin_e NUMERIC(12, 3), vitamin_k NUMERIC(12, 3), thiamin NUMERIC(12, 3),
    riboflavin NUMERIC(12, 3), niacin NUMERIC(12, 3), pantothenic_acid NUMERIC(12, 3), vitamin_b6 NUMERIC(12, 3),
    biotin NUMERIC(12, 3), folate NUMERIC(12, 3), vitamin_b12 NUMERIC(12, 3), caffeine NUMERIC(12, 3), cholesterol NUMERIC(12, 3),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS user_weight_history (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    weight_kg NUMERIC(12, 3) NOT NULL CHECK (weight_kg > 0),
    effective_from DATE NOT NULL,
    effective_to DATE,
    effective_range DATERANGE GENERATED ALWAYS AS (daterange(effective_from, COALESCE(effective_to, 'infinity'::date), '[)')) STORED,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (effective_to IS NULL OR effective_from < effective_to),
    EXCLUDE USING gist (user_id WITH =, effective_range WITH &&)
);
CREATE UNIQUE INDEX IF NOT EXISTS user_weight_history_current_idx ON user_weight_history (user_id) WHERE effective_to IS NULL;

CREATE TABLE IF NOT EXISTS user_diet_plan_history (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    diet_plan SMALLINT NOT NULL CHECK (diet_plan BETWEEN 1 AND 5),
    calorie_adjustment_percent NUMERIC(12, 3) NOT NULL,
    effective_from DATE NOT NULL,
    effective_to DATE,
    effective_range DATERANGE GENERATED ALWAYS AS (daterange(effective_from, COALESCE(effective_to, 'infinity'::date), '[)')) STORED,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (effective_to IS NULL OR effective_from < effective_to),
    EXCLUDE USING gist (user_id WITH =, effective_range WITH &&)
);
CREATE UNIQUE INDEX IF NOT EXISTS user_diet_plan_history_current_idx ON user_diet_plan_history (user_id) WHERE effective_to IS NULL;

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

CREATE INDEX IF NOT EXISTS foods_normalized_name_idx ON foods (normalized_name);
CREATE INDEX IF NOT EXISTS foods_normalized_name_trgm_idx ON foods USING GIN (normalized_name gin_trgm_ops);