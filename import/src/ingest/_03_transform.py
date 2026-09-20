from pyspark.sql import DataFrame
from pyspark.sql import functions as F


def join_foods_with_nutrients(
    food_df: DataFrame,
    food_nutrient_df: DataFrame,
    application_nutrients_df: DataFrame,
) -> DataFrame:

    required_food_columns = {
        "fdc_id",
        "description",
    }

    required_measurement_columns = {
        "fdc_id",
        "nutrient_id",
        "amount",
    }

    required_nutrient_columns = {
        "nutrient_id",
        "field",
    }

    missing_food = required_food_columns - set(food_df.columns)
    if missing_food:
        raise ValueError(
            f"food_df is missing columns: {sorted(missing_food)}"
        )

    missing_measurements = (
        required_measurement_columns
        - set(food_nutrient_df.columns)
    )

    if missing_measurements:
        raise ValueError(
            "food_nutrient_df is missing columns: "
            f"{sorted(missing_measurements)}"
        )

    missing_nutrients = (
        required_nutrient_columns
        - set(application_nutrients_df.columns)
    )

    if missing_nutrients:
        raise ValueError(
            "application_nutrients_df is missing columns: "
            f"{sorted(missing_nutrients)}"
        )

    foods = food_df.select(
        F.col("fdc_id").cast("long").alias("fdc_id"),
        F.col("description").alias("name"),
    )

    measurements = food_nutrient_df.select(
        F.col("fdc_id").cast("long").alias("fdc_id"),
        F.col("nutrient_id").cast("long").alias("nutrient_id"),
        F.col("amount").cast("double").alias("amount"),
    )

    nutrients = application_nutrients_df.select(
        F.col("nutrient_id").cast("long").alias("nutrient_id"),
        "field",
    )

    return (
        foods
        .join(
            measurements,
            on="fdc_id",
            how="inner",
        )
        .join(
            nutrients,
            on="nutrient_id",
            how="inner",
        )
    )

from ingest._00_config import APPLICATION_NUTRIENTS


def pivot_nutrients(joined_df: DataFrame) -> DataFrame:
    """
    Convert long-form nutrient measurements into one row per food.
    """

    required_columns = {
        "fdc_id",
        "name",
        "field",
        "amount",
    }

    missing = required_columns - set(joined_df.columns)

    if missing:
        raise ValueError(
            f"joined_df is missing columns: {sorted(missing)}"
        )

    nutrient_fields = [
        config["field"]
        for config in sorted(
            APPLICATION_NUTRIENTS.values(),
            key=lambda item: item["display_order"],
        )
    ]

    pivoted = (
        joined_df
        .groupBy("fdc_id", "name")
        .pivot("field", nutrient_fields)
        .agg(F.first("amount"))
    )

    pivoted = pivoted.withColumn(
        "carbohydrates",
        F.greatest(
            F.col("carbohydrates"),
            F.lit(0.0),
        ),
    )

    return pivoted.select(
        "fdc_id",
        "name",
        *nutrient_fields,
    )