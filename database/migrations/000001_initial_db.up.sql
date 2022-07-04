create table if not exists contestant
(
    username varchar not null,
    user_id  varchar not null
        constraint contestant_pk
            primary key,
    team_id  varchar
);