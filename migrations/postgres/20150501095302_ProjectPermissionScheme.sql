
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table project add idpermissionscheme uuid null references permission_scheme(id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

alter table project drop column idpermissionscheme;