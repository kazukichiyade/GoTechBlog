package repository

import (
	// model パッケージで定義した Article 構造体を利用するため、パッケージをインポート
	"database/sql"
	"go-tech-blog/model"
	"math"
	"time"
)

/* 戻り値ありの関数 */
// func ArticleList() ([]*model.Article, error) {
// 	// 実行する SQL
// 	query := `SELECT * FROM articles;`

// 	// データベースから取得した値を格納する変数を宣言
// 	var articles []*model.Article

// 	// Query を実行して、取得した値を変数に格納
// 	if err := db.Select(&articles, query); err != nil {
// 		return nil, err
// 	}

// 	return articles, nil
// }

func ArticleListByCursor(cursor int) ([]*model.Article, error) {
	// 引数で渡されたカーソルの値が 0 以下の場合は、代わりに int 型の最大値で置き換え
	if cursor <= 0 {
		cursor = math.MaxInt32
	}

	// ID の降順に記事データを 10 件取得するクエリ文字列を生成
	query := `SELECT *
	FROM articles
	WHERE id < ?
	ORDER BY id desc
	LIMIT 10`

	// クエリ結果を格納するスライスを初期化
	// 10 件取得すると決まっているため、サイズとキャパシティを指定
	articles := make([]*model.Article, 0, 10)

	// クエリ結果を格納する変数、クエリ文字列、パラメータを指定してクエリを実行
	if err := db.Select(&articles, query, cursor); err != nil {
		return nil, err
	}

	return articles, nil
}

func ArticleCreate(article *model.Article) (sql.Result, error) {
	// 現在日時を取得
	now := time.Now()

	// 構造体に現在日時を設定
	article.Created = now
	article.Updated = now

	// クエリ文字列を生成
	query := `INSERT INTO articles (title, body, created, updated)
	VALUES (:title, :body, :created, :updated);`

	// トランザクションを開始
	tx := db.MustBegin()

	// クエリ文字列と構造体を引数に渡して SQL を実行
	// クエリ文字列内の「:title」「:body」「:created」「:updated」は構造体の値で置換
	// 構造体タグで指定してあるフィールドが対象（`db:"title"` など）
	res, err := tx.NamedExec(query, article)
	if err != nil {
		// エラーが発生した場合はロールバック
		tx.Rollback()

		// エラー内容を返却
		return nil, err
	}
	// SQL の実行に成功した場合はコミット
	tx.Commit()

	// SQL の実行結果を返却
	return res, nil
}
