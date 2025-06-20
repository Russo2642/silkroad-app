package tour

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/lib/pq"
)

type TourType string

const (
	OneDayTour    TourType = "Однодневный тур"
	MultiDayTour  TourType = "Многодневный тур"
	CityTour      TourType = "Сити-тур"
	ExclusiveTour TourType = "Эксклюзивный тур"
	InfoTour      TourType = "Инфо-тур"
	AuthorsTour   TourType = "Авторский тур"
)

type TourStatus string

const (
	StatusActive   TourStatus = "active"
	StatusInactive TourStatus = "inactive"
	StatusArchived TourStatus = "archived"
)

type Difficulty int

const (
	DifficultyEasy     Difficulty = 1
	DifficultyModerate Difficulty = 2
	DifficultyHard     Difficulty = 3
	DifficultyVeryHard Difficulty = 4
	DifficultyExtreme  Difficulty = 5
)

type TourRoute struct {
	Points      []RoutePoint `json:"points"`
	Description []string     `json:"description"`
	Distance    *float64     `json:"distance,omitempty"`
	Elevation   *int         `json:"elevation,omitempty"`
}

type RoutePoint struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Latitude    *float64 `json:"latitude,omitempty"`
	Longitude   *float64 `json:"longitude,omitempty"`
	Altitude    *int     `json:"altitude,omitempty"`
}

type TourIncluded struct {
	Transport     bool     `json:"transport"`
	Accommodation bool     `json:"accommodation"`
	Meals         []string `json:"meals"`
	Guide         bool     `json:"guide"`
	Equipment     []string `json:"equipment"`
	Insurance     bool     `json:"insurance"`
	Permits       bool     `json:"permits"`
}

type TourRequirements struct {
	MinAge             *int     `json:"min_age,omitempty"`
	MaxAge             *int     `json:"max_age,omitempty"`
	PhysicalFitness    string   `json:"physical_fitness"`
	Experience         string   `json:"experience,omitempty"`
	Equipment          []string `json:"equipment"`
	MedicalLimitations []string `json:"medical_limitations,omitempty"`
}

type TourPricing struct {
	BasePrice       int             `json:"base_price"`
	Currency        string          `json:"currency"`
	PriceIncludes   []string        `json:"price_includes"`
	PriceExcludes   []string        `json:"price_excludes"`
	GroupDiscounts  []GroupDiscount `json:"group_discounts,omitempty"`
	SeasonalPricing []SeasonalPrice `json:"seasonal_pricing,omitempty"`
}

type GroupDiscount struct {
	MinPeople int `json:"min_people"`
	Discount  int `json:"discount"`
}

type SeasonalPrice struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Price     int       `json:"price"`
}

type TourSchedule struct {
	DayByDay    []TourDay   `json:"day_by_day"`
	Timeline    []TimeEvent `json:"timeline,omitempty"`
	Flexibility string      `json:"flexibility"`
}

