create table if not exists tms_order (
  id uuid primary key default uuid_generate_v4(),
  order_no varchar(32) not null,
  customer_id uuid not null,
  customer_name varchar(128),
  shipper_name varchar(128),
  shipper_phone varchar(32),
  shipper_address varchar(256),
  receiver_name varchar(128),
  receiver_phone varchar(32),
  receiver_address varchar(256),
  origin_station_id uuid,
  origin_name varchar(128),
  dest_station_id uuid,
  dest_name varchar(128),
  plan_pickup_time timestamptz,
  plan_delivery_time timestamptz,
  actual_pickup_time timestamptz,
  actual_delivery_time timestamptz,
  cargo_name varchar(128),
  cargo_weight numeric(12,2),
  cargo_volume numeric(12,3),
  cargo_quantity int,
  transport_requirement varchar(512),
  status varchar(32) not null default 'draft',
  remark varchar(512),
  created_by uuid,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create unique index uk_tms_order_no on tms_order(order_no) where deleted_at is null;
create index idx_tms_order_customer on tms_order(customer_id);
create index idx_tms_order_status on tms_order(status);
create index idx_tms_order_plan_pickup on tms_order(plan_pickup_time);
comment on table tms_order is '运输订单';
