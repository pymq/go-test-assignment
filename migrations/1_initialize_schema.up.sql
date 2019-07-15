create table tv
(
    id serial primary key,
    brand varchar(255) null,
    manufacturer varchar(255) not null,
    model varchar(255) not null,
    year smallint
);

create table sold_tv
(
    id serial primary key,
    tv_id integer not null references tv(id) on delete cascade,
    returned boolean not null default false,
    sale_date timestamp not null,
    quantity integer not null
);

