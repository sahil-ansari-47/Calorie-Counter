from pathlib import Path

# Project root: ingest/
PROJECT_ROOT = Path(__file__).resolve().parents[2]

# Data directories
DATA_DIR = PROJECT_ROOT / "data"
RAW_DIR = DATA_DIR / "raw"
PROCESSED_DIR = DATA_DIR / "processed"

# USDA datasets
FOUNDATION_DIR = RAW_DIR / "foundation_foods"
SURVEY_DIR = RAW_DIR / "survey"

# Pipeline outputs
OUTPUT_FILE = PROCESSED_DIR / "foods.parquet"

"""
Application-level nutrition configuration.

The USDA nutrient table contains many nutrients that are not required
by the calorie tracker. This mapping defines the subset exposed by
the application and gives each USDA nutrient a stable application
field name.
"""

APPLICATION_NUTRIENTS = {
    # Macronutrients
    1008: {
        "field": "calories",
        "unit": "KCAL",
        "display_order": 1,
    },
    1003: {
        "field": "protein",
        "unit": "G",
        "display_order": 2,
    },
    1005: {
        "field": "carbohydrates",
        "unit": "G",
        "display_order": 3,
    },
    1004: {
        "field": "fat",
        "unit": "G",
        "display_order": 4,
    },
    1079: {
        "field": "fiber",
        "unit": "G",
        "display_order": 5,
    },
    2000: {
        "field": "sugars",
        "unit": "G",
        "display_order": 6,
    },
    1258: {
        "field": "saturated_fat",
        "unit": "G",
        "display_order": 7,
    },

    # Minerals
    1093: {
        "field": "sodium",
        "unit": "MG",
        "display_order": 8,
    },
    1087: {
        "field": "calcium",
        "unit": "MG",
        "display_order": 9,
    },
    1089: {
        "field": "iron",
        "unit": "MG",
        "display_order": 10,
    },
    1090: {
        "field": "magnesium",
        "unit": "MG",
        "display_order": 11,
    },
    1091: {
        "field": "phosphorus",
        "unit": "MG",
        "display_order": 12,
    },
    1092: {
        "field": "potassium",
        "unit": "MG",
        "display_order": 13,
    },
    1095: {
        "field": "zinc",
        "unit": "MG",
        "display_order": 14,
    },
    1098: {
        "field": "copper",
        "unit": "MG",
        "display_order": 15,
    },
    1101: {
        "field": "manganese",
        "unit": "MG",
        "display_order": 16,
    },
    1103: {
        "field": "selenium",
        "unit": "UG",
        "display_order": 17,
    },

    # Vitamins
    1106: {
        "field": "vitamin_a",
        "unit": "UG",
        "display_order": 18,
    },
    1162: {
        "field": "vitamin_c",
        "unit": "MG",
        "display_order": 19,
    },
    1114: {
        "field": "vitamin_d",
        "unit": "UG",
        "display_order": 20,
    },
    1109: {
        "field": "vitamin_e",
        "unit": "MG",
        "display_order": 21,
    },
    1185: {
        "field": "vitamin_k",
        "unit": "UG",
        "display_order": 22,
    },
    1165: {
        "field": "thiamin",
        "unit": "MG",
        "display_order": 23,
    },
    1166: {
        "field": "riboflavin",
        "unit": "MG",
        "display_order": 24,
    },
    1167: {
        "field": "niacin",
        "unit": "MG",
        "display_order": 25,
    },
    1170: {
        "field": "pantothenic_acid",
        "unit": "MG",
        "display_order": 26,
    },
    1175: {
        "field": "vitamin_b6",
        "unit": "MG",
        "display_order": 27,
    },
    1176: {
        "field": "biotin",
        "unit": "UG",
        "display_order": 28,
    },
    1190: {
        "field": "folate",
        "unit": "UG",
        "display_order": 29,
    },
    1178: {
        "field": "vitamin_b12",
        "unit": "UG",
        "display_order": 30,
    },

    # Other
    1057: {
        "field": "caffeine",
        "unit": "MG",
        "display_order": 31,
    },
    1253: {
        "field": "cholesterol",
        "unit": "MG",
        "display_order": 32,
    },
}