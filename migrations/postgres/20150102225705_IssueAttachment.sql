
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create table issue_attachment
(
	id uuid primary key default uuid_generate_v4(),
	idissue uuid not null references issue(id),
	idfileitem uuid not null references file_item(id),
	iduser uuid not null references users(id),
	datecreated timestamp with time zone NOT NULL
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

drop table issue_attachment;