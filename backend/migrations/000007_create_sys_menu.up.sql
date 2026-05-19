create table if not exists sys_menu (
  id uuid primary key default uuid_generate_v4(),
  parent_id uuid,
  name varchar(128) not null,
  path varchar(256),
  component varchar(256),
  icon varchar(64),
  permission_code varchar(128),
  type varchar(16) not null default 'menu',
  sort_no int not null default 0,
  visible boolean not null default true,
  status varchar(16) not null default 'active',
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create index idx_sys_menu_parent_id on sys_menu(parent_id);
create index idx_sys_menu_sort_no on sys_menu(sort_no);
comment on table sys_menu is '系统菜单';
comment on column sys_menu.type is 'directory: 目录, menu: 菜单, button: 按钮';
