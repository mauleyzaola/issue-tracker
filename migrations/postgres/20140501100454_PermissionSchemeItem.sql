
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create table permission_scheme_item
(
	id uuid primary key default uuid_generate_v4(),
	idpermissionscheme uuid not null references permission_scheme(id),
	idpermissionname uuid not null references permission_name(id),	
	idrole uuid null references roles(id),
	idgroup uuid null references groups(id),
	iduser uuid null references users(id)
);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

drop table permission_scheme_item;