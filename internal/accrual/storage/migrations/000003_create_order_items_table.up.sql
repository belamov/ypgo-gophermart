create table if not exists accrual_order_items(
    order_id bigint not null,
    description text not null,
    price decimal not null,
    constraint fk_order_id foreign key(order_id) references accrual_orders(id)
)