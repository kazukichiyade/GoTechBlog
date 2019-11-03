package main

import (
	// HTTPを扱うパッケージ(標準パッケージ)
	"net/http"

	// GolangのWeb FWでAPIサーバーによく使われる(外部パッケージ)
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// 実行順(依存パッケージの読み込み > グローバル定数 > グローバル変数 > init() > main() の順に実行)

var e = createMux()

func main() {
	// `/` というパス（URL）と `articleIndex` という処理を結びつける
	e.GET("/", articleIndex)

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

	// アプリケーションインスタンスを返却
	return e
}

func articleIndex(c echo.Context) error {
	// ステータスコード 200 で、"Hello, World!" という文字列をレスポンス
	return c.String(http.StatusOK, "Hello, World!")
}
