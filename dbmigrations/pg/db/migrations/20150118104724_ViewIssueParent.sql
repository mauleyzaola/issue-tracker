
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

create or replace view view_issue_parents as

with recursive cte(id, pkey, idparent, level, path) as
(
select t.id, t.pkey, t.idparent, 1::INT as level, t.id::TEXT as path 
from issue t 
where t.idparent is null
union all
select c.id, c.pkey, c.idparent, p.level + 1 as level, (p.path || ':' || c.id::TEXT) 
from cte p
join issue c on c.idparent = p.id
)
select 	t.*
from cte t;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

drop view view_issue_parents;