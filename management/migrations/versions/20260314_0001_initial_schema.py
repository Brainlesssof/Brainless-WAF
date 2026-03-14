"""Initial management API schema.

Revision ID: 20260314_0001
Revises:
Create Date: 2026-03-14 00:00:00.000000
"""

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = "20260314_0001"
down_revision: Union[str, Sequence[str], None] = None
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.create_table(
        "config",
        sa.Column("id", sa.Integer(), nullable=False),
        sa.Column("key", sa.String(), nullable=False),
        sa.Column("value", sa.JSON(), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=True),
        sa.PrimaryKeyConstraint("id"),
    )
    op.create_index(op.f("ix_config_id"), "config", ["id"], unique=False)
    op.create_index(op.f("ix_config_key"), "config", ["key"], unique=True)

    op.create_table(
        "events",
        sa.Column("id", sa.Integer(), nullable=False),
        sa.Column("request_id", sa.String(), nullable=True),
        sa.Column(
            "timestamp",
            sa.DateTime(timezone=True),
            server_default=sa.text("CURRENT_TIMESTAMP"),
            nullable=True,
        ),
        sa.Column("src_ip", sa.String(), nullable=True),
        sa.Column("method", sa.String(), nullable=True),
        sa.Column("path", sa.String(), nullable=True),
        sa.Column("action", sa.String(), nullable=True),
        sa.Column("severity", sa.String(), nullable=True),
        sa.Column("anomaly_score", sa.Integer(), nullable=True),
        sa.Column("triggered_rules", sa.JSON(), nullable=True),
        sa.Column("user_agent", sa.String(), nullable=True),
        sa.PrimaryKeyConstraint("id"),
    )
    op.create_index(op.f("ix_events_id"), "events", ["id"], unique=False)
    op.create_index(op.f("ix_events_request_id"), "events", ["request_id"], unique=False)
    op.create_index(op.f("ix_events_severity"), "events", ["severity"], unique=False)
    op.create_index(op.f("ix_events_src_ip"), "events", ["src_ip"], unique=False)

    op.create_table(
        "ip_lists",
        sa.Column("id", sa.Integer(), nullable=False),
        sa.Column("cidr", sa.String(), nullable=False),
        sa.Column("list_type", sa.String(), nullable=False),
        sa.Column("reason", sa.String(), nullable=True),
        sa.Column("expires_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column(
            "created_at",
            sa.DateTime(timezone=True),
            server_default=sa.text("CURRENT_TIMESTAMP"),
            nullable=True,
        ),
        sa.PrimaryKeyConstraint("id"),
    )
    op.create_index(op.f("ix_ip_lists_cidr"), "ip_lists", ["cidr"], unique=False)
    op.create_index(op.f("ix_ip_lists_id"), "ip_lists", ["id"], unique=False)

    op.create_table(
        "rules",
        sa.Column("id", sa.Integer(), nullable=False),
        sa.Column("rule_id", sa.Integer(), nullable=False),
        sa.Column("rule_text", sa.Text(), nullable=False),
        sa.Column("description", sa.String(), nullable=True),
        sa.Column("enabled", sa.Boolean(), nullable=True),
        sa.Column("source", sa.String(), nullable=True),
        sa.Column("phase", sa.Integer(), nullable=True),
        sa.Column("severity", sa.String(), nullable=True),
        sa.Column(
            "created_at",
            sa.DateTime(timezone=True),
            server_default=sa.text("CURRENT_TIMESTAMP"),
            nullable=True,
        ),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=True),
        sa.PrimaryKeyConstraint("id"),
    )
    op.create_index(op.f("ix_rules_id"), "rules", ["id"], unique=False)
    op.create_index(op.f("ix_rules_rule_id"), "rules", ["rule_id"], unique=True)

    op.create_table(
        "users",
        sa.Column("id", sa.Integer(), nullable=False),
        sa.Column("username", sa.String(), nullable=False),
        sa.Column("email", sa.String(), nullable=False),
        sa.Column("hashed_password", sa.String(), nullable=False),
        sa.Column("is_active", sa.Boolean(), nullable=True),
        sa.Column("role", sa.String(), nullable=True),
        sa.Column(
            "created_at",
            sa.DateTime(timezone=True),
            server_default=sa.text("CURRENT_TIMESTAMP"),
            nullable=True,
        ),
        sa.PrimaryKeyConstraint("id"),
    )
    op.create_index(op.f("ix_users_email"), "users", ["email"], unique=True)
    op.create_index(op.f("ix_users_id"), "users", ["id"], unique=False)
    op.create_index(op.f("ix_users_username"), "users", ["username"], unique=True)


def downgrade() -> None:
    op.drop_index(op.f("ix_users_username"), table_name="users")
    op.drop_index(op.f("ix_users_id"), table_name="users")
    op.drop_index(op.f("ix_users_email"), table_name="users")
    op.drop_table("users")

    op.drop_index(op.f("ix_rules_rule_id"), table_name="rules")
    op.drop_index(op.f("ix_rules_id"), table_name="rules")
    op.drop_table("rules")

    op.drop_index(op.f("ix_ip_lists_id"), table_name="ip_lists")
    op.drop_index(op.f("ix_ip_lists_cidr"), table_name="ip_lists")
    op.drop_table("ip_lists")

    op.drop_index(op.f("ix_events_src_ip"), table_name="events")
    op.drop_index(op.f("ix_events_severity"), table_name="events")
    op.drop_index(op.f("ix_events_request_id"), table_name="events")
    op.drop_index(op.f("ix_events_id"), table_name="events")
    op.drop_table("events")

    op.drop_index(op.f("ix_config_key"), table_name="config")
    op.drop_index(op.f("ix_config_id"), table_name="config")
    op.drop_table("config")