"""
C2 Server Entry Point.
This is the main entry point for the C2 server application.
"""
from api import app

# This allows Gunicorn to import the app
# The app is defined in api.py