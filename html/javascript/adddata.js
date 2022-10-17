function adddata(table,json_data) {
    var xhr = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    xhr.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
      if(xhr.readyState == 4 && xhr.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
          var data = xhr.responseText;
          var tmp = JSON.parse(data)
          console.log(tmp.result);		          // 取得した ファイルの中身を表示
      }
    };
    var url = "/v1/add/" + table
    xhr.open('POST',url);
    xhr.setRequestHeader('Authorization', 'bearer '+TOKEN);
    xhr.send(json_data);

}

function createaddform(table,output) {
  var data = "";
  var table_listdata = table_list[selectdata] 
  data += table +"<br>"
  data += "<table>"+"<tr>"+"<th>Key名</th>"+"<th>値</th>"+"</tr>"
  for (var i = 0;i<table_listdata.length;i++){
    if (table_listdata[i] == "id") {
      continue
    }
    data += "<tr>"
    data += "<td>"
    data += table_listdata[i]
    data += "</td>"
    data += "<td>"
    data += "<input type=\"text\" name=\""+table+"\" id=\""+table+"_"+table_listdata[i]+"\">"
    data += "</td>"
    data += "</tr>"
  }
  data += "</table>"
  data +="<input type=\"button\" value=\"push\" onclick=\"sendadddata('"+table+"');closeedit('"+output+"')\">"
  document.getElementById(output).innerHTML = data;
}