create table if not exists sys_dict_type (
  id uuid primary key default uuid_generate_v4(),
  code varchar(64) not null,
  name varchar(128) not null,
  status varchar(16) not null default 'active',
  remark varchar(512),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create unique index uk_sys_dict_type_code on sys_dict_type(code) where deleted_at is null;

create table if not exists sys_dict_item (
  id uuid primary key default uuid_generate_v4(),
  type_code varchar(64) not null,
  item_code varchar(64) not null,
  item_name varchar(128) not null,
  sort_no int not null default 0,
  status varchar(16) not null default 'active',
  remark varchar(512),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);
create unique index uk_sys_dict_item on sys_dict_item(type_code, item_code);
create index idx_sys_dict_item_type on sys_dict_item(type_code);

comment on table sys_dict_type is '系统字典类型';
comment on table sys_dict_item is '系统字典项';
