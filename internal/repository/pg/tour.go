package pg

import (
	"database/sql"
	"fmt"
	"silkroad/m/internal/domain/tour"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
)

type TourRepository struct {
	db *sqlx.DB
}

func NewTour(db *sqlx.DB) *TourRepository {
	return &TourRepository{db: db}
}

func (r *TourRepository) Create(t tour.Tour) (int, error) {
	if t.Slug == "" {
		t.Slug = r.generateUniqueSlug(t.Title)
	}

	var id int
	query := `
		INSERT INTO tours (
			type, status, slug, title, subtitle, description,
			country, region, start_point, end_point,
			duration, min_participants, max_participants, difficulty,
			available_from, available_to, season,
			activities, categories,
			route, included, requirements, pricing, schedule, safety,
			keywords, meta_title, meta_description,
			is_popular, is_featured, sort_order
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10,
			$11, $12, $13, $14,
			$15, $16, $17,
			$18, $19,
			$20, $21, $22, $23, $24, $25,
			$26, $27, $28,
			$29, $30, $31
		) RETURNING id`

	err := r.db.QueryRow(query,
		t.Type, t.Status, t.Slug, t.Title, t.Subtitle, t.Description,
		t.Country, t.Region, t.StartPoint, t.EndPoint,
		t.Duration, t.MinParticipants, t.MaxParticipants, t.Difficulty,
		t.AvailableFrom, t.AvailableTo, t.Season,
		t.Activities, t.Categories,
		t.Route, t.Included, t.Requirements, t.Pricing, t.Schedule, t.Safety,
		t.Keywords, t.MetaTitle, t.MetaDesc,
		t.IsPopular, t.IsFeatured, t.SortOrder,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create tour: %w", err)
	}

	return id, nil
}

func (r *TourRepository) GetByID(id int) (tour.Tour, error) {
	var t tour.Tour
	query := `
		SELECT id, type, status, slug, title, subtitle, description,
			   country, region, start_point, end_point,
			   duration, min_participants, max_participants, difficulty,
			   available_from, available_to, season,
			   activities, categories,
			   route, included, requirements, pricing, schedule, safety,
			   keywords, meta_title, meta_description,
			   is_popular, is_featured, sort_order,
			   created_at, updated_at, deleted_at
		FROM tours 
		WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.Get(&t, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return t, fmt.Errorf("tour not found")
		}
		return t, fmt.Errorf("failed to get tour: %w", err)
	}

	photosData, err := r.loadTourPhotos(t.ID)
	if err != nil {
		fmt.Printf("Warning: failed to load photos for tour %d: %v\n", t.ID, err)
	} else {
		t.PhotosData = photosData
	}

	return t, nil
}

func (r *TourRepository) GetBySlug(slug string) (tour.Tour, error) {
	var t tour.Tour
	query := `
		SELECT id, type, status, slug, title, subtitle, description,
			   country, region, start_point, end_point,
			   duration, min_participants, max_participants, difficulty,
			   available_from, available_to, season,
			   activities, categories,
			   route, included, requirements, pricing, schedule, safety,
			   keywords, meta_title, meta_description,
			   is_popular, is_featured, sort_order,
			   created_at, updated_at, deleted_at
		FROM tours 
		WHERE slug = $1 AND deleted_at IS NULL`

	err := r.db.Get(&t, query, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return t, fmt.Errorf("tour not found")
		}
		return t, fmt.Errorf("failed to get tour: %w", err)
	}

	photosData, err := r.loadTourPhotos(t.ID)
	if err != nil {
		fmt.Printf("Warning: failed to load photos for tour %d: %v\n", t.ID, err)
	} else {
		t.PhotosData = photosData
	}

	return t, nil
}

func (r *TourRepository) GetAll(filter tour.TourFilter) ([]tour.Tour, int, error) {
	query, countQuery, args, countArgs := r.buildFilterQuery(filter, false)

	var total int
	err := r.db.Get(&total, countQuery, countArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tours: %w", err)
	}

	var tours []tour.Tour
	err = r.db.Select(&tours, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get tours: %w", err)
	}

	if len(tours) > 0 {
		tourIDs := make([]int, len(tours))
		for i, t := range tours {
			tourIDs[i] = t.ID
		}

		photosMap, err := r.loadTourPhotosForMultiple(tourIDs)
		if err != nil {
			fmt.Printf("Warning: failed to load photos for tours: %v\n", err)
		} else {
			for i := range tours {
				if photos, exists := photosMap[tours[i].ID]; exists {
					tours[i].PhotosData = photos
				}
			}
		}
	}

	return tours, total, nil
}

func (r *TourRepository) GetSummaries(filter tour.TourFilter) ([]tour.TourSummary, int, error) {
	query, countQuery, args, countArgs := r.buildFilterQuery(filter, true)

	var total int
	err := r.db.Get(&total, countQuery, countArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tours: %w", err)
	}

	var summaries []tour.TourSummary
	err = r.db.Select(&summaries, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get tour summaries: %w", err)
	}

	if len(summaries) > 0 {
		tourIDs := make([]int, len(summaries))
		for i, s := range summaries {
			tourIDs[i] = s.ID
		}

		photosMap, err := r.loadTourPhotosForMultiple(tourIDs)
		if err != nil {
			fmt.Printf("Warning: failed to load photos for tour summaries: %v\n", err)
		} else {
			for i := range summaries {
				if photos, exists := photosMap[summaries[i].ID]; exists {
					summaries[i].PhotosData = photos
				}
			}
		}
	}

	return summaries, total, nil
}

func (r *TourRepository) Update(t tour.Tour) error {
	query := `
		UPDATE tours SET 
			type = $1, status = $2, slug = $3, title = $4, subtitle = $5, description = $6,
			country = $7, region = $8, start_point = $9, end_point = $10,
			duration = $11, min_participants = $12, max_participants = $13, difficulty = $14,
			available_from = $15, available_to = $16, season = $17,
			activities = $18, categories = $19,
			route = $20, included = $21, requirements = $22, pricing = $23, schedule = $24, safety = $25,
			keywords = $26, meta_title = $27, meta_description = $28,
			is_popular = $29, is_featured = $30, sort_order = $31,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $32 AND deleted_at IS NULL`

	result, err := r.db.Exec(query,
		t.Type, t.Status, t.Slug, t.Title, t.Subtitle, t.Description,
		t.Country, t.Region, t.StartPoint, t.EndPoint,
		t.Duration, t.MinParticipants, t.MaxParticipants, t.Difficulty,
		t.AvailableFrom, t.AvailableTo, t.Season,
		t.Activities, t.Categories,
		t.Route, t.Included, t.Requirements, t.Pricing, t.Schedule, t.Safety,
		t.Keywords, t.MetaTitle, t.MetaDesc,
		t.IsPopular, t.IsFeatured, t.SortOrder,
		t.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update tour: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tour not found")
	}

	return nil
}

func (r *TourRepository) Delete(id int) error {
	query := `
		UPDATE tours 
		SET deleted_at = CURRENT_TIMESTAMP 
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tour: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tour not found")
	}

	return nil
}

