package pg

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"silkroad/m/internal/domain/tour"
	"strings"
)

type TourPostgres struct {
	db *sqlx.DB
}

func NewTour(db *sqlx.DB) *TourPostgres {
	return &TourPostgres{db: db}
}

func (r *TourPostgres) Create(tour tour.Tour) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	calendarJSON, err := json.Marshal(tour.Calendar)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	descriptionRouteJSON, err := json.Marshal(tour.DescriptionRoute)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tour.Slug = slug.Make(tour.Title)

	var id int
	createTourQuery := fmt.Sprintf("INSERT INTO %s (tour_type, slug, title, tour_place, season, quantity, duration, "+
		"physical_rating, description_excursion, description_route, price, currency, activity, tariff, tour_date, calendar)"+
		" VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING id",
		tourTable)

	row := tx.QueryRow(createTourQuery, tour.TourType, tour.Slug, tour.Title, tour.TourPlace, tour.Season, tour.Quantity, tour.Duration,
		tour.PhysicalRating, tour.DescriptionExcursion, descriptionRouteJSON, tour.Price, tour.Currency,
		pq.Array(tour.Activity), tour.Tariff, tour.TourDate, calendarJSON)

	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TourPostgres) GetAll(tourPlace, tourDate, searchTitle string, quantity []int, priceMin, priceMax, duration, limit, offset int) ([]tour.Tour, int, int, int, int, []string, error) {
	query, args := r.buildQuery(tourPlace, tourDate, searchTitle, quantity, priceMin, priceMax, duration, limit, offset)

	totalQuery := "SELECT COUNT(*) FROM (" + query + ") AS count_query"
	var totalItems int
	err := r.db.QueryRow(totalQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, 0, 0, 0, 0, nil, err
	}

	tours, err := r.executeQuery(query, args)
	if err != nil {
		return nil, 0, 0, 0, 0, nil, err
	}

	var tourPlaces []string
	placesQuery := fmt.Sprintf("SELECT DISTINCT tour_place FROM %s ORDER BY tour_place", tourTable)
	rows, err := r.db.Query(placesQuery)
	if err != nil {
		return nil, 0, 0, 0, 0, nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			panic(err)
		}
	}()

	for rows.Next() {
		var place string
		if err := rows.Scan(&place); err != nil {
			return nil, 0, 0, 0, 0, nil, err
		}
		tourPlaces = append(tourPlaces, place)
	}

	currentPage := (offset / limit) + 1
	totalPages := (totalItems + limit - 1) / limit

	return tours, currentPage, limit, totalItems, totalPages, tourPlaces, nil
}

