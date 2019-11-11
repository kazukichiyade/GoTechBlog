'use strict';

// DOM Tree の構築が完了したら処理を開始
document.addEventListener('DOMContentLoaded', () => {
  // DOM API を利用して HTML 要素を取得
  const inputs = document.getElementsByTagName('input');
  const form = document.forms.namedItem('article-form');
  const saveBtn = document.querySelector('.article-form__save');
  const cancelBtn = document.querySelector('.article-form__cancel');
  const previewOpenBtn = document.querySelector('.article-form__open-preview');
  const previewCloseBtn = document.querySelector('.article-form__close-preview');
  const articleFormBody = document.querySelector('.article-form__body');
  const articleFormPreview = document.querySelector('.article-form__preview');
  const articleFormBodyTextArea = document.querySelector('.article-form__input--body');
  const articleFormPreviewTextArea = document.querySelector('.article-form__preview-body-contents');
  const errors = document.querySelector('.article-form__errors');
  const errorTmpl = document.querySelector('.article-form__error-tmpl').firstElementChild;

  // CSRF トークンを取得
  const csrfToken = document.getElementsByName('csrf')[0].content;

  // 新規作成画面か編集画面かを URL から判定
  const mode = { method: '', url: '' };
  if (window.location.pathname.endsWith('new')) {
    // 新規作成時の HTTP メソッドは POST を利用
    mode.method = 'POST';
    // 作成リクエスト、および戻るボタンの遷移先のパスは "/articles" になる
    mode.url = '/articles';
  } else if (window.location.pathname.endsWith('edit')) {
    // 更新時の HTTP メソッドは PATCH を利用
    mode.method = 'PATCH';
    // 更新リクエスト、および戻るボタンの遷移先のパスは "/articles/:articleID" になる
    mode.url = `/articles/${window.location.pathname.split('/')[2]}`;
  }
  const { method, url } = mode;

  // input 要素にフォーカスが合った状態で Enter が押されると form が送信
  // 今回は Enter キーで form が送信されないように挙動を制御
  for (let elm of inputs) {
    elm.addEventListener('keydown', event => {
      if (event.keyCode && event.keyCode === 13) {
        // キーが押された際のデフォルトの挙動をキャンセル
        event.preventDefault();

        // 何もせず処理を終了
        return false;
      }
    });
  }

  // プレビューを開くイベントを設定
  previewOpenBtn.addEventListener('click', event => {
    // form の「本文」に入力された Markdown を HTML に変換してプレビューに埋め込む
    articleFormPreviewTextArea.innerHTML = md.render(articleFormBodyTextArea.value);

    // 入力フォームを非表示
    articleFormBody.style.display = 'none';

    // プレビューを表示
    articleFormPreview.style.display = 'grid';
  });

  // プレビューを閉じるイベントを設定
  previewCloseBtn.addEventListener('click', event => {
    // 入力フォームを表示
    articleFormBody.style.display = 'grid';

    // プレビューを非表示
    articleFormPreview.style.display = 'none';
  });

  // 前のページに戻るイベントを設定
  cancelBtn.addEventListener('click', event => {
    // <button> 要素クリック時のデフォルトの挙動をキャンセル
    event.preventDefault();

    // URL を指定して画面を遷移
    window.location.href = url;
  });

  // 保存処理を実行するイベントを設定
  saveBtn.addEventListener('click', event => {
    event.preventDefault();

    // 前回のバリデーションエラーの表示が残っている場合は削除
    errors.innerHTML = null;

    // フォームに入力された内容を取得
    const fd = new FormData(form);

    let status;

    // fetch API(非同期リクエスト) を利用してリクエストを送信
    fetch(`/api${url}`, {
      method: method,
      headers: { 'X-CSRF-Token': csrfToken },
      body: fd
    })
      .then(res => {
        status = res.status;
        return res.json();
      })
      .then(body => {
        console.log(JSON.stringify(body));

        if (status === 200) {
          // 成功時は一覧画面に遷移
          window.location.href = url;
        }

        if (body.ValidationErrors) {
          // バリデーションエラーがある場合の処理をここに記載
          showErrors(body.ValidationErrors);
        }
      })
      .catch(err => console.error(err));
  });

  // バリデーションエラーを表示する関数
  const showErrors = messages => {
    // 引数の値が配列であることを確認します。
    if (Array.isArray(messages) && messages.length != 0) {
      // 複数メッセージを格納するためのフラグメントを作成します。
      const fragment = document.createDocumentFragment();

      // メッセージをループ処理します。
      messages.forEach(message => {
        // 単一メッセージを格納するためのフラグメントを作成します。
        const frag = document.createDocumentFragment();

        // テンプレートをクローンしてフラグメントに追加します。
        frag.appendChild(errorTmpl.cloneNode(true));

        // エラー要素にメッセージをセットします。
        frag.querySelector('.article-form__error').innerHTML = message;

        // エラー要素を親フラグメントに追加します。
        fragment.appendChild(frag);
      });

      // エラーメッセージの表示エリア（要素）にメッセージを追加します。
      errors.appendChild(fragment);
    }
  };
});
