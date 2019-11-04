package main

import (
	"log"
	"os"

	"go-tech-blog/handler"
	"go-tech-blog/repository"

	// Using MySQL driver
	_ "github.com/go-sql-driver/mysql"

	// ORM(DBから引っ張ってきたデータを構造体、マップ、スライスに当てはめる)
	"github.com/jmoiron/sqlx"

	// GolangのWeb FWでAPIサーバーによく使われる(外部パッケージ)
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

/* 実行順(依存パッケージの読み込み > グローバル定数 > グローバル変数 > init() > main() の順に実行) */

var db *sqlx.DB
var e = createMux()

func main() {
	// 自前実装の関数である connectDB() の戻り値をグローバル変数に格納
	db = connectDB()
	repository.SetDB(db)

	// `/` というパス（URL）と `articleIndex` という処理を結びつける(ルーティング追加)
	e.GET("/", handler.ArticleIndex)
	e.GET("/new", handler.ArticleNew)
	e.GET("/:id", handler.ArticleShow)
	e.GET("/:id/edit", handler.ArticleEdit)

	// Webサーバーをポート番号 8080 で起動する
	e.Logger.Fatal(e.Start(":8080"))
}

func createMux() *echo.Echo {
	// アプリケーションインスタンスを生成
	e := echo.New()

	// アプリケーションに各種ミドルウェアを設定
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())

	// `src/css` ディレクトリ配下のファイルに `/css` のパスでアクセスできるようにする
	e.Static("/css", "src/css")
	e.Static("/js", "src/js")

	// アプリケーションインスタンスを返却
	return e
}

// 開発用から本番用に切り替える場合でもソースコード自体を変更する必要は無くなる
func connectDB() *sqlx.DB {
	// os パッケージの Getenv() 関数を利用して環境変数から取得
	dsn := os.Getenv("DSN")
	db, err := sqlx.Open("mysql", dsn)
	// DBへ接続が失敗した場合エラーをだす処理(if文)
	if err != nil {
		e.Logger.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		e.Logger.Fatal(err)
	}
	// DBへ接続が成功した場合の処理
	log.Println("db connection succeeded")
	return db
}
