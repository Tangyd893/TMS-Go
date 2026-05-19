create table if not exists tms_exception_log (
  id uuid primary key default uuid_generate_v4(),
  exception_id uuid not null,
  action varchar(32) not null,
  content text,
  operator_id uuid,
  operator_name varchar(64),
  created_at timestamptz not null default now()
);
create index idx_tms_exception_log_exception on tms_exception_log(exception_id);
comment on table tms_exception_log is '异常处理日志';
