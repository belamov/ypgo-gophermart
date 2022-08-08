create table if not exists users(
    id serial unique,
    login varchar unique not null,
    password varchar not null
)