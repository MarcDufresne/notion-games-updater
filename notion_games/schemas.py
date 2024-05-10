from datetime import datetime, date
from enum import IntEnum, StrEnum
from typing import Annotated

from pydantic import BaseModel, HttpUrl, Field


class NotionGameProp(StrEnum):
    RATING = "Rating"
    RELEASE_DATE = "Release Date"
    RELEASE_DATE_HUMAN = "Release Date (Human)"
    STORE_URL = "Store URL"
    PLATFORMS = "Platforms"
    GENRES = "Genres"
    IGDB_ID = "IGDB ID"
    OTHER_RESULTS = "_Other Results"


class Cover(BaseModel):
    id: int
    image_id: str

    def cover_big_2x_url(self) -> str:
        return f"https://images.igdb.com/igdb/image/upload/t_cover_big_2x/{self.image_id}.jpg"


class Platform(BaseModel):
    id: int
    name: str
    abbreviation: str


class Genre(BaseModel):
    id: int
    name: str


class WebsiteCategory(IntEnum):
    official = 1
    steam = 13
    itch = 15
    epicgames = 16
    gog = 17


class Website(BaseModel):
    category: WebsiteCategory | int
    url: HttpUrl


class ReleaseCategory(IntEnum):
    YYYYMMMMDD = 0
    YYYYMMMM = 1
    YYYY = 2
    YYYYQ1 = 3
    YYYYQ2 = 4
    YYYYQ3 = 5
    YYYYQ4 = 6
    TBD = 7


class ReleaseRegion(IntEnum):
    north_america = 2
    worldwide = 8


class ReleaseDateStatusType(StrEnum):
    offline = "Offline"
    full_release = "Full Release"
    beta = "Beta"
    alpha = "Alpha"
    early_access = "Early Access"
    cancelled = "Cancelled"


class ReleaseDateStatus(BaseModel):
    id: int
    name: ReleaseDateStatusType | str


class ReleaseDate(BaseModel):
    category: ReleaseCategory | int
    date_: Annotated[date | None, Field(alias="date")] = None
    human: str | None = None
    y: int | None = None
    m: int | None = None
    region: ReleaseRegion | int
    status: ReleaseDateStatus | None = None
    platform: Platform | None = None


class GameCategory(IntEnum):
    main_game = 0
    dlc_addon = 1
    expansion = 2
    bundle = 3
    standalone_expansion = 4
    mod = 5
    episode = 6
    season = 7
    remake = 8
    remaster = 9
    expanded_game = 10
    port = 11
    fork = 12
    pack = 13
    update = 14


class Game(BaseModel):
    id: int
    aggregated_rating: float | None = None
    category: GameCategory | int
    cover: Cover | None = None
    first_release_date: int | None = None
    genres: list[Genre] = []
    name: str
    platforms: list[Platform] = []
    release_dates: list[ReleaseDate] = []
    updated_at: datetime | None = None
    url: HttpUrl | None = None
    websites: list[Website] = []
