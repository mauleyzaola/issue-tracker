
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table issue
(
	id uuid primary key default uuid_generate_v4(),
	pkey text not null,
	name text not null,
	description text not null,
	datecreated timestamp with time zone not null,
	lastmodified timestamp with time zone,
	idstatus uuid not null references status(id),
	idworkflow uuid not null references workflow(id),
	idpriority uuid not null references priority(id),
	idproject uuid null references project(id),
	duedate timestamp with time zone not null,
	resolveddate timestamp with time zone,
	cancelleddate timestamp with time zone,
	idassignee uuid references users(id),
	idreporter uuid not null references users(id),
	idparent uuid references issue(id),
	constraint unique_issue_pkey unique (pkey)
);

create sequence seq_issue;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

drop sequence seq_issue;
drop table issue;