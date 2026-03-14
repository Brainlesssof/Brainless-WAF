from pathlib import Path

from alembic import command
from alembic.config import Config
from sqlalchemy import create_engine, inspect, text


ROOT_DIR = Path(__file__).resolve().parents[1]


def build_alembic_config(database_url: str) -> Config:
    config = Config(str(ROOT_DIR / "alembic.ini"))
    config.set_main_option("script_location", str(ROOT_DIR / "migrations"))
    config.set_main_option("sqlalchemy.url", database_url)
    return config


def test_alembic_upgrade_creates_management_schema(tmp_path):
    database_path = tmp_path / "migrations.db"
    database_url = f"sqlite:///{database_path.as_posix()}"
    config = build_alembic_config(database_url)

    command.upgrade(config, "head")

    engine = create_engine(database_url)
    inspector = inspect(engine)

    assert set(inspector.get_table_names()) == {
        "alembic_version",
        "config",
        "events",
        "ip_lists",
        "rules",
        "users",
    }

    with engine.connect() as connection:
        version = connection.execute(text("SELECT version_num FROM alembic_version")).scalar_one()

    assert version == "20260314_0001"


def test_alembic_downgrade_removes_management_schema(tmp_path):
    database_path = tmp_path / "migrations.db"
    database_url = f"sqlite:///{database_path.as_posix()}"
    config = build_alembic_config(database_url)

    command.upgrade(config, "head")
    command.downgrade(config, "base")

    engine = create_engine(database_url)
    inspector = inspect(engine)

    remaining_tables = set(inspector.get_table_names())

    assert remaining_tables.isdisjoint({"config", "events", "ip_lists", "rules", "users"})