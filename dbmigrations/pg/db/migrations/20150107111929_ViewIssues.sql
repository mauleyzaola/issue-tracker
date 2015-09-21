
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE OR REPLACE VIEW view_issues AS 
 SELECT i.id,
    i.pkey,
    i.name,
    i.datecreated,
    i.lastmodified,
    i.duedate,
    i.resolveddate,
    i.cancelleddate,
    (r.name || ' '::text) || r.lastname AS reporter,
    pr.name AS priority,
    s.name AS status,
    COALESCE(p.name, ''::text) AS project,
    (COALESCE(a.name, ''::text) || ' '::text) || COALESCE(a.lastname, ''::text) AS assignee,
    p.id AS idproject,
    a.id AS idassignee,
    r.id AS idreporter,
    s.id AS idstatus,
    pr.id AS idpriority,
    i.idparent,
    COALESCE(pa.pkey, ''::text) AS parent
   FROM issue i
     JOIN users r ON r.id = i.idreporter
     JOIN priority pr ON pr.id = i.idpriority
     JOIN status s ON s.id = i.idstatus
     LEFT JOIN project p ON p.id = i.idproject
     LEFT JOIN users a ON a.id = i.idassignee
     LEFT JOIN issue pa ON pa.id = i.idparent;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

drop view view_issues;