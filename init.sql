create table "user"
(
    id   serial primary key,
    name varchar
);

create table "account"
(
    id      serial,
    user_id int primary key references "user" (id),
    balance int
        constraint positive_balance CHECK (balance >= 0)
);

create table "operation"
(
    id         serial primary key,
    member     int references "user" (id),
    created_at timestamp default now(),
    meta       json
);

-- create users
insert into "user" (name)
values ('kolya'),
       ('petya'),
       ('vasya'),
       ('tanya');