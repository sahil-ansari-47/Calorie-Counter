from pyspark.sql import DataFrame
from pyspark.sql import functions as F

from ingest._00_config import APPLICATION_NUTRIENTS


def validate_foods(food_df: DataFrame) -> None:
    """
    Validate the final application-level food dataset.

    Raises ValueError when the dataset violates the expected schema
    or contains invalid nutrition values.
    """

    required_columns = {
        "fdc_id",
        "name",
        *[
            config["field"]
            for config in APPLICATION_NUTRIENTS.values()
        ],
    }

    missing = required_columns - set(food_df.columns)

    if missing:
        raise ValueError(
            f"Final food dataset is missing columns: {sorted(missing)}"
        )

    if food_df.filter(F.col("fdc_id").isNull()).limit(1).count():
        raise ValueError("Final food dataset contains null fdc_id")

    if food_df.filter(F.col("name").isNull()).limit(1).count():
        raise ValueError("Final food dataset contains null name")

    duplicate_ids = (
        food_df
        .groupBy("fdc_id")
        .count()
        .filter(F.col("count") > 1)
        .limit(1)
        .count()
    )

    if duplicate_ids:
        raise ValueError(
            "Final food dataset contains duplicate fdc_id values"
        )

    nutrient_columns = [
        config["field"]
        for config in APPLICATION_NUTRIENTS.values()
    ]

    for column in nutrient_columns:
        negative_count = (
            food_df
            .filter(
                F.col(column).isNotNull()
                & (F.col(column) < 0)
            )
            .limit(1)
            .count()
        )

        if negative_count:
            raise ValueError(
                f"Final food dataset contains negative values "
                f"for nutrient '{column}'"
            )