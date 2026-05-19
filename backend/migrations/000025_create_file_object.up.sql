create table if not exists file_object (
  id uuid primary key default uuid_generate_v4(),
  bucket varchar(64) not null,
  object_key varchar(512) not null,
  original_name varchar(256) not null,
  content_type varchar(128),
  size bigint not null default 0,
  sha256 varchar(64),
  biz_type varchar(32),
  biz_id uuid,
  uploaded_by uuid,
  created_at timestamptz not null default now()
);
create index idx_file_object_biz on file_object(biz_type, biz_id);
create index idx_file_object_bucket on file_object(bucket);
comment on table file_object is '文件对象元数据';
