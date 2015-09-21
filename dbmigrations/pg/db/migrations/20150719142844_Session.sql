
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied


create table sessions
(
	id uuid primary key default uuid_generate_v4(),
	iduser uuid not null references users(id),
	datecreated timestamp with time zone not null,
	expires timestamp with time zone not null,
	ipaddress text not null
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

drop table sessions;
