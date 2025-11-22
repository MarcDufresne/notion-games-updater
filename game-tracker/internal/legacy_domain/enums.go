package domain

type NotionGameProp string

const (
	PropRating           NotionGameProp = "Rating"
	PropReleaseDate      NotionGameProp = "Release Date"
	PropReleaseDateHuman NotionGameProp = "Release Date (Human)"
	PropStoreURL         NotionGameProp = "Store URL"
	PropPlatforms        NotionGameProp = "Platforms"
	PropGenres           NotionGameProp = "Genres"
	PropIGDBID           NotionGameProp = "IGDB ID"
	PropOtherResults     NotionGameProp = "_Other Results"
)

type WebsiteCategory int

const (
	WebsiteCategoryOfficial  WebsiteCategory = 1
	WebsiteCategorySteam     WebsiteCategory = 13
	WebsiteCategoryItch      WebsiteCategory = 15
	WebsiteCategoryEpicGames WebsiteCategory = 16
	WebsiteCategoryGOG       WebsiteCategory = 17
)

type ReleaseCategory int

const (
	ReleaseCategoryYYYYMMMMDD ReleaseCategory = 0
	ReleaseCategoryYYYYMMMM   ReleaseCategory = 1
	ReleaseCategoryYYYY       ReleaseCategory = 2
	ReleaseCategoryYYYYQ1     ReleaseCategory = 3
	ReleaseCategoryYYYYQ2     ReleaseCategory = 4
	ReleaseCategoryYYYYQ3     ReleaseCategory = 5
	ReleaseCategoryYYYYQ4     ReleaseCategory = 6
	ReleaseCategoryTBD        ReleaseCategory = 7
)

type ReleaseRegion int

const (
	ReleaseRegionNorthAmerica ReleaseRegion = 2
	ReleaseRegionWorldwide    ReleaseRegion = 8
)

type ReleaseDateStatusType string

const (
	ReleaseStatusOffline     ReleaseDateStatusType = "Offline"
	ReleaseStatusFullRelease ReleaseDateStatusType = "Full Release"
	ReleaseStatusBeta        ReleaseDateStatusType = "Beta"
	ReleaseStatusAlpha       ReleaseDateStatusType = "Alpha"
	ReleaseStatusEarlyAccess ReleaseDateStatusType = "Early Access"
	ReleaseStatusCancelled   ReleaseDateStatusType = "Cancelled"
)

type GameCategory int

const (
	GameCategoryMainGame            GameCategory = 0
	GameCategoryDLCAddon            GameCategory = 1
	GameCategoryExpansion           GameCategory = 2
	GameCategoryBundle              GameCategory = 3
	GameCategoryStandaloneExpansion GameCategory = 4
	GameCategoryMod                 GameCategory = 5
	GameCategoryEpisode             GameCategory = 6
	GameCategorySeason              GameCategory = 7
	GameCategoryRemake              GameCategory = 8
	GameCategoryRemaster            GameCategory = 9
	GameCategoryExpandedGame        GameCategory = 10
	GameCategoryPort                GameCategory = 11
	GameCategoryFork                GameCategory = 12
	GameCategoryPack                GameCategory = 13
	GameCategoryUpdate              GameCategory = 14
)
