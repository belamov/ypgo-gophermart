create table if not exists users(
    id serial,
    login varchar unique not null,
    password varchar not null
)