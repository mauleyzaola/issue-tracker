
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

create or replace view view_projects as
select p.id,p.pkey,p.name,p.datecreated,p.begins,p.ends,
	(l.name || ' '::text) || l.lastname AS projectlead,p.issuecount,p.notresolvedcount,
	coalesce(1::double precision - p.notresolvedcount::double precision / NULLIF(p.issuecount, 0), 0)::double precision AS percentagecompleted,
	coalesce(ps.name, '') as permissionscheme,
	p.idprojectlead, coalesce(p.idpermissionscheme	, uuid_generate_v4()) as idpermissionscheme
FROM project p
JOIN users l ON l.id = p.idprojectlead
left outer join permission_scheme ps on ps.id = p.idpermissionscheme;



-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

drop view view_projects;