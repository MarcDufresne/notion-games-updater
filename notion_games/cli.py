import time

import typer
from loguru import logger

from notion_games.use_cases import update_games_db

cli = typer.Typer()


@cli.command()
def update_games_db_cmd(run_forever: bool = False, interval_min: int = 15):
    try:
        if run_forever:
            while True:
                update_games_db.run()
                logger.info(f"Sleeping for {interval_min} minutes")
                time.sleep(interval_min * 60)
        else:
            update_games_db.run()
    except Exception as e:
        logger.exception(e)
        raise typer.Exit(code=1) from e
