create table if not exists orders(
    id int unique not null,
    created_at timestamp not null default now(),
    status smallint not null
)