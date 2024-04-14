var tablelist;

function getTablelist(output){

    var req = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    req.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
      if(req.readyState == 4 && req.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
          var data = req.responseText;
          var jata = JSON.parse(data);
          console.log(jata);		          // 取得した JSON ファイルの中身を表示
          document.getElementById(output).innerHTML = outputTable(jata.Result);
          tablelist = jata.Result;
      }else if (req.readyState == 4 && req.status != 200){ 
          var data = req.responseText;
          var jata = JSON.parse(data);
          console.log(jata);		          // 取得した JSON ファイルの中身を表示
      }
    };
    var url = HOSTURL + "/v1/read/" + searchurl[selectdata];
    req.open("GET", url, true); // HTTPメソッドとアクセスするサーバーの　URL　を指定
    req.send(null);					    // 実際にサーバーへリクエストを送信
  }
  