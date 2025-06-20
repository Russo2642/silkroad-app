-- Создание триггерной функции для обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Создание таблицы стран
CREATE TABLE IF NOT EXISTS countries (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    code VARCHAR(2) NOT NULL UNIQUE,
    flag_url VARCHAR(500),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Полная переработка таблицы tours (перемещено сюда)
DROP TABLE IF EXISTS tours CASCADE;

CREATE TABLE tours (
    id SERIAL PRIMARY KEY,
    
    -- Базовая информация
    type VARCHAR(100) NOT NULL CHECK (type IN ('Однодневный тур', 'Многодневный тур', 'Сити-тур', 'Эксклюзивный тур', 'Инфо-тур', 'Авторский тур')),
    status VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'archived')),
    slug VARCHAR(255) UNIQUE,
    title VARCHAR(500) NOT NULL,
    subtitle VARCHAR(1000),
    description TEXT NOT NULL,
    
    -- Локация и география
    country VARCHAR(255) NOT NULL,
    region VARCHAR(255),
    start_point VARCHAR(500),
    end_point VARCHAR(500),
    
    -- Характеристики тура
    duration INTEGER NOT NULL,
    min_participants INTEGER DEFAULT 1,
    max_participants INTEGER NOT NULL,
    difficulty INTEGER NOT NULL CHECK (difficulty >= 1 AND difficulty <= 5),
    
    -- Расписание и доступность
    available_from DATE,
    available_to DATE,
    season TEXT[], -- массив месяцев
    
    -- Активности и категории
    activities TEXT[] NOT NULL,
    categories TEXT[],
    
    -- Структурированная информация (JSONB поля)
    route JSONB,
    included JSONB,
    requirements JSONB,
    pricing JSONB NOT NULL,
    schedule JSONB,
    safety JSONB,
    
    -- Метадата и SEO
    keywords TEXT[],
    meta_title VARCHAR(1000),
    meta_description TEXT,
    
    -- Управление
    is_popular BOOLEAN DEFAULT false,
    is_featured BOOLEAN DEFAULT false,
    sort_order INTEGER DEFAULT 0,
    
    -- Временные метки
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Создание индексов для tours
CREATE INDEX idx_tours_type ON tours(type);
CREATE INDEX idx_tours_country ON tours(country);
CREATE INDEX idx_tours_status ON tours(status);
CREATE INDEX idx_tours_difficulty ON tours(difficulty);
CREATE INDEX idx_tours_duration ON tours(duration);
CREATE INDEX idx_tours_popular ON tours(is_popular);
CREATE INDEX idx_tours_featured ON tours(is_featured);
CREATE INDEX idx_tours_available_from ON tours(available_from);
CREATE INDEX idx_tours_available_to ON tours(available_to);
CREATE INDEX idx_tours_activities ON tours USING GIN(activities);
CREATE INDEX idx_tours_categories ON tours USING GIN(categories);
CREATE INDEX idx_tours_season ON tours USING GIN(season);
CREATE INDEX idx_tours_keywords ON tours USING GIN(keywords);
CREATE INDEX idx_tours_pricing ON tours USING GIN(pricing);
CREATE INDEX idx_tours_created_at ON tours(created_at);
CREATE INDEX idx_tours_deleted_at ON tours(deleted_at);
CREATE INDEX idx_tours_search ON tours USING GIN(to_tsvector('russian', title || ' ' || description));

-- Создание таблицы фотографий туров (теперь таблица tours уже существует)
CREATE TABLE IF NOT EXISTS tour_photos (
    id SERIAL PRIMARY KEY,
    tour_id INTEGER NOT NULL,
    photo_url VARCHAR(1000) NOT NULL,
    photo_type VARCHAR(50) NOT NULL CHECK (photo_type IN ('preview', 'gallery', 'route', 'booking')),
    title VARCHAR(500),
    description TEXT,
    alt_text VARCHAR(500),
    display_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tour_id) REFERENCES tours(id) ON DELETE CASCADE
);

-- Создание индексов для tour_photos
CREATE INDEX idx_tour_photos_tour_id ON tour_photos(tour_id);
CREATE INDEX idx_tour_photos_type ON tour_photos(photo_type);
CREATE INDEX idx_tour_photos_active ON tour_photos(is_active);
CREATE INDEX idx_tour_photos_display_order ON tour_photos(display_order);


-- Создание триггеров для обновления updated_at

CREATE TRIGGER update_tours_updated_at
    BEFORE UPDATE ON tours
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tour_photos_updated_at
    BEFORE UPDATE ON tour_photos
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_countries_updated_at
    BEFORE UPDATE ON countries
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Вставка примерных данных стран
INSERT INTO countries (name, code, flag_url, is_active) VALUES 
('Казахстан', 'KZ', NULL, true),
('Кыргызстан', 'KG', NULL, true),
('Узбекистан', 'UZ', NULL, true),
('Таджикистан', 'TJ', NULL, true),
('Туркменистан', 'TM', NULL, true),
('Россия', 'RU', NULL, true),
('Китай', 'CN', NULL, true),
('Монголия', 'MN', NULL, true)
ON CONFLICT (code) DO NOTHING; 