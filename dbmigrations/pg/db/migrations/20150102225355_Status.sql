
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table status
(
	id uuid primary key default uuid_generate_v4(),
	idworkflow uuid not null references workflow(id),
	name text not null,
	description text,
	datecreated timestamp with time zone not null,
	lastmodified timestamp with time zone
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

drop table status;