from app.models import Playlist, Project

OPEN_SOURCE = [
    Project(
        title="mlops-python-package",
        href="https://github.com/fmind/mlops-python-package",
        description="Kickstart your MLOps initiative with a flexible, robust, and productive Python package.",
        type="github",
    ),
    Project(
        title="cookiecutter-mlops-package",
        href="https://github.com/fmind/cookiecutter-mlops-package",
        description="Start building and deploying Python packages and Docker images for MLOps.",
        type="github",
    ),
    Project(
        title="MLOps Coding Course",
        href="https://mlops-coding-course.fmind.dev/",
        repo="https://github.com/MLOps-Courses/mlops-coding-course",
        description="Learn to create, develop, and maintain a state-of-the-art MLOps code base.",
        type="github",
    ),
]

YOUTUBE_SERIES = [
    Playlist(
        title="Bleeding Agent",
        url="https://www.youtube.com/playlist?list=PLPCnNL6Y2PbTckW80gDLnznFMDEz18HBS",
        description="Technical deep dives into the Black Box of AI Agents and emerging autonomous systems.",
        cta="View Podcast",
    ),
    Playlist(
        title="AI Agents in a Nut$HELL",
        url="https://www.youtube.com/playlist?list=PLPCnNL6Y2PbT1aKOx2fMFBpicTRzeMS6f",
        description="Brief, high-signal deep dives into the core architecture and inner workings of AI Agents.",
        cta="View Playlist",
    ),
    Playlist(
        title="MLOps Coding Course",
        url="https://www.youtube.com/playlist?list=PLPCnNL6Y2PbQplCczUFhtQpCznEXqDZnh",
        description="Bridge the gap between robust software engineering and cutting-edge data science.",
        cta="View Course",
    ),
]
