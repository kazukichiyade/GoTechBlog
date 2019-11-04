package handler

import (
	// repository パッケージを利用するためインポート
	"go-tech-blog/repository"
	"log"

	// HTTPを扱うパッケージ(標準パッケージ)
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// ArticleIndex ...
// ハンドラ関数 テンプレートファイルとデータを指定して render() 関数を呼び出し
func ArticleIndex(c echo.Context) error {
	// データベースから記事データの一覧を取得する
	articles, err := repository.ArticleList()
	// データベース操作でエラーが発生した場合の処理(500)
	if err != nil {
		log.Println(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}
	data := map[string]interface{}{
		// HTMLでこれを使って表示する{{  }}
		"Message":  "Article Index Updated",
		"Now":      time.Now(),
		"Articles": articles, // 記事データをテンプレートエンジンに渡す
	}
	return render(c, "article/index.html", data)
}

// ArticleNew ...
func ArticleNew(c echo.Context) error {
	data := map[string]interface{}{
		"Message": "Article New",
		"Now":     time.Now(),
	}
	return render(c, "article/new.html", data)
}

// ArticleShow ...
func ArticleShow(c echo.Context) error {
	// パスパラメータを抽出(id=999でアクセスがあった場合c.Param("id")によって取り出す)
	// c.Param()で取り出した値は文字列型になるのでstrconvパッケージのAtoi()関数で数値型にキャスト
	id, _ := strconv.Atoi(c.Param("id"))
	data := map[string]interface{}{
		"Message": "Article Show",
		"Now":     time.Now(),
		"ID":      id,
	}
	return render(c, "article/show.html", data)
}

// ArticleEdit ...
func ArticleEdit(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	data := map[string]interface{}{
		"Message": "Article Edit",
		"Now":     time.Now(),
		"ID":      id,
	}
	return render(c, "article/edit.html", data)
}