func (r *TourRepository) GetMinMaxPrice() (int, int, error) {
	query := `
		SELECT 
			COALESCE(MIN(CAST(pricing->>'base_price' AS INTEGER)), 0) as min_price,
			COALESCE(MAX(CAST(pricing->>'base_price' AS INTEGER)), 0) as max_price
		FROM tours 
		WHERE status = 'active' AND deleted_at IS NULL`

	var minPrice, maxPrice int
	err := r.db.QueryRow(query).Scan(&minPrice, &maxPrice)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get min/max prices: %w", err)
	}

	return minPrice, maxPrice, nil
}

func (r *TourRepository) GetFilterValues() (map[string][]string, error) {
	filters := make(map[string][]string)

	countries := []string{}
	err := r.db.Select(&countries, "SELECT DISTINCT country FROM tours WHERE status = 'active' AND deleted_at IS NULL ORDER BY country")
	if err != nil {
		return nil, fmt.Errorf("failed to get countries: %w", err)
	}
	filters["countries"] = countries

	regions := []string{}
	err = r.db.Select(&regions, "SELECT DISTINCT region FROM tours WHERE status = 'active' AND deleted_at IS NULL AND region IS NOT NULL ORDER BY region")
	if err != nil {
		return nil, fmt.Errorf("failed to get regions: %w", err)
	}
	filters["regions"] = regions

	activitiesQuery := `
		SELECT DISTINCT unnest(activities) as activity 
		FROM tours 
		WHERE status = 'active' AND deleted_at IS NULL AND activities IS NOT NULL
		ORDER BY activity`
	activities := []string{}
	err = r.db.Select(&activities, activitiesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get activities: %w", err)
	}
	filters["activities"] = activities

	categoriesQuery := `
		SELECT DISTINCT unnest(categories) as category 
		FROM tours 
		WHERE status = 'active' AND deleted_at IS NULL AND categories IS NOT NULL
		ORDER BY category`
	categories := []string{}
	err = r.db.Select(&categories, categoriesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	filters["categories"] = categories

	return filters, nil
}

func (r *TourRepository) buildFilterQuery(filter tour.TourFilter, summaryOnly bool) (string, string, []interface{}, []interface{}) {
	var selectFields string
	if summaryOnly {
		selectFields = `
			id, type, slug, title, subtitle, country, region, duration, max_participants, 
			difficulty, (pricing->>'base_price')::int as base_price, (pricing->>'currency') as currency,
			activities, is_popular, is_featured,
			available_from, available_to`
	} else {
		selectFields = `
			id, type, status, slug, title, subtitle, description,
			country, region, start_point, end_point,
			duration, min_participants, max_participants, difficulty,
			available_from, available_to, season,
			activities, categories,
			route, included, requirements, pricing, schedule, safety,
			keywords, meta_title, meta_description,
			is_popular, is_featured, sort_order,
			created_at, updated_at, deleted_at`
	}

	baseQuery := fmt.Sprintf("SELECT %s FROM tours", selectFields)
	countQuery := "SELECT COUNT(*) FROM tours"

	var conditions []string
	var args []interface{}
	argIndex := 1

	conditions = append(conditions, "deleted_at IS NULL")

	if len(filter.Type) > 0 {
		placeholders := make([]string, len(filter.Type))
		for i, t := range filter.Type {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, t)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("type IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(filter.Country) > 0 {
		placeholders := make([]string, len(filter.Country))
		for i, c := range filter.Country {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, c)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("country IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(filter.Region) > 0 {
		placeholders := make([]string, len(filter.Region))
		for i, r := range filter.Region {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, r)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("region IN (%s)", strings.Join(placeholders, ",")))
	}

	if filter.Duration != nil {
		if filter.Duration.Min != nil {
			conditions = append(conditions, fmt.Sprintf("duration >= $%d", argIndex))
			args = append(args, *filter.Duration.Min)
			argIndex++
		}
		if filter.Duration.Max != nil {
			conditions = append(conditions, fmt.Sprintf("duration <= $%d", argIndex))
			args = append(args, *filter.Duration.Max)
			argIndex++
		}
	}

	if len(filter.Difficulty) > 0 {
		placeholders := make([]string, len(filter.Difficulty))
		for i, d := range filter.Difficulty {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, d)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("difficulty IN (%s)", strings.Join(placeholders, ",")))
	}

	if filter.PriceMin != nil {
		conditions = append(conditions, fmt.Sprintf("CAST(pricing->>'base_price' AS INTEGER) >= $%d", argIndex))
		args = append(args, *filter.PriceMin)
		argIndex++
	}

	if filter.PriceMax != nil {
		conditions = append(conditions, fmt.Sprintf("CAST(pricing->>'base_price' AS INTEGER) <= $%d", argIndex))
		args = append(args, *filter.PriceMax)
		argIndex++
	}

	if filter.Quantity != nil {
		conditions = append(conditions, fmt.Sprintf("min_participants <= $%d AND max_participants >= $%d", argIndex, argIndex+1))
		args = append(args, *filter.Quantity, *filter.Quantity)
		argIndex += 2
	}

	if len(filter.Activities) > 0 {
		conditions = append(conditions, fmt.Sprintf("activities && $%d", argIndex))
		args = append(args, filter.Activities)
		argIndex++
	}

	if len(filter.Categories) > 0 {
		conditions = append(conditions, fmt.Sprintf("categories && $%d", argIndex))
		args = append(args, filter.Categories)
		argIndex++
	}

	if len(filter.Season) > 0 {
		conditions = append(conditions, fmt.Sprintf("season && $%d", argIndex))
		args = append(args, filter.Season)
		argIndex++
	}

	if filter.Popular != nil {
		conditions = append(conditions, fmt.Sprintf("is_popular = $%d", argIndex))
		args = append(args, *filter.Popular)
		argIndex++
	}

	if filter.Featured != nil {
		conditions = append(conditions, fmt.Sprintf("is_featured = $%d", argIndex))
		args = append(args, *filter.Featured)
		argIndex++
	}

	if filter.Available != nil && *filter.Available {
		now := time.Now()
		conditions = append(conditions, fmt.Sprintf("available_from <= $%d AND available_to >= $%d", argIndex, argIndex+1))
		args = append(args, now, now)
		argIndex += 2
	}

	if filter.SearchQuery != "" {
		conditions = append(conditions, fmt.Sprintf(`(
			to_tsvector('russian', title || ' ' || COALESCE(subtitle, '') || ' ' || description || ' ' || 
			array_to_string(activities, ' ') || ' ' || array_to_string(categories, ' ') || ' ' || 
			country || ' ' || COALESCE(region, '')) @@ plainto_tsquery('russian', $%d)
			OR LOWER(title) LIKE LOWER($%d) 
			OR LOWER(subtitle) LIKE LOWER($%d)
			OR LOWER(country) LIKE LOWER($%d)
		)`, argIndex, argIndex+1, argIndex+2, argIndex+3))

		searchPattern := "%" + filter.SearchQuery + "%"
		args = append(args, filter.SearchQuery, searchPattern, searchPattern, searchPattern)
		argIndex += 4
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	baseQuery += whereClause
	countQuery += whereClause

	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)

	orderBy := " ORDER BY "
	switch filter.SortBy {
	case "price":
		orderBy += "CAST(pricing->>'base_price' AS INTEGER)"
	case "duration":
		orderBy += "duration"
	case "popularity":
		orderBy += "is_popular DESC, is_featured DESC"
	case "created_at":
		orderBy += "created_at"
	default:
		orderBy += "sort_order, created_at"
	}

	if filter.SortDesc {
		orderBy += " DESC"
	} else {
		orderBy += " ASC"
	}

	baseQuery += orderBy

	if filter.Limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++

		baseQuery += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
		argIndex++
	}

	return baseQuery, countQuery, args, countArgs
}

