create table if not exists tms_order_cargo (
  id uuid primary key default uuid_generate_v4(),
  order_id uuid not null,
  cargo_name varchar(128),
  cargo_code varchar(64),
  cargo_type varchar(32),
  quantity int default 1,
  weight numeric(12,2),
  volume numeric(12,3),
  unit varchar(16),
  remark varchar(256),
  created_at timestamptz not null default now()
);
create index idx_tms_order_cargo_order on tms_order_cargo(order_id);
comment on table tms_order_cargo is '订单货物明细';
