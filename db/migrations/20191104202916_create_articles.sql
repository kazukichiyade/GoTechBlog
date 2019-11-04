-- goose upコマンドで記載されているSQL文が実行
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE articles (
  id int AUTO_INCREMENT,
  title varchar(100),
  PRIMARY KEY(id)
);

-- goose downコマンドで記載されているSQL文が実行
-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE articles;
