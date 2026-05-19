create table if not exists tms_exception (
  id uuid primary key default uuid_generate_v4(),
  task_id uuid not null,
  order_id uuid not null,
  exception_no varchar(32) not null,
  exception_type varchar(64) not null,
  description text not null,
  severity varchar(16) not null default 'normal',
  report_by uuid not null,
  report_by_name varchar(64),
  handler_id uuid,
  handler_name varchar(64),
  handle_result text,
  handle_at timestamptz,
  status varchar(32) not null default 'pending',
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create unique index uk_tms_exception_no on tms_exception(exception_no) where deleted_at is null;
create index idx_tms_exception_task on tms_exception(task_id);
create index idx_tms_exception_order on tms_exception(order_id);
create index idx_tms_exception_status on tms_exception(status);
comment on table tms_exception is '运输异常';
comment on column tms_exception.status is 'pending:待处理, processing:处理中, resolved:已解决, closed:已关闭';
comment on column tms_exception.severity is 'normal:一般, serious:严重, critical:紧急';
