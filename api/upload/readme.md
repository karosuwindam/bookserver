## upload

## 設定について
設定できる環境変数リストについて

|環境変数名|説明|デフォルト|備考|
|--|--|--|--|
|PDF_FILEPASS|pdfアップロードフォルダ|./upload/pdf||
|ZIP_FILEPASS|zipアップロードフォルダ|./upload/zip||

## 主な機能について

|API|Method|説明|備考|
|--|--|--|--|
|/upload|POST|file名で指定したファイルをアップロードする||
|/upload|PUT|{"Name":"ファイル名"}でデータを送信するとファイル状態を確認する||
|/upload|GET|{"Name":"ファイル名"}でデータを送信するとファイル状態を確認する||
|/upload/|GET|動作なし||
|/upload/|POST|file名で指定したファイルをアップロードする||
|/upload/zip|LIST|zipフォルダで指定したファイルリストを取得||
|/upload/pdf|LIST|pdfフォルダで指定したファイルリストを取得||