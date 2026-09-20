from ingest._05_export import export_foods


def test_export_creates_readable_parquet(
    spark,
    valid_food_df,
    tmp_path,
):
    output_path = tmp_path / "foods.parquet"

    export_foods(
        valid_food_df,
        output_path,
    )

    assert output_path.exists()

    reloaded = spark.read.parquet(str(output_path))

    assert reloaded.count() == valid_food_df.count()
    assert reloaded.columns == valid_food_df.columns