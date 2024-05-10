from typing import Any

import pendulum
from loguru import logger
from notion_client import Client

from notion_games import config
from notion_games.external.igdb_client import IGDB
from notion_games.schemas import Game, NotionGameProp, ReleaseRegion, WebsiteCategory

notion = Client(auth=config.NOTION_TOKEN)
igdb_client = IGDB(config.IGDB_CLIENT_ID, config.IGDB_CLIENT_SECRET)


def _get_game(page: dict[str, Any]) -> tuple[Game | None, dict[str, Any]]:
    page_title = page["properties"]["Game"]["title"][0]["plain_text"]

    properties_updates = {}

    try:
        igdb_id_prop = page["properties"].get(NotionGameProp.IGDB_ID, {})
        igdb_id = igdb_id_prop["rich_text"][0]["plain_text"]
    except (KeyError, IndexError):
        igdb_id = None

    if igdb_id is not None:
        if not igdb_id.startswith(":"):
            properties_updates[NotionGameProp.OTHER_RESULTS] = {
                "type": "rich_text",
                "rich_text": [{"text": {"content": ""}}],
            }
        else:
            igdb_id = igdb_id[1:]

        game_res = igdb_client.request(
            "games",
            (
                f"fields name,url,aggregated_rating,category,first_release_date,"
                f"platforms.*,cover.*,genres.*,websites.*,release_dates.*,"
                f"parent_game.id,parent_game.name,url,updated_at; "
                f"where id = {igdb_id};"
            ),
        )

        if len(game_res) == 0:
            logger.warning(f"Game {igdb_id} not found in IGDB, skipping page '{page_title}'")
            return None, properties_updates

        game = Game(**game_res[0])
    else:
        search_res = igdb_client.request(
            "search",
            (
                f"fields game.name,game.url,game.aggregated_rating,game.category,game.first_release_date,"
                f"game.platforms.*,game.cover.*,game.genres.*,game.websites.*,game.release_dates.*,"
                f"game.url,game.updated_at; "
                f'search "{page_title}"; '
                f"where game != null & game.category != (13) & game.version_parent = null;"
            ),
        )

        if len(search_res) == 0:
            logger.warning(f"No results found on IGDB for page '{page_title}', skipping page")
            return None, properties_updates
        elif len(search_res) > 1:
            logger.warning(f"Multiple results found for {page_title}; defaulting to the first one")
            items = []
            for res in search_res:
                g: dict[str, Any] = res["game"]
                items.append(
                    {
                        "text": {
                            "content": f"{g['name']}",
                            "link": {"url": g["url"]},
                        }
                    }
                )
                items.append({"text": {"content": f" ({g['id']})\n"}})

            properties_updates[NotionGameProp.OTHER_RESULTS] = {"type": "rich_text", "rich_text": items}

        game = Game(**search_res[0]["game"])

    return game, properties_updates


