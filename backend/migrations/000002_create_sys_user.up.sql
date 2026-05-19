create table if not exists sys_user (
  id uuid primary key default uuid_generate_v4(),
  username varchar(64) not null,
  password_hash varchar(256) not null,
  real_name varchar(64),
  phone varchar(32),
  email varchar(128),
  avatar_url varchar(512),
  status varchar(16) not null default 'active',
  org_id uuid,
  last_login_at timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create unique index uk_sys_user_username on sys_user(username) where deleted_at is null;
create index idx_sys_user_status on sys_user(status);
comment on table sys_user is '系统用户';
comment on column sys_user.status is 'active: 启用, disabled: 禁用';
