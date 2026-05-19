create table if not exists sys_user_role (
  id uuid primary key default uuid_generate_v4(),
  user_id uuid not null,
  role_id uuid not null,
  created_at timestamptz not null default now()
);
create unique index uk_sys_user_role on sys_user_role(user_id, role_id);
create index idx_sys_user_role_user_id on sys_user_role(user_id);
create index idx_sys_user_role_role_id on sys_user_role(role_id);
comment on table sys_user_role is '用户角色关联';
