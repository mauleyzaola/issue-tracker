
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table file_item
(
	id uuid primary key default uuid_generate_v4(),
	iduser uuid not null references users(id),
	filedata bytea not null,
	mimetype text not null,
	bytes bigint not null,
	extension text not null,
	name text not null,
	datecreated timestamp with time zone not null
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

drop table file_item;