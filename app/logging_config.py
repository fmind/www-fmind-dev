import logging
import sys

import structlog

from app.config import settings


def setup_logging() -> None:
    """Configures structured logging for local development (Console) and production (GCP Cloud Logging JSON)."""
    is_prod = settings.is_production

    shared_processors = [
        structlog.contextvars.merge_contextvars,
        structlog.processors.add_log_level,
        structlog.processors.TimeStamper(fmt="iso" if is_prod else "%Y-%m-%d %H:%M:%S"),
        structlog.processors.StackInfoRenderer(),
        structlog.processors.format_exc_info,
    ]

    if is_prod:

        def gcp_formatter(_logger, _method_name, event_dict):
            # Map 'level' to 'severity' for GCP Cloud Logging
            if "level" in event_dict:
                event_dict["severity"] = event_dict.pop("level").upper()
            else:
                event_dict["severity"] = "INFO"

            # Map 'timestamp' to 'time'
            if "timestamp" in event_dict:
                event_dict["time"] = event_dict.pop("timestamp")

            return event_dict

        processors = [
            *shared_processors,
            gcp_formatter,
            structlog.processors.JSONRenderer(),
        ]
    else:
        processors = [
            *shared_processors,
            structlog.dev.ConsoleRenderer(),
        ]

    structlog.configure(
        processors=processors,
        logger_factory=structlog.PrintLoggerFactory(),
        cache_logger_on_first_use=True,
    )

    # Route python standard logging to structlog
    logging.basicConfig(
        format="%(message)s",
        stream=sys.stdout,
        level=logging.INFO,
    )
