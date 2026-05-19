create table if not exists base_station (
  id uuid primary key default uuid_generate_v4(),
  code varchar(64) not null,
  name varchar(128) not null,
  province varchar(64),
  city varchar(64),
  district varchar(64),
  address varchar(256),
  contact_name varchar(64),
  contact_phone varchar(32),
  latitude numeric(10,7),
  longitude numeric(10,7),
  status varchar(16) not null default 'active',
  remark varchar(512),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create unique index uk_base_station_code on base_station(code) where deleted_at is null;
comment on table base_station is '站点';
