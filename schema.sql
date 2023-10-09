create table if not exists cars (
    id bigserial not null primary key,
    title text not null,
    user_id bigint not null references users,
    created_at timestamptz not null default current_timestamp,
    active boolean not null default true,
    updated_at timestamptz not null default current_timestamp,
    deleted_at timestamptz
);

select * from orders where driver_id = 16 and id = 6 and status ='Заказ принят';
select * from orders where driver_id = 15 and (created_at >= '2023-09-26 00:00:00' and created_at <='2023-09-27 12:00:00');
select * from orders where driver_id = 15 and (created_at >= '2023-09-26 00:00:00 +0000 UTC' and created_at <='2023-09-27 12:00:00 +0000 UTC');
select * from orders where driver_id = 15 and created_at in('2023-09-26 00:00:00', '2023-09-27 12:00:00');

select full_name, phone_number from users u
                                    inner join orders on orders.customer_id = u.id
where orders.driver_id = 15;


select price,source,destination,duration, users.full_name Clients_Name, users.phone_number Clients_phone, orders.created_at
from orders join users on users.id = orders.customer_id
where driver_id = 15 limit 10 offset 0;

select o.created_at, o.price, o.source, o.destination, o.duration, Customers.full_name, Customers.phone_number, Drivers.full_name, Drivers.phone_number
from orders o
    left join users as Customers on o.customer_id = Customers.id
    left join users as Drivers on o.driver_id= Drivers.id;

select o.created_at, o.price, o.source, o.destination, o.duration, Customers.full_name, Customers.phone_number, Drivers.full_name, Drivers.phone_number
from orders o
         left join users as Customers on o.customer_id = Customers.id
         left join users as Drivers on o.driver_id= Drivers.id
where (o.created_at>='2023-09-20 00:00:00' and o.created_at <='2023-10-02 24:00:00')
order by o.created_at desc
limit 15 offset 0;

select role from users where id =19;

create table if not exists users (
    id bigserial not null primary key,
    full_name text not null,
    login text not null,
    password text not null,
    phone_number text not null,
    role text not null,
    created_at timestamptz not null default current_timestamp,
    active boolean not null default true,
    updated_at timestamptz not null default current_timestamp,
    deleted_at timestamptz
);

drop table users cascade;

create table if not exists billing (
    id bigserial not null primary key,
    fare text not null,
    min_price bigint not null,
    car_id bigint not null references cars,
    created_at timestamptz not null default current_timestamp,
    active boolean not null default true,
    updated_at timestamptz not null default current_timestamp,
    deleted_at timestamptz
);
select distinct fare,car_id from billing where fare = 'economy' limit ;

select * from users u
    inner join cars c on c.user_id = u.id
where c.id = 1;

insert into billing (fare, min_price, car_id)
values ('premium',15,8);
       ('standard',12,5)
values ('economy',10,1),
       ('economy',10,2),
       ('economy',10,3),
       ('standard',12,4);
       ('standard',12,4),
       ('standard',12,5),
       ('standard',12,6),
       ('premium',15,7),
       ('premium',15,8),
       ('premium',15,9);

create table if not exists roles (
    id bigserial not null primary key,
    title text not null,
    user_id bigint not null references users
);

insert into roles (title, user_id) values ('driver',1);

create table if not exists orders (
    id bigserial not null primary key,
    price float not null,
    source text not null,
    destination text not null,
    duration float not null,
    distance float not null,
    driver_id bigint references users,
    customer_id bigint not null references users,
    created_at timestamptz not null default current_timestamp,
    status text not null,
    fare text not null,
    active boolean not null default true,
    updated_at timestamptz not null default current_timestamp,
    deleted_at timestamptz
);
insert into orders (price, source, destination, duration, driver_id, customer_id, status)
values (1.3, 'adasd', 'asdas', 1.4, null,1,'ожидание');
drop table orders;

select orders.id, source, destination, duration, price, status, u.full_name customer_name, u.phone_number customer_phone
from orders join users u on u.id = orders.customer_id
where status = 'В ожидании';

alter table orders add status text ;
create table users_tokens
(
    id bigserial not null primary key,
    user_id bigint      not null references users,
    token   text        not null unique,
    expire_at  timestamptz not null default current_timestamp + interval '1 day',
    created_at timestamptz not null default current_timestamp
);
