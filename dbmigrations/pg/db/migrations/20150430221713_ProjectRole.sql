
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table project_role
(
	id uuid primary key default uuid_generate_v4(),
	idproject uuid not null references project(id),
	idrole uuid not null references roles(id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

drop table project_role;