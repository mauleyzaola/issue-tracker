
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table project
(
	 id uuid primary key default uuid_generate_v4(),
	 pkey text not null,
	 name text,
	 idprojectlead uuid not null references users(id),	
	 datecreated timestamp with time zone NOT NULL,
	 issuecount integer not null default(0),
	 notresolvedcount integer not null default(0),
	 lastmodified timestamp with time zone,
	 begins timestamp with time zone,
	 ends timestamp with time zone,
	 next bigint not null default(1),
	 constraint uix_project_pkey unique(pkey)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

drop table project;