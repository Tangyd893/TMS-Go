create table if not exists tms_transport_task (
  id uuid primary key default uuid_generate_v4(),
  task_no varchar(32) not null,
  order_id uuid not null,
  dispatch_plan_id uuid,
  carrier_id uuid,
  vehicle_id uuid,
  driver_id uuid,
  driver_name varchar(64),
  plate_no varchar(32),
  origin_name varchar(128),
  dest_name varchar(128),
  plan_depart_time timestamptz,
  plan_arrive_time timestamptz,
  actual_depart_time timestamptz,
  actual_arrive_time timestamptz,
  status varchar(32) not null default 'pending',
  remark varchar(512),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create unique index uk_tms_transport_task_no on tms_transport_task(task_no) where deleted_at is null;
create index idx_tms_transport_task_order on tms_transport_task(order_id);
create index idx_tms_transport_task_status on tms_transport_task(status);
create index idx_tms_transport_task_driver on tms_transport_task(driver_id);
create index idx_tms_transport_task_vehicle on tms_transport_task(vehicle_id);
comment on table tms_transport_task is '运输任务';
