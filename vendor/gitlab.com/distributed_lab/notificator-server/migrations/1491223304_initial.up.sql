create table requests (
  id serial primary key,
  created_at timestamp without time zone default current_timestamp,

  type int not null,
  token varchar(255) not null,
  payload jsonb not null,

  priority int default 0
)