
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied



create table user_meta
(
	id uuid primary key default uuid_generate_v4(),
	iduser uuid not null references users(id),
	emailnotifications boolean not null,
	realtimenotifications boolean not null,
	recieveownchanges boolean not null
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

drop table user_meta;