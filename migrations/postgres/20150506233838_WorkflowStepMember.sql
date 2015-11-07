
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create table workflow_step_member
(
	id uuid primary key default uuid_generate_v4(),
	idworkflowstep uuid not null references workflow_step(id),
	iduser uuid null references users(id),
	idgroup uuid null references groups(id)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

drop table workflow_step_member;