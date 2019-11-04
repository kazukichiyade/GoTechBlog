package model

// Article ...
// タグによって各フィールド（ID や Title）にメタ情報を付与(sqlxがsql実行結果と紐付け)
type Article struct {
	ID    int    `db:"id"`
	Title string `db:"title"`
}
