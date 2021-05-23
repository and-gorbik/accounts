CREATE TABLE IF NOT EXISTS city (
    id      int not null,
    name    varchar(50) not null
);

CREATE TABLE IF NOT EXISTS country (
    id      int not null,
    name    varchar(50) not null
);

CREATE TABLE IF NOT EXISTS interest (
    account_id  int not null,
    name        varchar(100) not null
);

CREATE TABLE IF NOT EXISTS likes (
    liker_id    int not null,
    likee_id    int not null,
    ts          timestamp not null
);

CREATE TABLE IF NOT EXISTS account (
    id          int not null,
    email       varchar(100) not null,
    sex         varchar(1) not null,
    birth       timestamp not null,
    joined      timestamp not null,
    status      varchar(10) not null,
    name        varchar(50) default null,
    surname     varchar(50) default null,
    phone       varchar(16) default null,
    country_id  int default null,
    city_id     int default null,
    prem_start  timestamp default null,
    prem_end    timestamp default null
);