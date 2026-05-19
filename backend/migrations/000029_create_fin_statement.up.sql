create table if not exists fin_statement (
  id uuid primary key default uuid_generate_v4(),
  statement_no varchar(32) not null,
  partner_id uuid not null,
  partner_name varchar(128),
  partner_type varchar(16) not null,
  statement_period_start date not null,
  statement_period_end date not null,
  total_receivable numeric(18,2) default 0,
  total_payable numeric(18,2) default 0,
  status varchar(32) not null default 'draft',
  confirmed_at timestamptz,
  confirmed_by uuid,
  remark varchar(512),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create unique index uk_fin_statement_no on fin_statement(statement_no) where deleted_at is null;
create index idx_fin_statement_partner on fin_statement(partner_id, partner_type);
create index idx_fin_statement_status on fin_statement(status);
comment on table fin_statement is '对账单';
comment on column fin_statement.partner_type is 'customer:客户, carrier:承运商';
comment on column fin_statement.status is 'draft:草稿, confirmed:已确认, settled:已结算';
