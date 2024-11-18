create table if not exists users."user"
(
    id            bigserial
        primary key,
    login         text                                      not null,
    password_hash text                                      not null,
    name          text                     default ''::text not null,
    surname       text                     default ''::text not null,
    created_at    timestamp with time zone default now()    not null,
    deleted_at    timestamp with time zone,
    is_baned      boolean                  default false
);

alter table users."user"
    owner to "user";

create unique index if not exists user_login_idx
    on users."user" (login);

create table if not exists cls.role
(
    id   bigserial
        primary key,
    name text not null
);

alter table cls.role
    owner to "user";

create table if not exists users.user_role
(
    id      bigserial
        primary key,
    user_id bigserial
        references users."user"
            on delete cascade,
    role_id bigserial
        references cls.role
            on delete cascade
);

alter table users.user_role
    owner to "user";

create unique index if not exists user_role_idx
    on users.user_role (user_id, role_id);

create table if not exists materials.material
(
    id             bigserial
        primary key,
    name           text                                   not null,
    price          numeric(10, 2)                         not null,
    description    text,
    created_at     timestamp with time zone default now() not null,
    deleted_at     timestamp with time zone,
    address        text                                   not null,
    reject_percent numeric(10, 2)                         not null,
    sending_date   timestamp with time zone,
    count          numeric                  default 0     not null
);

alter table materials.material
    owner to "user";

create table if not exists products.product
(
    id             bigserial
        primary key,
    name           text           not null,
    description    text,
    price          numeric(10, 2) not null,
    reject_percent numeric(10, 2) not null
);

alter table products.product
    owner to "user";

create table if not exists products.material
(
    id          bigserial
        primary key,
    product_id  bigserial
        references products.product
            on delete cascade,
    material_id bigserial
        references materials.material
            on delete cascade,
    count       numeric default 0 not null
);

alter table products.material
    owner to "user";

create table if not exists orders."order"
(
    id         bigserial
        primary key,
    product_id bigint                                 not null
        references products.product,
    count      numeric                                not null,
    created_at timestamp with time zone default now() not null,
    status     text
);

alter table orders."order"
    owner to "user";

