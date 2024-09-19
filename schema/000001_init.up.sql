CREATE TABLE IF NOT EXISTS contact_form
(
    id          serial       not null unique,
    name        varchar(255) not null,
    phone       varchar(255) not null,
    email       varchar(255) not null,
    description varchar(255) not null,
    tour_id     int,
    created_at  timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS help_with_tour_form
(
    id         serial       not null unique,
    name       varchar(255) not null,
    phone      varchar(255) not null,
    place      varchar(255) not null,
    when_date  date         not null,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
