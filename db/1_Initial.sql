-- +migrate Up
CREATE TABLE products(
  id    INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE, 
  model  TEXT,
  company TEXT,
  price INTEGER
);

-- +migrate Down
DROp TABLE products;