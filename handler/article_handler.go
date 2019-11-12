package handler

import (
	// repository パッケージを利用するためインポート
	"go-tech-blog/model"
	"go-tech-blog/repository"
	"strings"

	// HTTPを扱うパッケージ(標準パッケージ)
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// ハンドラ関数 テンプレートファイルとデータを指定して render() 関数を呼び出し
func ArticleIndex(c echo.Context) error {
	// "/articles" のパスでリクエストがあったら "/" にリダイレクト
	// Google Analytics などでのアクセス解析時にパスが統一されて分析がしやすくなる
	if c.Request().URL.Path == "/articles" {
		c.Redirect(http.StatusPermanentRedirect, "/")
	}

	// リポジトリの処理を呼び出して記事の一覧データを取得
	articles, err := repository.ArticleListByCursor(0)

	// データベース操作でエラーが発生した場合の処理(500)
	// エラーが発生した場合
	if err != nil {
		// エラー内容をサーバーのログに出力
		c.Logger().Error(err.Error())

		// クライアントにステータスコード 500 でレスポンスを返す
		return c.NoContent(http.StatusInternalServerError)
	}

	// 取得できた最後の記事の ID をカーソルとして設定
	var cursor int
	if len(articles) != 0 {
		cursor = articles[len(articles)-1].ID
	}

	// テンプレートに渡すデータを map に格納
	data := map[string]interface{}{
		// HTMLでこれを使って表示する{{  }}
		"Articles": articles, // 記事データをテンプレートエンジンに渡す
		"Cursor":   cursor,
	}

	// テンプレートファイルとデータを指定して HTML を生成し、クライアントに返却
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

func ArticleShow(c echo.Context) error {
	// パスパラメータから記事IDを取得(id=999でアクセスがあった場合c.Param("id")によって取り出す)
	// c.Param()で取り出した値は文字列型になるのでstrconvパッケージのAtoi()関数で数値型にキャスト
	id, _ := strconv.Atoi(c.Param("articleID"))

	// 記事データを取得
	article, err := repository.ArticleGetByID(id)

	if err != nil {
		// エラー内容をサーバーのログに出力
		c.Logger().Error(err.Error())

		// ステータスコード500でレスポンスを返却
		return c.NoContent(http.StatusInternalServerError)
	}

	// テンプレートに渡すデータをmapに格納
	data := map[string]interface{}{
		"Article": article,
	}
	// テンプレートファイルとデータを指定して HTML を生成し、クライアントに返却
	return render(c, "article/show.html", data)
}

// ArticleEdit ...
func ArticleEdit(c echo.Context) error {
	// パスパラメータから記事 ID を取得
	// 文字列型で取得されるので、strconv パッケージを利用して数値型にキャスト
	id, _ := strconv.Atoi(c.Param("articleID"))

	// 編集フォームの初期値として表示するために記事データを取得
	article, err := repository.ArticleGetByID(id)

	if err != nil {
		// エラー内容をサーバーのログに出力
		c.Logger().Error(err.Error())

		// ステータスコード 500 でレスポンスを返却
		return c.NoContent(http.StatusInternalServerError)
	}

	// テンプレートに渡すデータを map に格納
	data := map[string]interface{}{
		"Article": article,
	}

	// テンプレートファイルとデータを指定して HTML を生成し、クライアントに返却
	return render(c, "article/edit.html", data)
}

type ArticleCreateOutput struct {
	Article          *model.Article
	Message          string
	ValidationErrors []string
}

func ArticleCreate(c echo.Context) error {
	// 送信されてくるフォームの内容を格納する構造体を宣言
	var article model.Article

	// レスポンスとして返却する構造体を宣言
	var out ArticleCreateOutput

	// フォームの内容を構造体に埋め込む
	if err := c.Bind(&article); err != nil {
		c.Logger().Error(err.Error())

		// リクエストの解釈に失敗した場合は 400 エラーを返却
		return c.JSON(http.StatusBadRequest, out)
	}

	// バリデーションチェックを実行
	if err := c.Validate(&article); err != nil {
		// エラーの内容をサーバーのログに出力
		c.Logger().Error(err.Error())

		// エラー内容を検査してカスタムエラーメッセージを取得
		out.ValidationErrors = article.ValidationErrors(err)

		// 解釈できたパラメータが許可されていない値の場合は422エラーを返却
		return c.JSON(http.StatusUnprocessableEntity, out)
	}

	// repository を呼び出して保存処理を実行
	res, err := repository.ArticleCreate(&article)
	if err != nil {
		// エラーの内容をサーバーのログに出力
		c.Logger().Error(err.Error())

		// サーバー内の処理でエラーが発生した場合は 500 エラーを返却
		return c.JSON(http.StatusInternalServerError, out)
	}

	// SQL 実行結果から作成されたレコードの ID を取得
	id, _ := res.LastInsertId()

	// 構造体に ID をセット
	article.ID = int(id)

	// レスポンスの構造体に保存した記事のデータを格納
	out.Article = &article

	// 処理成功時はステータスコード 200 でレスポンスを返却
	return c.JSON(http.StatusOK, out)
}

func ArticleDelete(c echo.Context) error {
	// パスパラメータから記事 ID を取得
	// 文字列型で取得されるので、strconv パッケージを利用して数値型にキャスト
	id, _ := strconv.Atoi(c.Param("articleID"))

	// repositoryの記事削除処理を呼び出し
	if err := repository.ArticleDelete(id); err != nil {
		// サーバーのログにエラー内容を出力
		c.Logger().Error(err.Error())

		// サーバーサイドでエラーが発生した場合は500エラーを返却
		return c.JSON(http.StatusInternalServerError, "")
	}

	// 成功時はステータスコード200を返却
	return c.JSON(http.StatusOK, fmt.Sprintf("Article %d is deleted.", id))
}

func ArticleList(c echo.Context) error {
	// クエリパラメータからカーソルの値を取得
	// 文字列型で取得できるので strconv パッケージを用いて数値型にキャスト
	cursor, _ := strconv.Atoi(c.QueryParam("cursor"))

	// リポジトリの処理を呼び出して記事の一覧データを取得
	// 引数にカーソルの値を渡して、ID のどの位置から 10 件取得するかを指定
	articles, err := repository.ArticleListByCursor(cursor)

	// エラーが発生した場合
	if err != nil {
		// サーバーのログにエラー内容を出力
		c.Logger().Error(err.Error())

		// クライアントにステータスコード 500 でレスポンスを返却
		// HTML ではなく JSON 形式でデータのみを返却するため、
		// c.HTMLBlob() ではなく c.JSON() を呼び出し
		return c.JSON(http.StatusInternalServerError, "")
	}

	// エラーがない場合は、ステータスコード 200 でレスポンスを返却
	// JSON 形式で返却するため、c.HTMLBlob() ではなく c.JSON() を呼び出し
	return c.JSON(http.StatusOK, articles)
}

type ArticleUpdateOutput struct {
	Article          *model.Article
	Message          string
	ValidationErrors []string
}

func ArticleUpdate(c echo.Context) error {
	// リクエスト送信元のパスを取得
	ref := c.Request().Referer()

	// リクエスト送信元のパスから記事 ID を抽出
	refID := strings.Split(ref, "/")[4]

	// リクエスト URL のパスパラメータから記事 ID を抽出
	reqID := c.Param("articleID")

	// 編集画面で表示している記事と更新しようとしている記事が異なる場合は、
	// 更新処理をせずに 400 エラーを返却
	if reqID != refID {
		return c.JSON(http.StatusBadRequest, "")
	}

	// フォームで送信される記事データを格納する構造体を宣言
	var article model.Article

	// レスポンスするデータの構造体を宣言
	var out ArticleUpdateOutput

	// フォームで送信されたデータを変数に格納
	if err := c.Bind(&article); err != nil {
		// リクエストのパラメータの解釈に失敗した場合は 400 エラーを返却
		return c.JSON(http.StatusBadRequest, out)
	}

	// 入力値のチェック（バリデーションチェック）を行う
	if err := c.Validate(&article); err != nil {
		// エラー内容をレスポンスのフィールドに格納
		out.ValidationErrors = article.ValidationErrors(err)

		// 解釈できたパラメータが不正な値の場合は 422 エラーを返却
		return c.JSON(http.StatusUnprocessableEntity, out)
	}

	// 文字列型の ID を数値型にキャスト
	articleID, _ := strconv.Atoi(reqID)

	// フォームデータを格納した構造体に ID をセット
	article.ID = articleID

	// 記事を更新する処理を呼び出す
	_, err := repository.ArticleUpdate(&article)

	if err != nil {
		// レスポンスの構造体にエラー内容をセット
		out.Message = err.Error()

		// リクエスト自体は正しいにも関わらずサーバー側で処理が失敗した場合は 500 エラーを返却
		return c.JSON(http.StatusInternalServerError, out)
	}

	// レスポンスの構造体に記事データをセット
	out.Article = &article

	// 処理成功時はステータスコード 200 でレスポンスを返却
	return c.JSON(http.StatusOK, out)
}
