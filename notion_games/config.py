import os

from dotenv import load_dotenv

load_dotenv()


IGDB_CLIENT_ID = os.environ["IGDB_CLIENT_ID"]
IGDB_CLIENT_SECRET = os.environ["IGDB_CLIENT_SECRET"]
NOTION_TOKEN = os.environ["NOTION_TOKEN"]
NOTION_DATABASE_ID = os.environ["NOTION_DATABASE_ID"]
