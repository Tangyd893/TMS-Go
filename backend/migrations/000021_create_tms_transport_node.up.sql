create table if not exists tms_transport_node (
  id uuid primary key default uuid_generate_v4(),
  task_id uuid not null,
  node_type varchar(32) not null,
  node_name varchar(128),
  location_name varchar(256),
  longitude numeric(10,6),
  latitude numeric(10,6),
  arrived_at timestamptz,
  departed_at timestamptz,
  remark varchar(512),
  created_at timestamptz not null default now(),
  created_by uuid
);
create index idx_tms_transport_node_task on tms_transport_node(task_id);
create index idx_tms_transport_node_type on tms_transport_node(node_type);
comment on table tms_transport_node is '运输节点记录';
comment on column tms_transport_node.node_type is 'pickup:提货, depart:发车, transit:中转到达, transit_depart:中转离开, arrive:到达目的地, sign:签收';
