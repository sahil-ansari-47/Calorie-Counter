from pathlib import Path

import pytest
from pyspark.sql import SparkSession

from ingest._01_filters import filter_foundation_foods
from ingest._02_nutrients import select_application_nutrients
from ingest._03_transform import (
    join_foods_with_nutrients,
    pivot_nutrients,
)


PROJECT_ROOT = Path(__file__).resolve().parents[1]

DATA_DIR = (
    PROJECT_ROOT
    / "data"
    / "raw"
    / "foundation_foods"
)


def read_csv(spark, filename):
    return (
        spark.read
        .option("header", True)
        .option("quote", '"')
        .option("escape", '"')
        .option("multiLine", True)
        .csv(str(DATA_DIR / filename))
    )


@pytest.fixture(scope="session")
def spark():
    spark = (
        SparkSession.builder
        .master("local[*]")
        .appName("CalorieTrackerIngestionTests")
        .getOrCreate()
    )

    yield spark
    spark.stop()


@pytest.fixture(scope="session")
def food_df(spark):
    return read_csv(spark, "food.csv")


@pytest.fixture(scope="session")
def foundation_df(spark):
    return read_csv(spark, "foundation_food.csv")


@pytest.fixture(scope="session")
def foundation_food_df(food_df):
    return filter_foundation_foods(food_df)


@pytest.fixture(scope="session")
def food_nutrient_df(spark):
    return read_csv(spark, "food_nutrient.csv")


@pytest.fixture(scope="session")
def nutrient_df(spark):
    return read_csv(spark, "nutrient.csv")


@pytest.fixture(scope="session")
def joined_df(
    foundation_food_df,
    food_nutrient_df,
    nutrient_df,
):
    selected_nutrients = select_application_nutrients(
        nutrient_df
    )

    return join_foods_with_nutrients(
        foundation_food_df,
        food_nutrient_df,
        selected_nutrients,
    )


@pytest.fixture(scope="session")
def valid_food_df(joined_df):
    return pivot_nutrients(joined_df)


@pytest.fixture(scope="session")
def food_category_df(spark):
    return read_csv(spark, "food_category.csv")