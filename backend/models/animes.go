package models

// Translation представляет перевод
type Translation struct {
	ID    int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

type SeasonEpisode struct {
	Season   int   `json:"season"`
	Episodes []int `json:"episodes" gorm:"type:int[]"`
}

type BlockedSeasons struct {
	All interface{} `json:"all,omitempty" gorm:"type:jsonb"`
}

type MaterialData struct {
	Title            *string   `json:"title,omitempty"`
	AnimeTitle       *string   `json:"anime_title,omitempty"`
	TitleEn          *string   `json:"title_en,omitempty"`
	OtherTitles      *[]string `json:"other_titles,omitempty" gorm:"type:text[]"`
	OtherTitlesEn    *[]string `json:"other_titles_en,omitempty" gorm:"type:text[]"`
	OtherTitlesJp    *[]string `json:"other_titles_jp,omitempty" gorm:"type:text[]"`
	AnimeLicenseName *string   `json:"anime_license_name,omitempty"`
	AnimeLicensedBy  *[]string `json:"anime_licensed_by,omitempty" gorm:"type:text[]"`
	AnimeKind        *string   `json:"anime_kind,omitempty"`
	AllStatus        *string   `json:"all_status,omitempty"`
	AnimeStatus      *string   `json:"anime_status,omitempty"`
	Year             *int      `json:"year,omitempty"`
	Description      *string   `json:"description,omitempty"`
	PosterURL        *string   `json:"poster_url,omitempty"`
	Screenshots      *[]string `json:"screenshots,omitempty" gorm:"type:text[]"`
	Duration         *int      `json:"duration,omitempty"`
	Countries        *[]string `json:"countries,omitempty" gorm:"type:text[]"`
	AllGenres        *[]string `json:"all_genres,omitempty" gorm:"type:text[]"`
	Genres           *[]string `json:"genres,omitempty" gorm:"type:text[]"`
	AnimeGenres      *[]string `json:"anime_genres,omitempty" gorm:"type:text[]"`
	AnimeStudios     *[]string `json:"anime_studios,omitempty" gorm:"type:text[]"`
	KinopoiskRating  *float64  `json:"kinopoisk_rating,omitempty"`
	KinopoiskVotes   *int      `json:"kinopoisk_votes,omitempty"`
	ImdbRating       *float64  `json:"imdb_rating,omitempty"`
	ImdbVotes        *int      `json:"imdb_votes,omitempty"`
	ShikimoriRating  *float64  `json:"shikimori_rating,omitempty"`
	ShikimoriVotes   *int      `json:"shikimori_votes,omitempty"`
	AiredAt          *string   `json:"aired_at,omitempty"`
	ReleasedAt       *string   `json:"released_at,omitempty"`
	NextEpisodeAt    *string   `json:"next_episode_at,omitempty"`
	RatingMpaa       *string   `json:"rating_mpaa,omitempty"`
	MinimalAge       *int      `json:"minimal_age,omitempty"`
	EpisodesTotal    *int      `json:"episodes_total,omitempty"`
	EpisodesAired    *int      `json:"episodes_aired,omitempty"`
	Actors           *[]string `json:"actors,omitempty" gorm:"type:text[]"`
	Directors        *[]string `json:"directors,omitempty" gorm:"type:text[]"`
	Producers        *[]string `json:"producers,omitempty" gorm:"type:text[]"`
	Writers          *[]string `json:"writers,omitempty" gorm:"type:text[]"`
	Composers        *[]string `json:"composers,omitempty" gorm:"type:text[]"`
	Editors          *[]string `json:"editors,omitempty" gorm:"type:text[]"`
	Designers        *[]string `json:"designers,omitempty" gorm:"type:text[]"`
	Operators        *[]string `json:"operators,omitempty" gorm:"type:text[]"`
}

type Anime struct {
	ID               string           `json:"id" gorm:"primaryKey"`
	Title            string           `json:"title"`
	TitleOrig        string           `json:"title_orig"`
	OtherTitle       *string          `json:"other_title,omitempty"`
	Link             string           `json:"link"`
	Year             int              `json:"year"`
	KinopoiskID      *int             `json:"kinopoisk_id,omitempty"`
	ImdbID           *string          `json:"imdb_id,omitempty"`
	MdlID            *int             `json:"mdl_id,omitempty"`
	WorldArtLink     *string          `json:"worldart_link,omitempty"`
	ShikimoriID      *string          `json:"shikimori_id,omitempty"`
	Type             string           `json:"type"`
	Quality          string           `json:"quality"`
	Camrip           bool             `json:"camrip"`
	LGBT             bool             `json:"lgbt"`
	TranslationID    int              `json:"translation_id"` // Foreign key for Translation
	Translation      Translation      `json:"translation" gorm:"foreignKey:TranslationID"`
	CreatedAt        string           `json:"created_at"`
	UpdatedAt        string           `json:"updated_at"`
	BlockedCountries []string         `json:"blocked_countries" gorm:"type:text[]"`
	Seasons          *[]SeasonEpisode `json:"seasons,omitempty" gorm:"type:jsonb"`
	LastSeason       *int             `json:"last_season,omitempty"`
	LastEpisode      *int             `json:"last_episode,omitempty"`
	EpisodesCount    *int             `json:"episodes_count,omitempty"`
	BlockedSeasons   BlockedSeasons   `json:"blocked_seasons" gorm:"type:jsonb"`
	Screenshots      []string         `json:"screenshots" gorm:"type:text[]"`
	MaterialData     *MaterialData    `json:"material_data,omitempty" gorm:"type:jsonb"`
}

func (t Anime) Validate() error {
	var err error
	return err
}
