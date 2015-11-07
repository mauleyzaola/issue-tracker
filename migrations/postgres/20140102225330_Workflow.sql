
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create table workflow
(
	id uuid primary key default uuid_generate_v4(),
	name text not null,
	datecreated timestamp with time zone not null,
	lastmodified timestamp with time zone,
	constraint uix_workflow_name unique (name)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

drop table workflow;