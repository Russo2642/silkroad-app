package pg

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"silkroad/m/internal/domain/tour"
	"strconv"
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

	tour.Slug = slug.Make(tour.Title)

	var id int
	createTourQuery := fmt.Sprintf("INSERT INTO %s (tour_type, title, tour_place, season, quantity, duration, "+
		"physical_rating, description_excursion, description_route, price, currency, activity, tariff, tour_date, calendar)"+
		" VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING id",
		tourTable)

	row := tx.QueryRow(createTourQuery, tour.TourType, tour.Title, tour.TourPlace, tour.Season, tour.Quantity, tour.Duration,
		tour.PhysicalRating, tour.DescriptionExcursion, tour.DescriptionRoute, tour.Price, tour.Currency,
		pq.Array(tour.Activity), tour.Tariff, tour.TourDate, calendarJSON)

	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TourPostgres) GetAll(priceRange, tourPlace, tourDate, searchTitle string, quantity []int, duration, limit, offset int) ([]tour.Tour, int, error) {
	query, args := r.buildQuery(priceRange, tourPlace, tourDate, searchTitle, quantity, duration, limit, offset)

	totalQuery := "SELECT COUNT(*) FROM (" + query + ") AS count_query"
	var total int
	err := r.db.QueryRow(totalQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	tours, err := r.executeQuery(query, args)
	if err != nil {
		return nil, 0, err
	}

	return tours, total, nil
}

func (r *TourPostgres) buildQuery(priceRange, tourPlace, tourDate, searchTitle string, quantity []int, duration, limit, offset int) (string, []interface{}) {
	var filters []string
	var args []interface{}
	argCount := 1

	query := fmt.Sprintf("SELECT id, tour_type, slug, title, tour_place, season, quantity, duration, "+
		"physical_rating, description_excursion, description_route, price, currency, activity, tariff, tour_date, calendar "+
		"FROM %s", tourTable)

	if priceRange != "" {
		priceParts := strings.Split(priceRange, "-")
		if len(priceParts) == 2 {
			minPrice, err1 := strconv.Atoi(priceParts[0])
			maxPrice, err2 := strconv.Atoi(priceParts[1])
			if err1 == nil && err2 == nil {
				filters = append(filters, fmt.Sprintf("price BETWEEN $%d AND $%d", argCount, argCount+1))
				args = append(args, minPrice, maxPrice)
				argCount += 2
			}
		}
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
	defer rows.Close()

	for rows.Next() {
		var t tour.Tour
		var calendarJSON []byte

		err := rows.Scan(&t.Id, &t.TourType, &t.Slug, &t.Title, &t.TourPlace, &t.Season, &t.Quantity, &t.Duration, &t.PhysicalRating,
			&t.DescriptionExcursion, &t.DescriptionRoute, &t.Price, &t.Currency, pq.Array(&t.Activity), &t.Tariff, &t.TourDate, &calendarJSON)
		if err != nil {
			return nil, err
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

	query := fmt.Sprintf("SELECT id, tour_type, slug, title, tour_place, season, quantity, duration, physical_rating, "+
		"description_excursion, description_route, price, currency, activity, tariff, tour_date, calendar FROM %s WHERE id = $1",
		tourTable)

	row := r.db.QueryRow(query, tourId)
	err := row.Scan(&t.Id, &t.TourType, &t.Slug, &t.Title, &t.TourPlace, &t.Season, &t.Quantity, &t.Duration, &t.PhysicalRating,
		&t.DescriptionExcursion, &t.DescriptionRoute, &t.Price, &t.Currency, pq.Array(&t.Activity), &t.Tariff, &t.TourDate, &calendarJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return t, nil
		}
		return t, err
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

	query := fmt.Sprintf("SELECT id, tour_type, slug, title, tour_place, season, quantity, duration, physical_rating, "+
		"description_excursion, description_route, price, currency, activity, tariff, tour_date, calendar FROM %s WHERE slug = $1",
		tourTable)

	row := r.db.QueryRow(query, tourSlug)
	err := row.Scan(&t.Id, &t.TourType, &t.Slug, &t.Title, &t.TourPlace, &t.Season, &t.Quantity, &t.Duration, &t.PhysicalRating,
		&t.DescriptionExcursion, &t.DescriptionRoute, &t.Price, &t.Currency, pq.Array(&t.Activity), &t.Tariff, &t.TourDate, &calendarJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return t, nil
		}
		return t, err
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
