package model

// Article ...
// タグによって各フィールド（ID や Title）にメタ情報を付与(sqlxがsql実行結果と紐付け)
//引っ張ってきたデータを当てはめる構造体を用意。
//その際、バッククオート（`）で、どのカラムと紐づけるのかを明示する。
/* Struct(構造体) */
type Article struct {
	ID    int    `db:"id"`
	Title string `db:"title"`
}
