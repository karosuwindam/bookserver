var searchurl = ["booknames","filelists","copyfile"]
var bookname_row = ["id","name","title","writer","burand","booktype","ext"];
var copyfile_row = ["id","zippass","filesize","copyflag"];
var filelists_row =["id","name","pdfpass","zippass"];
var table_list = [bookname_row,filelists_row,copyfile_row];
var selectdata = 0

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

function loadlistdata(table,output) {
  var xhr = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
  xhr.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
    if(xhr.readyState == 4 && xhr.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
        var data = xhr.responseText;
        var tmp = JSON.parse(data)
        console.log(tmp.result);		          // 取得した ファイルの中身を表示
        if (table == `listdata`){
          document.getElementById(output).innerHTML = listoutput(data);
        }else if (tmp.result != "[]"){
          document.getElementById(output).innerHTML = jsonDataOutput(tmp.result);
        }else{
          document.getElementById(output).innerHTML = "";
        }
  }
  };
  var url = "/v1/read/" + table
  xhr.open('LIST',url);
  xhr.setRequestHeader('Authorization', 'bearer '+TOKEN);
  xhr.send(null);

}

function searchkeyword(keyword,outputid) {
  var xhr = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
  xhr.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
    if(xhr.readyState == 4 && xhr.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
        var data = xhr.responseText;
        var tmp = JSON.parse(data)
        console.log(tmp.result);		          // 取得した ファイルの中身を表示
        document.getElementById(outputid).innerHTML = JSON.stringify(tmp.result)
    }
  };
  var url = "/v1/search/" + "booknames" + "/"
  xhr.open('LIST',url+keyword);
  xhr.setRequestHeader('Authorization', 'bearer '+TOKEN);
  xhr.send(null);

}

function jsonDataOutput(jsondata) {
  var output = "";
  var table_title = table_list[selectdata] 
  var tableHeader = "<tr>";
  for (var i=0;i<table_title.length;i++){
      output += "<th>"+table_title[i]+"</th>"
  }
  tableHeader += "</tr>";
  var tablebody = "";
  for (var i=0;i<jsondata.length;i++){
    tablebody +="<tr>"
    for (var j=0;j<table_title.length;j++){
      if (table_title[j] == "id"){
        var id = jsondata[i][table_title[j]];
      }
      tablebody += "<td>"+jsondata[i][table_title[j]]+"</td>";
    }
    tablebody += "<td>"+edithtml(id)+"</td>"
    tablebody += "<td>"+deletehtml(id)+"</td>"
    tablebody +="</tr>"
  }
  output = "<table>"+tableHeader+tablebody+"</table>";
  return output;
}

function edithtml(id){
  var output = "";
  output += "<a href=javascript:void(0); onclick='vieweditform("+id+");return false'>"
  output += "edit</a>"
  return output
}

function deletehtml(id){
  var output = "";
  output += "<a href=javascript:void(0); onclick='destory("+id+");reviewform();return false'>"
  output += "destory</a>"
  return output
}