func (r *TourRepository) generateUniqueSlug(title string) string {
	baseSlug := slug.Make(title)
	uniqueSlug := baseSlug

	suffix := 1
	for {
		var count int
		err := r.db.QueryRow("SELECT COUNT(*) FROM tours WHERE slug = $1", uniqueSlug).Scan(&count)
		if err != nil {
			break
		}

		if count == 0 {
			break
		}

		uniqueSlug = fmt.Sprintf("%s-%d", baseSlug, suffix)
		suffix++
	}

	return uniqueSlug
}

func (r *TourRepository) loadTourPhotos(tourID int) (*tour.TourPhotosGrouped, error) {
	query := `
		SELECT id, tour_id, photo_url, photo_type, title, description, alt_text, 
		       display_order, is_active, created_at, updated_at
		FROM tour_photos 
		WHERE tour_id = $1 AND is_active = true
		ORDER BY photo_type, display_order, created_at`

	var photos []tour.TourPhoto
	err := r.db.Select(&photos, query, tourID)
	if err != nil {
		return nil, fmt.Errorf("failed to load tour photos: %w", err)
	}

	grouped := &tour.TourPhotosGrouped{
		Preview: make([]tour.TourPhoto, 0),
		Gallery: make([]tour.TourPhoto, 0),
		Route:   make([]tour.TourPhoto, 0),
		Booking: make([]tour.TourPhoto, 0),
	}

	for _, photo := range photos {
		switch photo.PhotoType {
		case tour.PhotoTypePreview:
			grouped.Preview = append(grouped.Preview, photo)
		case tour.PhotoTypeGallery:
			grouped.Gallery = append(grouped.Gallery, photo)
		case tour.PhotoTypeRoute:
			grouped.Route = append(grouped.Route, photo)
		case tour.PhotoTypeBooking:
			grouped.Booking = append(grouped.Booking, photo)
		}
	}

	return grouped, nil
}

