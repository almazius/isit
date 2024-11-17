package automigrations

import (
	"context"
	"github.com/jmoiron/sqlx"
)

const query = `
CREATE SCHEMA users;

-- users
create table "user"
(
    id            serial
        primary key,
    name          text not null,
    login         text not null,
    password_hash text not null,
    created_at    timestamp default CURRENT_TIMESTAMP,
    blocked       boolean   default false,
    blocker_at    timestamp
);

-- roles
alter table "user"
    owner to postgres;

create unique index login_uniq_idx
    on "user" (login);

create table roles
(
    id   serial
        primary key,
    name text not null
);

alter table roles
    owner to postgres;

-- user roles
create table user_role
(
    id      serial
        primary key,
    user_id integer not null
        references "user",
    role_id integer not null
        references roles
);

alter table user_role
    owner to postgres;

create unique index user_role_idx
    on user_role (user_id, role_id);
`

func InitRepository(ctx context.Context, db *sqlx.DB) error {
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
