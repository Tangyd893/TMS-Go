create table if not exists tms_receipt (
  id uuid primary key default uuid_generate_v4(),
  task_id uuid not null,
  order_id uuid not null,
  receipt_no varchar(32) not null,
  sign_by varchar(64),
  sign_at timestamptz,
  sign_image_url varchar(512),
  remark varchar(512),
  created_at timestamptz not null default now(),
  created_by uuid,
  updated_at timestamptz not null default now()
);
create unique index uk_tms_receipt_no on tms_receipt(receipt_no);
create index idx_tms_receipt_task on tms_receipt(task_id);
create index idx_tms_receipt_order on tms_receipt(order_id);
comment on table tms_receipt is '签收回单';
