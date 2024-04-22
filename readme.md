# Go言語による図書サーバ

# 概要について

https://github.com/karosuwindam/bookserver2
で作っていたけど、sqlite3モジュールをばらしたので、作り直し

# 機能について

* 機能APIについて

|url|Method|説明|権限|備考|
|--|--|--|--|--|
|/health|*|現時点の動作状態を確認|||
|/v1/image/:id/:filename|GET|ファイル名とidを指定してidと連携したzipファイル名内のファイルを読み取る||
|/v1/search|POST|送信したJSONデータをもとにデータベースから検索する|
|/v1/search/:table/:keyword|GET|keywordとテーブルを指定してデータベースから検索する|
|/v1/copy/:id|GET|copyfileテーブル内をidを指定してデータを読み取る|
|/v1/copy|POST|jsonデータを受け取ってその結果からファイルを共有フォルダへ追加削除を行う|
|/v1/add/:table|POST|テーブルを指定してデータベースにデータを追加
|/v1/read/:table|GET|テーブルを指定してデータベースのデータをすべて取得
|/v1/read/:table/:id|GET|テーブルとIDを指定してデータを取得
|/v1/edit/:table/:id|GET|テーブルとIDを指定してデータを取得
|/v1/edit/:table/:id|POST|テーブルとIDを指定してデータを更新
|/v1/edit/:table/:id|DELETE|テーブルとIDを指定してデータを削除
|/v1/delete/:table/:id|GET|テーブルとIDを指定してデータを削除できるか確認
|/v1/delete/:table/:id|DELETE|テーブルとIDを指定してデータを削除
|/v1/upload|POST|ファイルを送信することでデータを特定フォルダに保管したりアップロードテーブルを更新する|
|/v1/upload/:filename|GET|ファイル名に連携したファイル出力について返す|||
|/v1/upload/:filetype/:filename|GET|ファイルの種類とファイ名を指定して特定場所に保存していることを確認|
|/v1/download/:filetype/:id|GET|対象のファイルをダウンロード|
|/v1/history|GET|/view/:idにアクセスしたときの履歴を表示|
|/v1/history?n=x|GET|/view/:idにアクセスしたときの履歴をx個表示|


# データベースにデータを追加について
テーブルのデータベースの型に沿ってJSONデータを送信する

# 設定可能な環境変数

|名前|説明|初期値|備考|
|--|--|--|--|
|WEB_PROTOCOL|プロトコル名|tcp||
|WEB_HOST|ホスト名|空白||
|WEB_PORT|解放ポート|8080||
|DB_NAME|SQLのデータベースタイプ|mysql|sqlite3も可能|
|DB_HOST|SQLのIPアドレス|127.0.0.1||
|DB_PORT|SQLの接続ポート|3306||
|DB_USER|SQLの接続ユーザ||
|DB_PASS|SQLの接続ユーザパスワード||
|DB_FILE|SQLite3の接続ファイル名|test.db|
|DBROOTPASS|SQLlite3の相対ファイルパス|./db/|
|TMP_FILEPASS|テンプレートフォルダとして用意するフォルダ|./tmp|
|PDF_FILEPASS|PDFを保存するフォルダ|./upload/pdf|
|ZIP_FILEPASS|ZIPを保存するフォルダ|./upload/zip|
|IMG_FILEPASS|imageを保存するフォルダ|./html/img|
|PUBLIC_FILEPASS|公開用のフォルダ|./public|
|MAX_MULTI_MEMORY|複数アップロード時のメモリ|256M||

## CURLによるテスト

curl localhost:8080/v1/add/booknames/  -X POST -d "name=bagaet" --data-urlencode "title=はなび" --data-urlencode "ext=げた"
curl localhost:8080/v1/add/filelists/  -X POST --data-urlencode "name=はなび" --data-urlencode "pdfpass=げた" --data-urlencode "zippass=gaega.zip"
curl localhost:8080/v1/add/copyfile/  -X POST --data-urlencode "zippass=はなび" -d "filesize=130" -d "copyflag=1"


curl -H "Authorization: bearer <token>" localhost:8080/v1/
