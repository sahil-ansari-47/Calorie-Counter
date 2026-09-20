from pyspark.sql import DataFrame
from pyspark.sql import functions as F

from ingest._00_config import APPLICATION_NUTRIENTS


def select_application_nutrients(
    nutrient_df: DataFrame,
) -> DataFrame:
    """
    Select the USDA nutrient definitions used by the application.
    """

    required_columns = {"id", "name", "unit_name"}

    missing = required_columns - set(nutrient_df.columns)
    if missing:
        raise ValueError(
            f"nutrient_df is missing required columns: {sorted(missing)}"
        )

    nutrient_ids = list(APPLICATION_NUTRIENTS)

    mapping_rows = [
        (
            nutrient_id,
            config["field"],
            config["unit"],
            config["display_order"],
        )
        for nutrient_id, config in APPLICATION_NUTRIENTS.items()
    ]

    mapping_df = nutrient_df.sparkSession.createDataFrame(
        mapping_rows,
        schema=[
            "nutrient_id",
            "field",
            "expected_unit",
            "display_order",
        ],
    )

    return (
        nutrient_df
        .select(
            F.col("id").cast("long").alias("nutrient_id"),
            "name",
            "unit_name",
        )
        .filter(
            F.col("nutrient_id").isin(nutrient_ids)
        )
        .join(
            mapping_df,
            on="nutrient_id",
            how="inner",
        )
        .select(
            "nutrient_id",
            "name",
            "unit_name",
            "field",
            "display_order",
        )
        .orderBy("display_order")
    )