create table if not exists withdraws(
    order_id bigint unique not null,
    user_id int not null,
    amount float not null,
    created_at timestamp not null default now(),
    constraint fk_user foreign key(user_id) references users(id)
)