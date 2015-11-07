
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

-- +migrate StatementBegin

create or replace function func_project_issue_count_update(myid uuid) returns boolean
as $$
declare pr project%ROWTYPE;
	total bigint;
	notResolved bigint;
begin

select count(*) into total from issue where idproject=myid and cancelleddate is null;
select count(*) into notResolved from issue where idproject=myid and cancelleddate is null and resolveddate is null;

update project set issuecount=total, notresolvedcount = notResolved where id = myid;

return true;
 end;
$$ language plpgsql;

-- +migrate StatementEnd




-- +migrate StatementBegin

create or replace function func_trg_issue() returns trigger 
 as $$
 declare 	row issue%ROWTYPE;
		oldId project.id%TYPE;
		newId project.id%TYPE;
		dummy boolean;
 begin
   if TG_OP = 'DELETE' then
    row:= OLD;
    dummy := func_project_issue_count_update(OLD.idproject);
   elsif TG_OP = 'INSERT' then    
    row := NEW;
    dummy := func_project_issue_count_update(NEW.idproject);
   else --UPDATE
    row := NEW;
    dummy := func_project_issue_count_update(NEW.idproject);
    dummy := func_project_issue_count_update(OLD.idproject);
   end if;
 
   return row;
 end;
$$ language plpgsql;

create trigger trg_issue_project after insert or update or delete on issue
 for each row execute procedure func_trg_issue();


-- +migrate StatementEnd

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

drop trigger trg_issue_project on issue;

drop function func_trg_issue();
drop function func_project_issue_count_update(uuid);