create table if not exists base_customer (
  id uuid primary key default uuid_generate_v4(),
  code varchar(64) not null,
  name varchar(128) not null,
  short_name varchar(64),
  contact_name varchar(64),
  contact_phone varchar(32),
  contact_email varchar(128),
  address varchar(256),
  status varchar(16) not null default 'active',
  remark varchar(512),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create unique index uk_base_customer_code on base_customer(code) where deleted_at is null;
comment on table base_customer is '客户档案';
