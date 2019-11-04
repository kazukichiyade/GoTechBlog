package repository

import (
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

// 引数にデータベースとの接続情報を持った構造体を取り、repositoryパッケージのグローバル変数にセット
// repository パッケージ内でデータベースへのアクセスが可能になる
// SetDB ...
func SetDB(d *sqlx.DB) {
	db = d
}
