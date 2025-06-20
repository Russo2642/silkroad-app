CREATE TABLE IF NOT EXISTS contact_form
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    tour_id INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS help_with_tour_form
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL,
    when_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_contact_form_created_at ON contact_form(created_at);
CREATE INDEX IF NOT EXISTS idx_contact_form_tour_id ON contact_form(tour_id);
CREATE INDEX IF NOT EXISTS idx_contact_form_tour_id ON contact_form(phone);

CREATE INDEX IF NOT EXISTS idx_help_with_tour_form_country ON help_with_tour_form(country);
CREATE INDEX IF NOT EXISTS idx_help_with_tour_form_when_date ON help_with_tour_form(when_date);
CREATE INDEX IF NOT EXISTS idx_help_with_tour_form_created_at ON help_with_tour_form(created_at);
