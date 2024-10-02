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

CREATE TABLE IF NOT EXISTS tour
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
    created_at            timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tour_editor
(
    id         serial       not null unique,
    name       varchar(255) not null,
    phone      varchar(255) not null,
    email      varchar(255) not null,
    tour_date  date         not null,
    activity   text[]       not null,
    location   text[]       not null,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

