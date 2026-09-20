from ingest._00_config import APPLICATION_NUTRIENTS
from ingest._02_nutrients import select_application_nutrients


def test_all_application_nutrients_exist(nutrient_df):
    selected = select_application_nutrients(nutrient_df)

    configured_ids = set(APPLICATION_NUTRIENTS)

    actual_ids = {
        row["nutrient_id"]
        for row in selected.select("nutrient_id").collect()
    }

    assert actual_ids == configured_ids


def test_application_nutrient_units_match_config(nutrient_df):
    selected = select_application_nutrients(nutrient_df)

    rows = selected.select(
        "nutrient_id",
        "unit_name",
    ).collect()

    actual_units = {
        row["nutrient_id"]: row["unit_name"]
        for row in rows
    }

    for nutrient_id, config in APPLICATION_NUTRIENTS.items():
        assert actual_units[nutrient_id] == config["unit"]


def test_application_field_names_are_unique():
    fields = [
        config["field"]
        for config in APPLICATION_NUTRIENTS.values()
    ]

    assert len(fields) == len(set(fields))


def test_application_display_orders_are_unique():
    orders = [
        config["display_order"]
        for config in APPLICATION_NUTRIENTS.values()
    ]

    assert len(orders) == len(set(orders))


def test_application_metadata_matches_config(nutrient_df):
    selected = select_application_nutrients(nutrient_df)

    rows = selected.collect()

    assert len(rows) == len(APPLICATION_NUTRIENTS)

    for row in rows:
        config = APPLICATION_NUTRIENTS[row["nutrient_id"]]

        assert row["field"] == config["field"]
        assert row["unit_name"] == config["unit"]
        assert row["display_order"] == config["display_order"]