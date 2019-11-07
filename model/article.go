package model

import "time"

// タグによって各フィールド（ID や Title）にメタ情報を付与(sqlxがsql実行結果と紐付け)
//引っ張ってきたデータを当てはめる構造体を用意。
//その際、バッククオート（`）で、どのカラムと紐づけるのかを明示する。
// 構造体タグを利用して、フォームの name 属性と構造体のフィールドの紐付けを行なっている
// required(必須), max=50(最大50文字)
/* Struct(構造体) */
type Article struct {
	ID      int       `db:"id" form:"id"`
	Title   string    `db:"title" form:"title" validate:"required,max=50"`
	Body    string    `db:"body" form:"body" validate:"required"`
	Created time.Time `db:"created"`
	Updated time.Time `db:"updated"`
}
