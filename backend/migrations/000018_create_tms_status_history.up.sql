create table if not exists tms_status_history (
  id uuid primary key default uuid_generate_v4(),
  biz_type varchar(32) not null,
  biz_id uuid not null,
  from_status varchar(32),
  to_status varchar(32) not null,
  action varchar(32),
  operator_id uuid,
  operator_name varchar(64),
  remark varchar(512),
  created_at timestamptz not null default now()
);
create index idx_tms_status_history_biz on tms_status_history(biz_type, biz_id);
comment on table tms_status_history is '状态变更历史';
