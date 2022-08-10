create table if not exists withdraws(
    order_id int unique not null,
    user_id int not null,
    amount float not null,
    constraint fk_user foreign key(user_id) references users(id)
)