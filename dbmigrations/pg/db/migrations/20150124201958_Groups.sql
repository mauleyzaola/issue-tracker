
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table groups
(
	id uuid primary key default uuid_generate_v4(),
	name text not null,
	datecreated timestamp with time zone not null,
	lastmodified timestamp with time zone,
	constraint unique_erpgroup_name unique (name)
);

create table user_group
(
	id uuid primary key default uuid_generate_v4(),
	iduser uuid not null references users(id),
	idgroup uuid not null references groups(id)
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

drop table user_group;
drop table groups;