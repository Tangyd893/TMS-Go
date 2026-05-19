create table if not exists fin_settlement (
  id uuid primary key default uuid_generate_v4(),
  settlement_no varchar(32) not null,
  statement_id uuid not null,
  partner_id uuid not null,
  partner_name varchar(128),
  partner_type varchar(16) not null,
  settlement_amount numeric(18,2) not null default 0,
  settlement_method varchar(32),
  status varchar(32) not null default 'pending',
  settled_at timestamptz,
  remark varchar(512),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create unique index uk_fin_settlement_no on fin_settlement(settlement_no) where deleted_at is null;
create index idx_fin_settlement_statement on fin_settlement(statement_id);
create index idx_fin_settlement_partner on fin_settlement(partner_id);
comment on table fin_settlement is '结算单';
comment on column fin_settlement.status is 'pending:待结算, completed:已完成';
