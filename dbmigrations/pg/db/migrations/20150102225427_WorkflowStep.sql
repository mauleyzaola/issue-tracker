
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table workflow_step
(
	id uuid primary key default uuid_generate_v4(),
	name text not null,
	datecreated timestamp with time zone not null,
	idworkflow uuid not null references workflow(id),
	idprevstatus uuid references status(id),
	idnextstatus uuid not null references status(id),
	resolves boolean not null,
	cancels boolean not null
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

drop table workflow_step;