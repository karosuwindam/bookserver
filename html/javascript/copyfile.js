function ckboxupdate () {
    el_copyckbox = document.getElementsByClassName("copyckbox")
    for (var i=0;i < el_copyckbox.length;i++){
        var id = el_copyckbox[i].innerText-0;
        if ((id != NaN) && (id != 0)){
            viewckbox(el_copyckbox[i],id)
        }
    }
}

function viewckbox(output,id) {
    var req = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    req.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
      if(req.readyState == 4 && req.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
          var data = req.responseText;
          var jata = JSON.parse(data);
          console.log(jata);		          // 取得した JSON ファイルの中身を表示
          output.innerHTML = creatckbox(jata.Result[0],id)
      }else if (req.readyState == 4 && req.status != 200){ 
          var data = req.responseText;
          var jata = JSON.parse(data);
          console.log(jata);		          // 取得した JSON ファイルの中身を表示
          output.innerHTML = creatckbox(null,id)

      }
    };
    var url = HOSTURL + "/v1/copy/" + id
    req.open("GET", url, false); // HTTPメソッドとアクセスするサーバーの　URL　を指定
    req.send(null);					    // 実際にサーバーへリクエストを送信
}

function creatckbox(jdata,id) {
    var output = "<input type=\"checkbox\" onclick=\"sendckbox(this,"+id+")\""
    if (jdata == null) {

    }else{
        if (jdata.Copyflag == 1) {
            output += "checked"
        }
    }
    output += ">"    
    return output
}

function sendckbox(input, id) {
    var jsondata = {};
    jsondata["Id"] = id -0;
    jsondata["Flag"] = input.checked;
    var req = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    req.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
      if(req.readyState == 4 && req.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
          var data = req.responseText;
          var jata = JSON.parse(data);
          console.log(jata);		          // 取得した JSON ファイルの中身を表示
      }else if (req.readyState == 4 && req.status != 200){ 
          var data = req.responseText;
          var jata = JSON.parse(data);
          console.log(jata);		          // 取得した JSON ファイルの中身を表示
          output.innerHTML = creatckbox(null,id)

      }
    };
    var url = HOSTURL + "/v1/copy"
    req.open("POST", url, false); // HTTPメソッドとアクセスするサーバーの　URL　を指定
    req.send(JSON.stringify(jsondata));					    // 実際にサーバーへリクエストを送信
}