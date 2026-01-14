-- 初始化用户与凭证表（示例，后续可迁移至真实数据库）
create table if not exists users (
  id varchar(64) primary key,
  username varchar(255) not null unique,
  password_hash text not null,
  created_at timestamp default current_timestamp
);

create table if not exists refresh_tokens (
  token varchar(128) primary key,
  user_id varchar(64) not null,
  created_at timestamp default current_timestamp
);
