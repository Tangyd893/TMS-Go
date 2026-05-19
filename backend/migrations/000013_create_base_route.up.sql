create table if not exists base_route (
  id uuid primary key default uuid_generate_v4(),
  code varchar(64) not null,
  name varchar(128) not null,
  origin_station_id uuid,
  origin_name varchar(128),
  dest_station_id uuid,
  dest_name varchar(128),
  distance_km numeric(10,2),
  estimated_hours numeric(6,1),
  status varchar(16) not null default 'active',
  remark varchar(512),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create unique index uk_base_route_code on base_route(code) where deleted_at is null;
comment on table base_route is '运输线路';