function jsonOutput2(str){
  var output = ""
  var table_title = ["id","name","title","writer","brand","booktype","ext"]
  if (meta_suburl=="filelist"){
    table_title = ["id","name","pdfpass","zippass","tag"]
  }else if(meta_suburl=="copyfile"){
    table_title = ["id","zippass","copyflag","filesize"]
  }
  var tmp = JSON.parse(str)
  output += "<table>"
  output += "<tr>"
  for (var i=0;i<table_title.length;i++){
      output += "<th>"+table_title[i]+"</th>"
  }
  output += "</tr>"
  for (var i=0; i< tmp.length;i++){
  //   output += "<div>"
    output += "<tr>"
  //   output += tmp[i].name
    output += "<td>"+tmp[i].id+"</td>"
    if (meta_suburl !="copyfile"){
      output += "<td>"+tmp[i].name+"</td>"
    }
    if (meta_suburl=="filelist"){
      output += "<td>"+tmp[i].pdfpass+"</td>"
      output += "<td>"+tmp[i].Zippass+"</td>"
      output += "<td>"+tmp[i].tag+"</td>"
    }else if(meta_suburl=="copyfile"){
      output += "<td>"+tmp[i].Zippass+"</td>"
      output += "<td>"+tmp[i].copyflag+"</td>"
      output += "<td>"+tmp[i].Filesize+"</td>"
    }else{
      output += "<td>"+tmp[i].title+"</td>"
      output += "<td>"+tmp[i].Writer+"</td>"
      output += "<td>"+tmp[i].brand+"</td>"
      output += "<td>"+tmp[i].booktype+"</td>"
      output += "<td>"+tmp[i].ext+"</td>"
    }
  //   output += " <a href='edit/"+tmp[i].id+"'>"+"edit"+"</a>"
    output += "<td><a href='show/"
    if (meta_suburl!=""){
      output +=meta_suburl+"/"
    }
    output += tmp[i].id+"'>"+"show"+"</a></td>"
    output += "<td><a href='edit/"
    if (meta_suburl!=""){
      output +=meta_suburl+"/"
    }
    output += tmp[i].id+"'>"+"edit"+"</a></td>"
  //   output += " <a href='destory/"+tmp[i].id+"'>"+"destory"+"</a>"
    output += "<td><a href='javascript:destory("+tmp[i].id+");'>"+"destory"+"</a></td>"
    if (meta_suburl=="copyfile"){
      output += "<td>"+"<input type='checkbox' "
      if (tmp[i].copyflag == "1"){
        output += "checked='checked' "
      }
      output += "onclick=\"ck_copyfilebox(\'"+tmp[i].Zippass+"\',this.checked)\""
      output += ">"+"</td>"
    }
    output += "</tr>"
  //   output +="</div>"
  }
  output += "</table>"
  return output
}

function getJSON(url,output) {
  var req = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
  req.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
    if(req.readyState == 4 && req.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
        var data = req.responseText;
        console.log(data);		          // 取得した JSON ファイルの中身を表示
        if (meta_suburl == `listdata`){
          document.getElementById(output).innerHTML = listoutput(data);
        }else{
          document.getElementById(output).innerHTML = jsonOutput(data);
        }
    }
  };
  req.open("GET", url+"/"+meta_suburl, false); // HTTPメソッドとアクセスするサーバーの　URL　を指定
  req.send(null);					    // 実際にサーバーへリクエストを送信
}

function serchDataTagSplit(tag){
  var output = ""
  var tmp = tag.split(",")
  for(var i=0;i<tmp.length;i++){
    //updataserch
    output += "<a href='"+"javascript:void(0);"+"'"
    output += " onclick="+"\"updataserch('"+tmp[i]+"');\""
    output += ">" +tmp[i]+ "</a>"
    if (i==0){
      output += "<br>\n"
    }else{
      output += " "
    }
  }
  return output
}


function searchgetData(output){
  var out = document.getElementById(output)
  var url
  url = "/v1/search/" + searchurl[selectdata] +"/"
  
  var keyword = document.getElementById("serch").value
  if (keyword == "") {
    out.innerHTML = "";
    return;
  }
  var req = new XMLHttpRequest();
  req.onreadystatechange = function(){
      if(req.readyState == 4 && req.status == 200){
          nowserchpage = 1
          var data=req.responseText;
          out.innerHTML = outputhtmlJson(data)
          //out.innerHTML = data
      }else if (req.readyState == 4 && req.status != 200){
        nowserchpage = 1
        var data=req.responseText;
        out.innerHTML = outputhtmlJson(data)
      }
  };
  req.open("GET",url+keyword,false);
  req.send(null);
}

function outputhtmlJson(str){
  var tmp = JSON.parse(str);
  var ary;
  console.log(tmp.result)
  if (tmp.result == "[]") {
    tmp.result = ""
  }
  ary = readlistedit(selectdata,tmp.result);
  // console.log(ary)
  return outputTable(ary)
}

function readlistedit(selectcout,ary) {
  var table_listdata = table_list[selectcout] ;
  var output = []
  for (var i=0;i<ary.length;i++){
      var tmp =[]
      for (var j=0;j<table_listdata.length;j++){
        tmp.push(ary[i][table_listdata[j]])
      }
      output.push(tmp)
  }
  return output
}