def _update_page(page_id: str) -> None:
    logger.info(f"Processing page: {page_id}")
    page = notion.pages.retrieve(page_id)

    page_title = page["properties"]["Game"]["title"][0]["plain_text"]
    logger.info(f"Game: {page_title}")

    game, properties = _get_game(page)

    # TODO: Set attributes on the game in Notion (release date, platforms, cover image)
    properties.update(
        {
            NotionGameProp.IGDB_ID: {
                "type": "rich_text",
                "rich_text": [{"text": {"content": f":{game.id}"}}],
            },
        }
    )
    if game.aggregated_rating:
        logger.debug(f"Rating: {game.aggregated_rating}")
        properties[NotionGameProp.RATING] = {
            "number": round(game.aggregated_rating / 100, 2),
        }

    if game.platforms:
        platforms = []
        for platform in game.platforms:
            logger.debug(f"Platform: {platform.abbreviation}")
            platforms.append({"name": platform.abbreviation})

        properties[NotionGameProp.PLATFORMS] = {
            "type": "multi_select",
            "multi_select": platforms,
        }

    if game.genres:
        genres = []
        for genre in game.genres:
            logger.debug(f"Genre: {genre.name}")
            genres.append({"name": genre.name})

        properties[NotionGameProp.GENRES] = {
            "type": "multi_select",
            "multi_select": genres,
        }

    if game.websites:
        priority_order = [WebsiteCategory.steam, WebsiteCategory.itch, WebsiteCategory.gog, WebsiteCategory.epicgames]
        category_to_url = {website.category: str(website.url) for website in game.websites}
        url = None

        for category in priority_order:
            if category in category_to_url:
                url = category_to_url[category]
                break

        if url:
            logger.debug(f"Store URL: {url}")
            properties[NotionGameProp.STORE_URL] = {
                "type": "url",
                "url": url,
            }

    if game.first_release_date:
        from_timestamp = pendulum.from_timestamp(game.first_release_date)
        logger.debug(f"Release date: {from_timestamp.format('MMMM D, YYYY')}")
        properties[NotionGameProp.RELEASE_DATE] = {
            "type": "date",
            "date": {
                "start": from_timestamp.format("YYYY-MM-DD"),
            },
        }
        properties[NotionGameProp.RELEASE_DATE_HUMAN] = {
            "type": "rich_text",
            "rich_text": [{"text": {"content": from_timestamp.format("MMMM D, YYYY")}}],
        }
    elif game.release_dates:
        filtered_release_dates = [
            release
            for release in game.release_dates
            if release.region in [ReleaseRegion.worldwide, ReleaseRegion.north_america]
        ]
        dates_mapping: list[tuple[pendulum.DateTime, str]] = []
        for release in filtered_release_dates:
            if 0 <= release.category <= 6:
                if release.category == 0:
                    release_date = pendulum.instance(release.date)
                elif release.category == 1:
                    release_date = pendulum.date(release.y, release.m, 1).end_of("month")
                elif release.category == 2:
                    release_date = pendulum.date(release.y, 1, 1).end_of("year")
                elif release.category == 3:
                    release_date = pendulum.date(release.y, 3, 1).end_of("month")
                elif release.category == 4:
                    release_date = pendulum.date(release.y, 6, 1).end_of("month")
                elif release.category == 5:
                    release_date = pendulum.date(release.y, 9, 1).end_of("month")
                else:
                    release_date = pendulum.date(release.y, 12, 1).end_of("month")

                dates_mapping.append((release_date, release.human))
                break

        if dates_mapping:
            date = sorted(dates_mapping, key=lambda x: x[0])[0]
            logger.debug(f"Release date: {date[0].format('MMMM D, YYYY')}")
            properties[NotionGameProp.RELEASE_DATE] = {
                "type": "date",
                "date": {
                    "start": date[0].format("YYYY-MM-DD"),
                },
            }
            properties[NotionGameProp.RELEASE_DATE_HUMAN] = {
                "type": "rich_text",
                "rich_text": [{"text": {"content": date[1]}}],
            }
        else:
            logger.debug("No release date found")
            properties[NotionGameProp.RELEASE_DATE] = {
                "type": "date",
                "date": None,
            }
            properties[NotionGameProp.RELEASE_DATE_HUMAN] = {
                "type": "rich_text",
                "rich_text": [{"text": {"content": "TBD"}}],
            }

    properties["Game"] = {
        "type": "title",
        "title": [{"text": {"content": game.name}}],
    }

    update_kwargs = {"properties": properties}

    if game.cover:
        update_kwargs["cover"] = {
            "type": "external",
            "external": {"url": game.cover.cover_big_2x_url()},
        }

    notion.pages.update(page_id, **update_kwargs)
    logger.success(f"Updated page '{page_title}'")


def run():
    games_db_res = notion.databases.query(database_id=config.NOTION_DATABASE_ID)

    if games_db_res["object"] != "list":
        raise Exception("Not a list database")

    for res in games_db_res["results"]:
        if res["object"] == "page":
            try:
                _update_page(res["id"])
            except Exception as e:
                logger.exception(e)
                logger.error(f"Failed to update page: {res['id']}")

    logger.success("Database updated")
