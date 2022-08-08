create table if not exists orders(
    id int unique not null,
    created_by int not null,
    constraint fk_user foreign key(created_by) references users(id)
)