func (r *TourRepository) loadTourPhotosForMultiple(tourIDs []int) (map[int]*tour.TourPhotosGrouped, error) {
	if len(tourIDs) == 0 {
		return make(map[int]*tour.TourPhotosGrouped), nil
	}

	placeholders := make([]string, len(tourIDs))
	args := make([]interface{}, len(tourIDs))
	for i, id := range tourIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT id, tour_id, photo_url, photo_type, title, description, alt_text, 
		       display_order, is_active, created_at, updated_at
		FROM tour_photos 
		WHERE tour_id IN (%s) AND is_active = true
		ORDER BY tour_id, photo_type, display_order, created_at`,
		strings.Join(placeholders, ","))

	var photos []tour.TourPhoto
	err := r.db.Select(&photos, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to load tour photos: %w", err)
	}

	result := make(map[int]*tour.TourPhotosGrouped)

	for _, photo := range photos {
		if _, exists := result[photo.TourID]; !exists {
			result[photo.TourID] = &tour.TourPhotosGrouped{
				Preview: make([]tour.TourPhoto, 0),
				Gallery: make([]tour.TourPhoto, 0),
				Route:   make([]tour.TourPhoto, 0),
				Booking: make([]tour.TourPhoto, 0),
			}
		}

		switch photo.PhotoType {
		case tour.PhotoTypePreview:
			result[photo.TourID].Preview = append(result[photo.TourID].Preview, photo)
		case tour.PhotoTypeGallery:
			result[photo.TourID].Gallery = append(result[photo.TourID].Gallery, photo)
		case tour.PhotoTypeRoute:
			result[photo.TourID].Route = append(result[photo.TourID].Route, photo)
		case tour.PhotoTypeBooking:
			result[photo.TourID].Booking = append(result[photo.TourID].Booking, photo)
		}
	}

	return result, nil
}
