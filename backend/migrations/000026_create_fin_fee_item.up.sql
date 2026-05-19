create table if not exists fin_fee_item (
  id uuid primary key default uuid_generate_v4(),
  code varchar(64) not null,
  name varchar(128) not null,
  fee_type varchar(16) not null default 'transport',
  unit_price numeric(18,2) not null default 0,
  charge_unit varchar(32) not null default 'per_order',
  status varchar(16) not null default 'active',
  remark varchar(512),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create unique index uk_fin_fee_item_code on fin_fee_item(code) where deleted_at is null;
comment on table fin_fee_item is '费用项目';
comment on column fin_fee_item.fee_type is 'transport:运输费, loading:装卸费, waiting:等待费, toll:过路费, insurance:保险费, claim:赔付';
comment on column fin_fee_item.charge_unit is 'per_order:按单, per_kg:按重量, per_volume:按体积, per_km:按里程, fixed:固定';
