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
