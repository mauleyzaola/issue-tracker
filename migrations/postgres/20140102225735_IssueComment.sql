
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create table issue_comment
(
	id uuid primary key default uuid_generate_v4(),
	datecreated timestamp with time zone not null,
	lastmodified timestamp with time zone,
	idissue uuid not null references issue(id),
	iduser uuid not null references users(id),
	body text not null
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

drop table issue_comment;