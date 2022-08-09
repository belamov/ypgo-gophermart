create table if not exists orders(
    id int unique not null,
    created_by int not null,
    uploaded_at timestamp not null default now(),
    status smallint not null,
    accrual decimal,
    constraint fk_user foreign key(created_by) references users(id)
)