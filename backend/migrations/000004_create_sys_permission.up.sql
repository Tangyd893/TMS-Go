create table if not exists sys_permission (
  id uuid primary key default uuid_generate_v4(),
  code varchar(128) not null,
  name varchar(128) not null,
  type varchar(32) not null,
  parent_id uuid,
  sort_no int not null default 0,
  remark varchar(512),
  status varchar(16) not null default 'active',
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create unique index uk_sys_permission_code on sys_permission(code) where deleted_at is null;
comment on table sys_permission is '系统权限';
comment on column sys_permission.type is 'menu: 菜单, button: 按钮, api: 接口';
