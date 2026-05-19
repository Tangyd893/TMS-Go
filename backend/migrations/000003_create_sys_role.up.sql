create table if not exists sys_role (
  id uuid primary key default uuid_generate_v4(),
  code varchar(64) not null,
  name varchar(128) not null,
  remark varchar(512),
  status varchar(16) not null default 'active',
  sort_no int not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create unique index uk_sys_role_code on sys_role(code) where deleted_at is null;
comment on table sys_role is '系统角色';
