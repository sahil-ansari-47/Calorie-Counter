from pathlib import Path
from pyspark.sql import DataFrame

def export_foods(food_df: DataFrame, output_path: Path) -> None:
    """
    Write the final application food dataset as Parquet.
    """
    output_path = Path(output_path)
    output_path.parent.mkdir(parents=True, exist_ok=True)

    (
        food_df
        .write
        .mode("overwrite")
        .parquet(str(output_path)) 
    )

def export_foods_csv(food_df: DataFrame, output_path: Path) -> None:
    """
    Write the final application food dataset as a CSV file.
    """
    output_path = Path(output_path)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    
    (
        food_df
        .write
        .mode("overwrite")
        .option("header", "true")
        .option("sep", ",")
        .csv(str(output_path)) 
    )
