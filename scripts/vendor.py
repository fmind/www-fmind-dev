"""Script to download and vendor external static resources locally."""

from pathlib import Path

import httpx

# Map of remote URLs to local paths relative to repo root
DEPENDENCIES = {
    # JS Libraries
    "https://unpkg.com/htmx.org@2.0.10/dist/htmx.min.js": "static/js/htmx.min.js",
    "https://unpkg.com/alpinejs@3.15.12/dist/cdn.min.js": "static/js/alpine.min.js",
    # Fonts
    "https://fonts.gstatic.com/s/inter/v20/UcC73FwrK3iLTeHuS_nVMrMxCp50SjIa1ZL7W0Q5nw.woff2": "static/fonts/Inter-Variable.woff2",
    "https://fonts.gstatic.com/s/outfit/v15/QGYvz_MVcBeNP4NJtEtqUYLknw.woff2": "static/fonts/Outfit-Variable.woff2",
    "https://fonts.gstatic.com/s/jetbrainsmono/v24/tDbV2o-flEEny0FZhsfKu5WU4xD7OwGtT0rU.woff2": "static/fonts/JetBrainsMono-Variable.woff2",
}


def download_file(url: str, dest_path: Path) -> None:
    """Downloads a file from the given URL to the target destination path.

    Args:
        url: The source HTTP URL to download from.
        dest_path: The local filesystem path to write the downloaded file to.
    """
    print(f"Downloading {url} -> {dest_path}...")
    headers = {
        "User-Agent": (
            "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) "
            "AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
        )
    }
    response = httpx.get(url, headers=headers, follow_redirects=True)
    response.raise_for_status()

    # Ensure parent directory exists
    dest_path.parent.mkdir(parents=True, exist_ok=True)
    dest_path.write_bytes(response.content)
    print(f"Successfully saved {dest_path} ({len(response.content)} bytes)")


def main() -> None:
    """Main function to orchestrate downloading of all vendored resources."""
    root_dir = Path(__file__).parent.parent
    for url, rel_path in DEPENDENCIES.items():
        dest = root_dir / rel_path
        try:
            download_file(url, dest)
        except Exception as e:
            print(f"Error downloading {url}: {e}")
            raise e


if __name__ == "__main__":
    main()
