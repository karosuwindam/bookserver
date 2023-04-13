
function createeditform(id,table,output) {
    var xhr = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    xhr.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
      if(xhr.readyState == 4 && xhr.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
          var data = xhr.responseText;
          var tmp = JSON.parse(data)
          document.getElementById(output).innerHTML = createeditforminput(table,tmp.Result[0],output);
          console.log(tmp.Result);		          // 取得した ファイルの中身を表示
      }
    };
    var url = "/v1/edit/" + table + "/" + id
    xhr.open('GET',url);
    xhr.setRequestHeader('Authorization', 'bearer '+TOKEN);
    xhr.send(null);

}

function geteditjsondata(id,table) {
    var output="";
    var jsondata = {};
    var table_listdata = table_list[selectdata] ;
    for (var i=0;i<table_listdata.length;i++){
        if (table_listdata[i] == "Id") {
            jsondata[table_listdata[i]] = id
        }else {
            jsondata[table_listdata[i]] = document.getElementById(table+"_"+table_listdata[i]).value
            switch (table_listdata[i]){
                case "filesize":
                case "copyflag":
                    jsondata[table_listdata[i]] = jsondata[table_listdata[i]] -0;
                    break;
            }
        }
    }
    output = JSON.stringify(jsondata)
    return output;
}

function sendeditdata(id,table){
    var xhr = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    xhr.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
      if(xhr.readyState == 4 && xhr.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
          var data = xhr.responseText;
          var tmp = JSON.parse(data)
          console.log(tmp.Result);		          // 取得した ファイルの中身を表示
      }
    };
    var url = "/v1/edit/" + table + "/" + id;
    var jsondata = geteditjsondata(id,table);
    xhr.open('POST',url);
    xhr.setRequestHeader('Authorization', 'bearer '+TOKEN);
    xhr.send(jsondata);
}


function closeedit(iddata){
    document.getElementById(iddata).innerHTML = "";
    reviewform();
}

function createeditforminput(table,jsondata,tag) {
    var output = "";

    var table_listdata = table_list[selectdata] 
    output += table +"<br>"
    output += "<table>"+"<tr>"+"<th>Key名</th>"+"<th>値</th>"+"</tr>"
    for (var i = 0;i<table_listdata.length;i++){
      output += "<tr>"
      output += "<td>"
      output += table_listdata[i]
      output += "</td>"
      output += "<td>"
      if (table_listdata[i] == "Id") {
        output += jsondata[table_listdata[i]]
        var id =jsondata[table_listdata[i]];
      }else{
        output += "<input type=\"text\" name=\""+table+"\" id=\""+table+"_"+table_listdata[i]+"\" value=\""+jsondata[table_listdata[i]]+"\">"

      }
      output += "</td>"
      output += "</tr>"
    }
    output += "</table>"
    output +="<input type=\"button\" value=\"push\" onclick=\"sendeditdata("+id+",'"+table+"');closeedit('"+tag+"')\">"
    output +="<input type=\"button\" value=\"close\" onclick=\"closeedit('"+tag+"')\">"
    // document.getElementById(output).innerHTML = data;

    return output;
}