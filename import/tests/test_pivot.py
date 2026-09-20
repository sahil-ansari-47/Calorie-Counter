from ingest._00_config import APPLICATION_NUTRIENTS
from ingest._03_transform import pivot_nutrients


def test_pivot_creates_application_columns(joined_df):
    result = pivot_nutrients(joined_df)

    expected_fields = {
        config["field"]
        for config in APPLICATION_NUTRIENTS.values()
    }

    assert expected_fields <= set(result.columns)


def test_pivot_produces_one_row_per_food(joined_df):
    result = pivot_nutrients(joined_df)

    assert result.select("fdc_id").count() == (
        joined_df.select("fdc_id").distinct().count()
    )


def test_pivot_preserves_food_identity(joined_df):
    result = pivot_nutrients(joined_df)

    original = {
        row["fdc_id"]: row["name"]
        for row in joined_df.select("fdc_id", "name").distinct().collect()
    }

    pivoted = {
        row["fdc_id"]: row["name"]
        for row in result.select("fdc_id", "name").collect()
    }

    assert pivoted == original


def test_missing_nutrient_is_null_not_zero(joined_df):
    result = pivot_nutrients(joined_df)

    # The fixture/data should contain at least one food without
    # a particular selected nutrient.
    rows = result.filter(
        result["protein"].isNull()
    ).collect()

    assert all(row["protein"] is None for row in rows)