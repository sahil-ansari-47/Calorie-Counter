from pyspark.sql import DataFrame
from pyspark.sql import functions as F


FOUNDATION_FOOD_DATA_TYPE = "foundation_food"


def filter_foundation_foods(food_df: DataFrame) -> DataFrame:
    """
    Select canonical Foundation Food records from USDA food.csv.
    """
    required_columns = {"fdc_id", "data_type", "description"}

    missing = required_columns - set(food_df.columns)
    if missing:
        raise ValueError(
            f"food_df is missing required columns: {sorted(missing)}"
        )

    return food_df.filter(
        F.col("data_type") == FOUNDATION_FOOD_DATA_TYPE
    )