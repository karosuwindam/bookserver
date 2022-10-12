function loadlist(table) {
    var xhr = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    xhr.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
      if(xhr.readyState == 4 && xhr.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
          var data = xhr.responseText;
          var tmp = JSON.parse(data)
          console.log(tmp.result);		          // 取得した ファイルの中身を表示
      }
    };
    var url = "/v1/read/" + table
    xhr.open('LIST',url);
    xhr.setRequestHeader('Authorization', 'bearer '+TOKEN);
    xhr.send(null);

}