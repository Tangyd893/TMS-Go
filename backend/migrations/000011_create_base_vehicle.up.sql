create table if not exists base_vehicle (
  id uuid primary key default uuid_generate_v4(),
  plate_no varchar(32) not null,
  vehicle_type varchar(32),
  brand_model varchar(64),
  color varchar(16),
  max_load numeric(12,2),
  max_volume numeric(12,3),
  owner_type varchar(16),
  carrier_id uuid,
  driver_id uuid,
  status varchar(16) not null default 'active',
  remark varchar(512),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create unique index uk_base_vehicle_plate on base_vehicle(plate_no) where deleted_at is null;
create index idx_base_vehicle_carrier on base_vehicle(carrier_id);
comment on table base_vehicle is '车辆档案';
