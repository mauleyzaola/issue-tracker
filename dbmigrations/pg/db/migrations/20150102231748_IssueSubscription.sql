
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table issue_subscription
(
  id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
  datecreated timestamp with time zone NOT NULL,
  idissue uuid NOT NULL REFERENCES issue(id),
  iduser uuid NOT NULL REFERENCES users(id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

drop table issue_subscription;