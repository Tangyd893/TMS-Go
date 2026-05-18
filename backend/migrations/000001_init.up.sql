create extension if not exists "uuid-ossp";

create table if not exists sys_operation_log (
  id uuid primary key default uuid_generate_v4(),
  trace_id varchar(64),
  operator_id uuid,
  operator_name varchar(128),
  module varchar(64) not null,
  action varchar(64) not null,
  biz_type varchar(64),
  biz_id uuid,
  request_method varchar(16),
  request_path varchar(256),
  client_ip varchar(64),
  user_agent varchar(512),
  result varchar(32),
  error_message text,
  created_at timestamptz not null default now()
);
