create table if not exists drink
(
    name varchar not null
        constraint drink_pk
            primary key
);

create unique index if not exists drink_name_uindex
    on drink (name);