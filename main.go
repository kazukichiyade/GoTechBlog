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

	// Custom Validation(バリデーションチェック用のライブラリ)
	"gopkg.in/go-playground/validator.v9"
)

/* 実行順(依存パッケージの読み込み > グローバル定数 > グローバル変数 > init() > main() の順に実行) */

var db *sqlx.DB

/* グローバル変数 */
var e = createMux()

func main() {
	// 自前実装の関数である connectDB() の戻り値をグローバル変数に格納
	db = connectDB()
	repository.SetDB(db)

	// TOP ページに記事の一覧を表示
	// `/` というパス（URL）と `articleIndex` という処理を結びつける(ルーティング追加)
	/* handlerパッケージのArticle〇〇〇を呼び出し */
	e.GET("/", handler.ArticleIndex)

	// 記事に関するページは "/articles" で開始するようにする
	// 記事一覧画面には "/" と "/articles" の両方でアクセスできるようにする
	// パスパラメータの ":id" も ":articleID" と明確にしている
	e.GET("/articles", handler.ArticleIndex)                // 一覧画面
	e.GET("/articles/new", handler.ArticleNew)              // 新規作成画面
	e.GET("/articles/:articleID", handler.ArticleShow)      // 詳細画面
	e.GET("/articles/:articleID/edit", handler.ArticleEdit) // 編集画面

	// HTML ではなく JSON を返却する処理は "/api" で開始
	// 記事に関する処理なので "/articles"
	e.GET("/api/articles", handler.ArticleList)                 // 一覧
	e.POST("/api/articles", handler.ArticleCreate)              // 作成
	e.DELETE("/api/articles/:articleID", handler.ArticleDelete) //削除
	e.PATCH("/api/articles/:articleID", handler.ArticleUpdate)  //更新

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
	// CSRF対策(クロス・サイト・リクエスト・フォージェリ)
	e.Use(middleware.CSRF())

	// `src/css` ディレクトリ配下のファイルに `/css` のパスでアクセスできるようにする
	e.Static("/css", "src/css")
	e.Static("/js", "src/js")

	// バリデーションを使うのに必要
	e.Validator = &CustomValidator{validator: validator.New()}

	// アプリケーションインスタンスを返却
	return e
}

// 開発用から本番用に切り替える場合でもソースコード自体を変更する必要は無くなる
func connectDB() *sqlx.DB {
	// os パッケージの Getenv() 関数を利用して環境変数から取得
	dsn := os.Getenv("DSN")
	db, err := sqlx.Open("mysql", dsn)
	// DBへ接続が失敗した場合エラーをだす処理(if文)
	/* 例外処理 */
	if err != nil {
		e.Logger.Fatal(err)
	}
	/* 簡易文付きif文 */
	if err := db.Ping(); err != nil {
		e.Logger.Fatal(err)
	}
	// DBへ接続が成功した場合の処理
	log.Println("db connection succeeded")
	return db
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
