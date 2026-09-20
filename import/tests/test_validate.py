import pytest

from ingest._00_config import APPLICATION_NUTRIENTS
from ingest._04_validate import validate_foods


def test_valid_food_dataset_passes(valid_food_df):
    validate_foods(valid_food_df)


def test_missing_fdc_id_fails(valid_food_df):
    invalid = valid_food_df.drop("fdc_id")

    with pytest.raises(ValueError, match="fdc_id"):
        validate_foods(invalid)


def test_missing_nutrient_column_fails(valid_food_df):
    invalid = valid_food_df.drop("protein")

    with pytest.raises(ValueError, match="protein"):
        validate_foods(invalid)


def test_duplicate_food_ids_fail(spark, valid_food_df):
    duplicate = valid_food_df.union(
        valid_food_df.limit(1)
    )

    with pytest.raises(ValueError, match="duplicate"):
        validate_foods(duplicate)


def test_negative_nutrient_values_fail(spark, valid_food_df):
    invalid = valid_food_df.withColumn(
        "protein",
        valid_food_df["protein"] * -1,
    )

    with pytest.raises(ValueError, match="negative"):
        validate_foods(invalid)


def test_required_columns_are_defined():
    expected = {
        "fdc_id",
        "name",
        *[
            config["field"]
            for config in APPLICATION_NUTRIENTS.values()
        ],
    }

    assert len(expected) == 2 + len(APPLICATION_NUTRIENTS)