package model

import (
	"time"

	"gopkg.in/go-playground/validator.v9"
)

// タグによって各フィールド（ID や Title）にメタ情報を付与(sqlxがsql実行結果と紐付け)
//引っ張ってきたデータを当てはめる構造体を用意。
//その際、バッククオート（`）で、どのカラムと紐づけるのかを明示する。
// 構造体タグを利用して、フォームの name 属性と構造体のフィールドの紐付けを行なっている
// required(必須), max=50(最大50文字)
/* Struct(構造体) */
type Article struct {
	ID      int       `db:"id" form:"id" json:"id"`
	Title   string    `db:"title" form:"title" validate:"required,max=50" json:"title"`
	Body    string    `db:"body" form:"body" validate:"required" json:"body"`
	Created time.Time `db:"created" json:"created"`
	Updated time.Time `db:"updated" json:"updated"`
}

func (a *Article) ValidationErrors(err error) []string {
	// メッセージを格納するスライスを宣言
	var errMessages []string

	// 複数のエラーが発生する場合があるのでループ処理
	for _, err := range err.(validator.ValidationErrors) {
		// メッセージを格納する変数を宣言
		var message string

		// エラーになったフィールドを特定
		switch err.Field() {
		case "Title":

			// エラーになったバリデーションルールを特定
			switch err.Tag() {
			case "required":
				message = "タイトルは必須です。"
			case "max":
				message = "タイトルは最大50文字です。"
			}
		case "Body":
			message = "本文は必須です。"
		}
		// メッセージをスライスに追加
		if message != "" {
			errMessages = append(errMessages, message)
		}
	}
	return errMessages

}
