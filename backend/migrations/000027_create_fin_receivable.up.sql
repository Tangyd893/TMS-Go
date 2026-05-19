create table if not exists fin_receivable (
  id uuid primary key default uuid_generate_v4(),
  receivable_no varchar(32) not null,
  order_id uuid not null,
  task_id uuid,
  customer_id uuid not null,
  customer_name varchar(128),
  fee_item_id uuid not null,
  fee_item_name varchar(128),
  amount numeric(18,2) not null default 0,
  has_tax boolean not null default false,
  tax_rate numeric(5,4) default 0,
  tax_amount numeric(18,2) default 0,
  total_amount numeric(18,2) not null default 0,
  status varchar(32) not null default 'pending',
  remark varchar(512),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create unique index uk_fin_receivable_no on fin_receivable(receivable_no) where deleted_at is null;
create index idx_fin_receivable_order on fin_receivable(order_id);
create index idx_fin_receivable_customer on fin_receivable(customer_id);
create index idx_fin_receivable_status on fin_receivable(status);
comment on table fin_receivable is '应收费用';
comment on column fin_receivable.status is 'pending:待确认, confirmed:已确认, settled:已结算';
