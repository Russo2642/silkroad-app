-- Откат триггеров
DROP TRIGGER IF EXISTS update_tours_updated_at ON tours;
DROP TRIGGER IF EXISTS update_contact_form_updated_at ON contact_form;
DROP TRIGGER IF EXISTS update_help_with_tour_form_updated_at ON help_with_tour_form;
DROP TRIGGER IF EXISTS update_countries_updated_at ON countries;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Удаление индексов
DROP INDEX IF EXISTS idx_tours_type;
DROP INDEX IF EXISTS idx_tours_country;
DROP INDEX IF EXISTS idx_tours_status;
DROP INDEX IF EXISTS idx_tours_difficulty;
DROP INDEX IF EXISTS idx_tours_duration;
DROP INDEX IF EXISTS idx_tours_popular;
DROP INDEX IF EXISTS idx_tours_featured;
DROP INDEX IF EXISTS idx_tours_available_from;
DROP INDEX IF EXISTS idx_tours_available_to;
DROP INDEX IF EXISTS idx_tours_activities;
DROP INDEX IF EXISTS idx_tours_categories;
DROP INDEX IF EXISTS idx_tours_season;
DROP INDEX IF EXISTS idx_tours_keywords;
DROP INDEX IF EXISTS idx_tours_pricing;
DROP INDEX IF EXISTS idx_tours_created_at;
DROP INDEX IF EXISTS idx_tours_deleted_at;
DROP INDEX IF EXISTS idx_tours_search;

DROP INDEX IF EXISTS idx_contact_form_status;
DROP INDEX IF EXISTS idx_contact_form_created_at;
DROP INDEX IF EXISTS idx_contact_form_tour_id;
DROP INDEX IF EXISTS idx_contact_form_source;

DROP INDEX IF EXISTS idx_help_with_tour_form_status;
DROP INDEX IF EXISTS idx_help_with_tour_form_created_at;
DROP INDEX IF EXISTS idx_help_with_tour_form_when_date;
DROP INDEX IF EXISTS idx_help_with_tour_form_place;

-- Восстановление старой таблицы tours
DROP TABLE IF EXISTS tours CASCADE;

CREATE TABLE IF NOT EXISTS tours
(
    id                    serial       not null unique,
    tour_type             varchar(255) not null,
    slug                  varchar(255),
    title                 varchar(255) not null,
    tour_place            varchar(255) not null,
    season                varchar(255) not null,
    quantity              int          not null,
    duration              int          not null,
    physical_rating       int          not null,
    description_excursion varchar(255) not null,
    description_route     varchar(255) not null,
    price                 int          not null,
    currency              varchar(255) not null,
    activity              text[]       not null,
    tariff                varchar(255),
    tour_date             date         not null,
    calendar              jsonb,
    popular               boolean      DEFAULT false,
    created_at            timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

-- Удаление новых колонок из contact_form
ALTER TABLE contact_form DROP COLUMN IF EXISTS status;
ALTER TABLE contact_form DROP COLUMN IF EXISTS source;
ALTER TABLE contact_form DROP COLUMN IF EXISTS utm_source;
ALTER TABLE contact_form DROP COLUMN IF EXISTS utm_medium;
ALTER TABLE contact_form DROP COLUMN IF EXISTS utm_campaign;
ALTER TABLE contact_form DROP COLUMN IF EXISTS user_agent;
ALTER TABLE contact_form DROP COLUMN IF EXISTS ip_address;
ALTER TABLE contact_form DROP COLUMN IF EXISTS updated_at;
ALTER TABLE contact_form DROP COLUMN IF EXISTS processed_at;
ALTER TABLE contact_form DROP COLUMN IF EXISTS notes;

-- Удаление новых колонок из help_with_tour_form
ALTER TABLE help_with_tour_form DROP COLUMN IF EXISTS status;
ALTER TABLE help_with_tour_form DROP COLUMN IF EXISTS tour_type;
ALTER TABLE help_with_tour_form DROP COLUMN IF EXISTS budget;
ALTER TABLE help_with_tour_form DROP COLUMN IF EXISTS participants;
ALTER TABLE help_with_tour_form DROP COLUMN IF EXISTS preferences;
ALTER TABLE help_with_tour_form DROP COLUMN IF EXISTS comments;
ALTER TABLE help_with_tour_form DROP COLUMN IF EXISTS source;
ALTER TABLE help_with_tour_form DROP COLUMN IF EXISTS utm_source;
ALTER TABLE help_with_tour_form DROP COLUMN IF EXISTS utm_medium;
ALTER TABLE help_with_tour_form DROP COLUMN IF EXISTS utm_campaign;
ALTER TABLE help_with_tour_form DROP COLUMN IF EXISTS user_agent;
ALTER TABLE help_with_tour_form DROP COLUMN IF EXISTS ip_address;
ALTER TABLE help_with_tour_form DROP COLUMN IF EXISTS updated_at;
ALTER TABLE help_with_tour_form DROP COLUMN IF EXISTS processed_at;
ALTER TABLE help_with_tour_form DROP COLUMN IF EXISTS notes;

-- Переименование колонки обратно если нужно
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'help_with_tour_form' AND column_name = 'when_date') THEN
        ALTER TABLE help_with_tour_form RENAME COLUMN when_date TO date;
    END IF;
END $$;

-- Удаление таблицы стран
DROP TABLE IF EXISTS countries;