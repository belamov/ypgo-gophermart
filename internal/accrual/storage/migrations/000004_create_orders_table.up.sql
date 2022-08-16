create table if not exists accrual_orders(
    id bigint unique not null,
    created_at timestamp not null default now(),
    status smallint not null,
    accrual decimal
)