CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS foods (
    id UUID PRIMARY KEY,
    canonical_name TEXT NOT NULL,
    normalized_name TEXT NOT NULL,
    display_name TEXT NOT NULL,
    preparation_state TEXT,
    quantity NUMERIC(12, 3) NOT NULL CHECK (quantity > 0),
    quantity_unit TEXT NOT NULL CHECK (quantity_unit IN ('g', 'ml', 'oz', 'serving')),
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

CREATE INDEX IF NOT EXISTS foods_normalized_name_idx ON foods (normalized_name);
CREATE INDEX IF NOT EXISTS foods_normalized_name_trgm_idx ON foods USING GIN (normalized_name gin_trgm_ops);