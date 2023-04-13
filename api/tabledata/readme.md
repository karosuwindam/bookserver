# table操作機能

## webの主な機能について

|API|Method|説明|備考|
|--|--|--|--|
|/add/:table|POST|テーブルを指定してjsonデータを送信することでデータベースに書き込む||
|/edit/:table/:id|GET|テーブルを指定してその中のidのデータを取得する|
|/edit/:table/:id|POST|テーブルとidを指定してjsonデータを送信することでそのidを書き換える|
|/edit/:table/:id|DELETE|指定したテーブル内からidを削除する|
|/read/:table|LIST|指定したテーブルのデータを全て取得する|
|/read/:table/:id|GET|指定したテーブルからidが一致したデータを取得する|
|/search|POST|送信したjsonデータから指定したテーブルからkeywordに一致したデータを取得する|
|/search/:table/:keyword|GET|指定したテーブル内からkeywordに一致したデータを取得する|