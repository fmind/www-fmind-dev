"""Configuration module for www.fmind.dev settings and environment variables."""

from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    """Application settings loaded from environment variables and .env file."""

    env: str = "prod"
    session_secret: str = "dev-secret-key-change-me-in-production"  # noqa: S105

    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        extra="ignore",
    )

    @property
    def is_production(self) -> bool:
        """Returns True if the application is running in production mode."""
        return self.env == "prod"


settings = Settings()
