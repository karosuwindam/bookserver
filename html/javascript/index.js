

var HOSTURL = "";        //検索先のURLについて
var SEARCHTABLE = "filelists";

function formdataJSON(inputElement){
    var filelist = inputElement.files;
    var filename = filelist[0].name
    var req = new XMLHttpRequest();
    req.onreadystatechange = function(){
      if(req.readyState == 4 && req.status == 200){
        var data=req.responseText;
        var tmp = JSON.parse(data)
        console.log(tmp)
        document.getElementById("fileck").innerHTML = fileckdata(tmp.Result)
      }
    };
    var url = HOSTURL + "/v1/upload/" + filename
    req.open("GET",url,true);
    req.send(null);
  }
  
  function fileckdata(str){
    var output = "not file"
    if (str.Register) {
        output = str.Name + " 既存ファイルあり"
    }else {
        output = str.Name + " file is not"
    }
    if (str.Name.toLowerCase().indexOf('.pdf')>0) {
        output += " create file: " + str.ChangeName.Zip
    }else if (str.Name.toLowerCase().indexOf('.zip')>0) {
        output += " create file: " + str.ChangeName.Pdf
    }
    return output
  }


function postFile() {
  if (document.getElementById("file").files.length == 0){
      return 
  }
  document.getElementById("file").disabled = true;
  document.getElementById("post2").disabled = true;
  var formData = new FormData();
  formData.append("file", document.getElementById("file").files[0]);
  var url = HOSTURL + "/v1/upload"

  var request = new XMLHttpRequest();
  request.upload.addEventListener("progress", updateProgress, false);
  request.open("POST", url);
  request.send(formData);
}

function updateProgress(e) {
  if (e.lengthComputable) {
      var percent = e.loaded / e.total;
      document.getElementById("progress").value = percent * 100;
      if (percent == 1){
          document.getElementById("file").disabled = false;
          document.getElementById("post2").disabled = false;
          document.getElementById("progress").value = 0;
          document.getElementById("file").value = "";
          document.getElementById("fileck").innerHTML = "";
          document.getElementById("health").innerHTML = "";
      }
  }
}

function getSearchData(output) {
  var req = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
  req.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
    if(req.readyState == 4 && req.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
        var data = req.responseText;
        var jata = JSON.parse(data);
        console.log(jata);		          // 取得した JSON ファイルの中身を表示
        document.getElementById(output).innerHTML = viewSearchTable(jata.Result)
    }else if (req.readyState == 4 && req.status != 200){ 
        var data = req.responseText;
        var jata = JSON.parse(data);
        console.log(jata);		          // 取得した JSON ファイルの中身を表示
    }
  };
  var url = HOSTURL + "/v1/search"
  var jsondata = {};
  jsondata["Table"] = SEARCHTABLE;
  jsondata["Keyword"] = document.getElementById("keyword").value;
  req.open("POST", url, true); // HTTPメソッドとアクセスするサーバーの　URL　を指定
  req.send(JSON.stringify(jsondata));					    // 実際にサーバーへリクエストを送信
}

function viewSearchTable(jdata) {
  var output = ""
  tmpjdata = jdata
  for (var i=0;i<tmpjdata.length;i++){
    output += createViewCell(i)
  }
  return output
}


function serchDataTagSplit(tag){
  var output = ""
  var tmp = tag.split(",")
  for(var i=0;i<tmp.length;i++){
    //updataserch
    output += "<a class=\"button\" href='"+"javascript:void(0);"+"'"
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

function updataserch(str){
  for ( var i=1;i<str.length;i++){
      var tmp = str.slice(-i) -0;
      if (isNaN(tmp)){
          tmp = str.slice(0,str.length-i+1)
          break
      }
  }
  document.getElementById("keyword").value = tmp;
  getSearchData('output');
}

var tmpjdata
function createViewCell(count) {
  var jtmp = tmpjdata[count]
  var output = ""
  output += "<div class=\"serchdata\">"
  output += "<div>"+"<img class='cell' src=\"img/"+jtmp.Name+".jpg\">"+"</div>"
  output += "<div>" + serchDataTagSplit(jtmp.Tag)
  output +="</div>" + "pdf download" + "zip download"
  output += "<div>"
  output +="</div>" 
  output +="</div>"
  return output
}