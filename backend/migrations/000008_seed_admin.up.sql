-- 管理员角色
insert into sys_role (id, code, name, remark, status, sort_no)
values (uuid_generate_v4(), 'admin', '超级管理员', '系统超级管理员角色', 'active', 1);

-- 管理员用户 (密码: admin123)
insert into sys_user (id, username, password_hash, real_name, status)
values (
  uuid_generate_v4(),
  'admin',
  '$2a$12$zkLVEPXoGeIotHOyEX6c3uCV.tZC4eyXS1UjSfkvFdec9CmdekzzO',
  '系统管理员',
  'active'
);

-- 关联管理员用户到管理员角色
insert into sys_user_role (id, user_id, role_id)
select uuid_generate_v4(), u.id, r.id
from sys_user u, sys_role r
where u.username = 'admin' and r.code = 'admin';

-- 基础权限
insert into sys_permission (id, code, name, type, sort_no, status)
values
  (uuid_generate_v4(), 'dashboard', '工作台', 'menu', 1, 'active'),
  (uuid_generate_v4(), 'system:user:list', '用户管理', 'menu', 2, 'active'),
  (uuid_generate_v4(), 'system:user:create', '创建用户', 'button', 1, 'active'),
  (uuid_generate_v4(), 'system:user:update', '编辑用户', 'button', 2, 'active'),
  (uuid_generate_v4(), 'system:user:delete', '删除用户', 'button', 3, 'active'),
  (uuid_generate_v4(), 'system:role:list', '角色管理', 'menu', 3, 'active'),
  (uuid_generate_v4(), 'system:menu:list', '菜单管理', 'menu', 4, 'active');

-- 关联管理员角色到所有权限
insert into sys_role_permission (id, role_id, permission_id)
select uuid_generate_v4(), r.id, p.id
from sys_role r, sys_permission p
where r.code = 'admin';

-- 基础菜单
insert into sys_menu (id, parent_id, name, path, component, icon, permission_code, type, sort_no, visible, status)
select uuid_generate_v4(), null, '工作台', '/dashboard', 'views/dashboard/DashboardView', 'Monitor', 'dashboard', 'menu', 1, true, 'active'
where not exists (select 1 from sys_menu where name = '工作台');

insert into sys_menu (id, parent_id, name, path, component, icon, permission_code, type, sort_no, visible, status)
select uuid_generate_v4(), null, '系统管理', '/system', null, 'Setting', null, 'directory', 9, true, 'active'
where not exists (select 1 from sys_menu where name = '系统管理');

insert into sys_menu (id, parent_id, name, path, component, icon, permission_code, type, sort_no, visible, status)
select uuid_generate_v4(), (select id from sys_menu where name = '系统管理' limit 1), '用户管理', '/system/user', 'views/system/UserList', null, 'system:user:list', 'menu', 1, true, 'active'
where not exists (select 1 from sys_menu where name = '用户管理');
