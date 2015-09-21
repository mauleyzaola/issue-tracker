
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create domain tcurrency numeric(25,2);
create domain tquantity numeric(25,4);

-- +goose StatementBegin
create extension "uuid-ossp";
-- +goose StatementEnd


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

drop domain tcurrency;
drop domain tquantity;

drop extension "uuid-ossp";