type TourDay struct {
	DayNumber     int      `json:"day_number"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Activities    []string `json:"activities"`
	Meals         []string `json:"meals"`
	Accommodation string   `json:"accommodation,omitempty"`
	Distance      *float64 `json:"distance,omitempty"`
}

type TimeEvent struct {
	Time     string `json:"time"`
	Activity string `json:"activity"`
	Location string `json:"location,omitempty"`
	Duration *int   `json:"duration,omitempty"`
}

type TourSafety struct {
	RiskLevel      string   `json:"risk_level"`
	SafetyMeasures []string `json:"safety_measures"`
	EmergencyPlan  string   `json:"emergency_plan"`
	Insurance      bool     `json:"insurance"`
	FirstAidKit    bool     `json:"first_aid_kit"`
}

type Tour struct {
	ID          int        `json:"id" db:"id"`
	Type        TourType   `json:"type" db:"type" binding:"required"`
	Status      TourStatus `json:"status" db:"status"`
	Slug        string     `json:"slug" db:"slug"`
	Title       string     `json:"title" db:"title" binding:"required"`
	Subtitle    string     `json:"subtitle" db:"subtitle"`
	Description string     `json:"description" db:"description" binding:"required"`

	Country    string `json:"country" db:"country" binding:"required"`
	Region     string `json:"region" db:"region"`
	StartPoint string `json:"start_point" db:"start_point"`
	EndPoint   string `json:"end_point" db:"end_point"`

	Duration        int        `json:"duration" db:"duration" binding:"required"`
	MinParticipants int        `json:"min_participants" db:"min_participants"`
	MaxParticipants int        `json:"max_participants" db:"max_participants" binding:"required"`
	Difficulty      Difficulty `json:"difficulty" db:"difficulty" binding:"required,gte=1,lte=5"`

	AvailableFrom time.Time      `json:"available_from" db:"available_from"`
	AvailableTo   time.Time      `json:"available_to" db:"available_to"`
	Season        pq.StringArray `json:"season" db:"season"`

	Activities pq.StringArray `json:"activities" db:"activities"`
	Categories pq.StringArray `json:"categories" db:"categories"`

	Route        TourRoute        `json:"route" db:"route"`
	Included     TourIncluded     `json:"included" db:"included"`
	Requirements TourRequirements `json:"requirements" db:"requirements"`
	Pricing      TourPricing      `json:"pricing" db:"pricing"`
	Schedule     TourSchedule     `json:"schedule" db:"schedule"`
	Safety       TourSafety       `json:"safety" db:"safety"`

	PhotosData *TourPhotosGrouped `json:"photos_data,omitempty" db:"-"`

	Keywords  pq.StringArray `json:"keywords" db:"keywords"`
	MetaTitle string         `json:"meta_title" db:"meta_title"`
	MetaDesc  string         `json:"meta_description" db:"meta_description"`

	IsPopular  bool `json:"is_popular" db:"is_popular"`
	IsFeatured bool `json:"is_featured" db:"is_featured"`
	SortOrder  int  `json:"sort_order" db:"sort_order"`

	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type TourFilter struct {
	Type        []TourType     `json:"type" form:"type"`
	Country     pq.StringArray `json:"country" form:"country"`
	Region      pq.StringArray `json:"region" form:"region"`
	Duration    *RangeFilter   `json:"duration" form:"duration"`
	Difficulty  []Difficulty   `json:"difficulty" form:"difficulty"`
	PriceMin    *int           `json:"price_min" form:"price_min"`
	PriceMax    *int           `json:"price_max" form:"price_max"`
	Quantity    *int           `json:"quantity" form:"quantity"`
	Activities  pq.StringArray `json:"activities" form:"activities"`
	Categories  pq.StringArray `json:"categories" form:"categories"`
	Season      pq.StringArray `json:"season" form:"season"`
	Popular     *bool          `json:"popular" form:"popular"`
	Featured    *bool          `json:"featured" form:"featured"`
	Available   *bool          `json:"available" form:"available"`
	SearchQuery string         `json:"search_query" form:"q"`

	Limit    int    `json:"limit" form:"limit"`
	Offset   int    `json:"offset" form:"offset"`
	SortBy   string `json:"sort_by" form:"sort_by"`
	SortDesc bool   `json:"sort_desc" form:"sort_desc"`
}

type RangeFilter struct {
	Min *int `json:"min" form:"min"`
	Max *int `json:"max" form:"max"`
}

type TourSummary struct {
	ID              int            `json:"id"`
	Type            TourType       `json:"type"`
	Slug            string         `json:"slug"`
	Title           string         `json:"title"`
	Subtitle        string         `json:"subtitle"`
	Country         string         `json:"country"`
	Region          string         `json:"region"`
	Duration        int            `json:"duration"`
	MaxParticipants int            `json:"max_participants"`
	Difficulty      Difficulty     `json:"difficulty"`
	BasePrice       int            `json:"base_price"`
	Currency        string         `json:"currency"`
	Activities      pq.StringArray `json:"activities"`
	IsPopular       bool           `json:"is_popular"`
	IsFeatured      bool           `json:"is_featured"`
	AvailableFrom   time.Time      `json:"available_from"`
	AvailableTo     time.Time      `json:"available_to"`

	PhotosData *TourPhotosGrouped `json:"photos_data,omitempty"`
}

func (tr TourRoute) Value() (driver.Value, error) {
	return json.Marshal(tr)
}

func (tr *TourRoute) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, tr)
}

func (ti TourIncluded) Value() (driver.Value, error) {
	return json.Marshal(ti)
}

func (ti *TourIncluded) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, ti)
}

func (tr TourRequirements) Value() (driver.Value, error) {
	return json.Marshal(tr)
}

func (tr *TourRequirements) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, tr)
}

func (tp TourPricing) Value() (driver.Value, error) {
	return json.Marshal(tp)
}

func (tp *TourPricing) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, tp)
}

func (ts TourSchedule) Value() (driver.Value, error) {
	return json.Marshal(ts)
}

func (ts *TourSchedule) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, ts)
}

func (ts TourSafety) Value() (driver.Value, error) {
	return json.Marshal(ts)
}

func (ts *TourSafety) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, ts)
}

func IsValidTourType(t TourType) bool {
	switch t {
	case OneDayTour, MultiDayTour, CityTour, ExclusiveTour, InfoTour, AuthorsTour:
		return true
	}
	return false
}

func IsValidDifficulty(d Difficulty) bool {
	return d >= DifficultyEasy && d <= DifficultyExtreme
}

func IsValidTourStatus(s TourStatus) bool {
	switch s {
	case StatusActive, StatusInactive, StatusArchived:
		return true
	}
	return false
}
