create table if not exists sys_role_permission (
  id uuid primary key default uuid_generate_v4(),
  role_id uuid not null,
  permission_id uuid not null,
  created_at timestamptz not null default now()
);
create unique index uk_sys_role_permission on sys_role_permission(role_id, permission_id);
comment on table sys_role_permission is '角色权限关联';
