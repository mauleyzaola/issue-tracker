
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table project_role_member
(
	id uuid primary key default uuid_generate_v4(),
	idprojectrole uuid not null references project_role(id),
	idgroup uuid null references groups(id),
	iduser uuid null references users(id)
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

drop table project_role_member;