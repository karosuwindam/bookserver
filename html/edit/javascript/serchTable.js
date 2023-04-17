function searchGetData(output) {

  var req = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
  req.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
    if(req.readyState == 4 && req.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
        var data = req.responseText;
        var jata = JSON.parse(data);
        console.log(jata);		          // 取得した JSON ファイルの中身を表示
        document.getElementById(output).innerHTML = outputTable(jata.Result)
    }else if (req.readyState == 4 && req.status != 200){ 
        var data = req.responseText;
        var jata = JSON.parse(data);
        console.log(jata);		          // 取得した JSON ファイルの中身を表示
    }
  };
  var url = HOSTURL + "/v1/search"
  var jsondata = {};
  jsondata["Table"] = searchurl[selectdata];
  jsondata["Keyword"] = document.getElementById("search").value;
  req.open("POST", url, true); // HTTPメソッドとアクセスするサーバーの　URL　を指定
  req.send(JSON.stringify(jsondata));					    // 実際にサーバーへリクエストを送信
}

//テーブルリストの表示
function outputTable(data) {
    var output = "";
    if (data == "") {
        return output
    }
    var rowlist = table_list[selectdata];
    output += "<div class=\"table\">"
    output += "<div class=\"row\">"
    for(var j=0;j<rowlist.length;j++){
        output += "<div class=\"top-cell\">"+rowlist[j]+"</div>"
    }
    output += "</div>\n"
    for (var i=0;i<data.length;i++){
        var tmp = data[i]
        var id = tmp.Id
        output += "<div class=\"row\">"
        for(var j=0;j<rowlist.length;j++){
            output += "<div class=\"cell\">"+tmp[rowlist[j]]+"</div>"
        }
        output += "<div class=\"cell\">"+edithtml(id)+"</div>"
        output += "<div class=\"cell\">"+deletehtml(id)+"</div>"

        output += "</div>\n"
    }
    output += "</div>\n"
    return output
}


function edithtml(id){
    var output = "";
    output += "<a href=javascript:void(0); class='link-button'; onclick='vieweditforms("+"\"edit\""+","+id+");return false'>"
    output += "edit</a>"
    return output
  }
  
  function deletehtml(id){
    var output = "";
    output += "<a href=javascript:void(0); class='link-button'; onclick='destory("+id+");return false'>"
    output += "destory</a>"
    return output
  }