func (r *TourPostgres) buildQuery(tourPlace, tourDate, searchTitle string, quantity []int, priceMin, priceMax, duration, limit, offset int) (string, []interface{}) {
	var filters []string
	var args []interface{}
	argCount := 1

	query := fmt.Sprintf("SELECT id, tour_type, slug, title, tour_place, season, quantity, duration, "+
		"physical_rating, description_excursion, description_route, price, currency, activity, tariff, tour_date, calendar "+
		"FROM %s", tourTable)

	if priceMin > 0 && priceMax > 0 {
		filters = append(filters, fmt.Sprintf("price BETWEEN $%d AND $%d ", argCount, argCount+1))
		args = append(args, priceMin, priceMax)
		argCount++
	} else if priceMin > 0 {
		filters = append(filters, fmt.Sprintf("price >= $%d", argCount))
		args = append(args, priceMin)
		argCount++
	} else if priceMax > 0 {
		filters = append(filters, fmt.Sprintf("price <= $%d", argCount))
		args = append(args, priceMax)
		argCount++
	}

	if len(quantity) > 0 {
		placeholders := make([]string, len(quantity))
		for i, q := range quantity {
			placeholders[i] = fmt.Sprintf("$%d", argCount)
			args = append(args, q)
			argCount++
		}
		filters = append(filters, fmt.Sprintf("quantity IN (%s)", strings.Join(placeholders, ",")))
	}

	if tourPlace != "" {
		filters = append(filters, fmt.Sprintf("tour_place = $%d", argCount))
		args = append(args, tourPlace)
		argCount++
	}

	if duration > 0 {
		filters = append(filters, fmt.Sprintf("duration = $%d", argCount))
		args = append(args, duration)
		argCount++
	}

	if tourDate != "" {
		filters = append(filters, fmt.Sprintf("tour_date = $%d", argCount))
		args = append(args, tourDate)
		argCount++
	}

	if searchTitle != "" {
		filters = append(filters, fmt.Sprintf("title LIKE $%d", argCount))
		args = append(args, "%"+searchTitle+"%")
		argCount++
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	query += " ORDER BY created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	return query, args
}

func (r *TourPostgres) executeQuery(query string, args []interface{}) ([]tour.Tour, error) {
	var tours []tour.Tour

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			panic(err)
		}
	}()

	for rows.Next() {
		var t tour.Tour
		var calendarJSON []byte
		var descriptionRouteJSON []byte

		err := rows.Scan(&t.Id, &t.TourType, &t.Slug, &t.Title, &t.TourPlace, &t.Season, &t.Quantity, &t.Duration, &t.PhysicalRating,
			&t.DescriptionExcursion, &descriptionRouteJSON, &t.Price, &t.Currency, pq.Array(&t.Activity), &t.Tariff, &t.TourDate, &calendarJSON)
		if err != nil {
			return nil, err
		}

		if len(descriptionRouteJSON) > 0 {
			err = json.Unmarshal(descriptionRouteJSON, &t.DescriptionRoute)
			if err != nil {
				return nil, err
			}
		}

		if len(calendarJSON) > 0 {
			err = json.Unmarshal(calendarJSON, &t.Calendar)
			if err != nil {
				return nil, err
			}
		}

		tours = append(tours, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tours, nil
}

func (r *TourPostgres) GetById(tourId int) (tour.Tour, error) {
	var t tour.Tour
	var calendarJSON []byte
	var descriptionRouteJSON []byte

	query := fmt.Sprintf("SELECT id, tour_type, slug, title, tour_place, season, quantity, duration, physical_rating, "+
		"description_excursion, description_route, price, currency, activity, tariff, tour_date, calendar FROM %s WHERE id = $1",
		tourTable)

	row := r.db.QueryRow(query, tourId)
	err := row.Scan(&t.Id, &t.TourType, &t.Slug, &t.Title, &t.TourPlace, &t.Season, &t.Quantity, &t.Duration, &t.PhysicalRating,
		&t.DescriptionExcursion, &descriptionRouteJSON, &t.Price, &t.Currency, pq.Array(&t.Activity), &t.Tariff, &t.TourDate, &calendarJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return t, nil
		}
		return t, err
	}

	if len(descriptionRouteJSON) > 0 {
		err = json.Unmarshal(descriptionRouteJSON, &t.DescriptionRoute)
		if err != nil {
			return t, err
		}
	}

	if len(calendarJSON) > 0 {
		err = json.Unmarshal(calendarJSON, &t.Calendar)
		if err != nil {
			return t, err
		}
	}

	return t, nil
}

func (r *TourPostgres) GetBySlug(tourSlug string) (tour.Tour, error) {
	var t tour.Tour
	var calendarJSON []byte
	var descriptionRouteJSON []byte

	query := fmt.Sprintf("SELECT id, tour_type, slug, title, tour_place, season, quantity, duration, physical_rating, "+
		"description_excursion, description_route, price, currency, activity, tariff, tour_date, calendar FROM %s WHERE slug = $1",
		tourTable)

	row := r.db.QueryRow(query, tourSlug)
	err := row.Scan(&t.Id, &t.TourType, &t.Slug, &t.Title, &t.TourPlace, &t.Season, &t.Quantity, &t.Duration, &t.PhysicalRating,
		&t.DescriptionExcursion, &descriptionRouteJSON, &t.Price, &t.Currency, pq.Array(&t.Activity), &t.Tariff, &t.TourDate, &calendarJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return t, nil
		}
		return t, err
	}

	if len(descriptionRouteJSON) > 0 {
		err = json.Unmarshal(descriptionRouteJSON, &t.DescriptionRoute)
		if err != nil {
			return t, err
		}
	}

	if len(calendarJSON) > 0 {
		err = json.Unmarshal(calendarJSON, &t.Calendar)
		if err != nil {
			return t, err
		}
	}

	return t, nil
}

func (r *TourPostgres) GetMinMaxPrice() (int, int, error) {
	var minPrice, maxPrice int

	query := fmt.Sprintf("SELECT MIN(price), MAX(price) FROM %s", tourTable)
	row := r.db.QueryRow(query)

	err := row.Scan(&minPrice, &maxPrice)
	if err != nil {
		return 0, 0, err
	}

	return minPrice, maxPrice, nil
}
