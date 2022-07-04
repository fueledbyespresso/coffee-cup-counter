create table if not exists tally
(
    contestant varchar                             not null
        constraint tally_contestant_user_id_fk
            references contestant
            on delete cascade,
    drink      varchar   default 'drip'            not null
        constraint tally_drink_name_fk
            references drink
            on update cascade on delete set default,
    timestamp  timestamp default current_timestamp not null
);

