create table if not exists base_driver (
  id uuid primary key default uuid_generate_v4(),
  code varchar(64) not null,
  name varchar(64) not null,
  id_card varchar(32),
  phone varchar(32),
  license_type varchar(16),
  license_no varchar(32),
  license_expire_date date,
  status varchar(16) not null default 'active',
  remark varchar(512),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create unique index uk_base_driver_code on base_driver(code) where deleted_at is null;
comment on table base_driver is '司机档案';
