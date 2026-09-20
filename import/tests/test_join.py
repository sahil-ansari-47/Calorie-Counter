from ingest._02_nutrients import select_application_nutrients
from ingest._03_transform import join_foods_with_nutrients


def test_food_nutrient_join_uses_fdc_id(
    foundation_food_df,
    food_nutrient_df,
    nutrient_df,
):
    selected_nutrients = select_application_nutrients(
        nutrient_df
    )

    result = join_foods_with_nutrients(
        foundation_food_df,
        food_nutrient_df,
        selected_nutrients,
    )

    foundation_ids = {
        int(row["fdc_id"])
        for row in foundation_food_df
        .select("fdc_id")
        .collect()
    }

    result_ids = {
        row["fdc_id"]
        for row in result
        .select("fdc_id")
        .distinct()
        .collect()
    }

    assert result_ids <= foundation_ids


def test_join_contains_application_nutrient_field(
    foundation_food_df,
    food_nutrient_df,
    nutrient_df,
):
    selected_nutrients = select_application_nutrients(
        nutrient_df
    )

    result = join_foods_with_nutrients(
        foundation_food_df,
        food_nutrient_df,
        selected_nutrients,
    )

    assert "field" in result.columns
    assert "amount" in result.columns


def test_join_does_not_include_unselected_nutrients(
    foundation_food_df,
    food_nutrient_df,
    nutrient_df,
):
    selected_nutrients = select_application_nutrients(
        nutrient_df
    )

    result = join_foods_with_nutrients(
        foundation_food_df,
        food_nutrient_df,
        selected_nutrients,
    )

    selected_ids = {
        row["nutrient_id"]
        for row in selected_nutrients
        .select("nutrient_id")
        .collect()
    }

    result_ids = {
        row["nutrient_id"]
        for row in result
        .select("nutrient_id")
        .distinct()
        .collect()
    }

    assert result_ids <= selected_ids