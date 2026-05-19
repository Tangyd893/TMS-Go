create table if not exists tms_dispatch_plan (
  id uuid primary key default uuid_generate_v4(),
  plan_no varchar(32) not null,
  status varchar(32) not null default 'pending',
  remark varchar(512),
  created_by uuid,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create unique index uk_tms_dispatch_plan_no on tms_dispatch_plan(plan_no) where deleted_at is null;
create index idx_tms_dispatch_plan_status on tms_dispatch_plan(status);

create table if not exists tms_dispatch_detail (
  id uuid primary key default uuid_generate_v4(),
  plan_id uuid not null,
  order_id uuid not null,
  carrier_id uuid,
  vehicle_id uuid,
  driver_id uuid,
  seq int default 1,
  created_at timestamptz not null default now()
);
create index idx_tms_dispatch_detail_plan on tms_dispatch_detail(plan_id);
create index idx_tms_dispatch_detail_order on tms_dispatch_detail(order_id);
comment on table tms_dispatch_plan is '调度计划';
comment on table tms_dispatch_detail is '调度明细';
