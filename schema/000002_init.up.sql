CREATE TABLE IF NOT EXISTS tour
(
    id                    serial       not null unique,
    tour_type             varchar(255) not null,
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