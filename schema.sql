create extension if not exists pgcrypto;

create table if not exists users (
     id            bigserial primary key,
     login         text not null unique,
     password_hash text not null,
     created_at    timestamptz not null default now()
);

create table if not exists sessions (
    id         text primary key default encode(gen_random_bytes(16),'hex'),
    uid        bigint not null,
    ip_address text,
    created_at timestamptz not null default now()
);

create table if not exists assets (
    name       text not null,
    uid        bigint not null,
    data       bytea not null,
    created_at timestamptz not null default now(),
    primary key (name, uid)
);

-- Добавляем внешние ключи (FK), чтобы при удалении пользователя удалялись его сессии/файлы (on delete cascade).
alter table sessions
    add constraint sessions_uid_fk
    foreign key (uid) references users(id)
    on delete cascade;

alter table assets
    add constraint assets_uid_fk
    foreign key (uid) references users(id)
    on delete cascade;

-- Тестовый пользователь (login='alice', password='secret')
insert into users (login, password_hash)
values ('alice', encode(digest('secret','md5'),'hex'))
    on conflict do nothing;
