var HOSTURL = "";        //検索先のURLについて
var SEARCHTABLE = "filelists"; 
var TmpJdata;
var upflag = true
var MAXFILESIZE = "<%uploadsize%>" // 256*1024*1024 //Upload Size : <%uploadsize%>


function formdataJSON(inputElement){
  var filelist = inputElement.files;
  document.getElementById("fileck").innerHTML = ""
  var filelists = document.getElementById("file").files
  if (filelists.length > 0) {
    var countsize = 0
    for(var i=0;i<filelists.length;i++){
      countsize += filelists[i].size
    }
    
    if (countsize/(1<<30) > 1) {
      document.getElementById("fileck").innerHTML = "Sum size :" + Math.round(countsize/(1<<30)*10)/10 + "G"
    }else if (countsize/(1<<20) > 1) {
      document.getElementById("fileck").innerHTML = "Sum size :" + Math.round(countsize/(1<<20)*10)/10 + "M"
    }else if (countsize/(1<<10) > 1) {
      document.getElementById("fileck").innerHTML = "Sum size :" + Math.round(countsize/(1<<10)*10)/10 + "K"
    }else {
      document.getElementById("fileck").innerHTML = "Sum size :" + countsize
    }
  }
  for(var i=0;i<filelist.length;i++){
    getformdataJSON(filelist[i].name)
  }
}

function getformdataJSON(filename) {
  var req = new XMLHttpRequest();
  req.onreadystatechange = function(){
    if(req.readyState == 4 && req.status == 200){
      var data=req.responseText;
      var tmp = JSON.parse(data)
      console.log(tmp)
      if (document.getElementById("fileck").innerHTML==""){
        document.getElementById("fileck").innerHTML = fileckdata(tmp.Result)
      }else {
        document.getElementById("fileck").innerHTML += "<br>"+fileckdata(tmp.Result)
      }
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
  if (str.Name.toLowerCase().indexOf('.pdf')>0) {       //.pdfは含まない
      output += " create file: " + str.ChangeName.Zip
  }else if (str.Name.toLowerCase().indexOf('.zip')>0) { //.zipは含まない
      output += " create file: " + str.ChangeName.Pdf
  }
  if (str.Overwrite) {
    output += " テーブル上書きあり"
  }
  return output
}

var postFileFlag = false

function sizeread() {
  var output = 0
  var tmp = MAXFILESIZE.toLowerCase()
  var i =0  
  if ((i=tmp.indexOf("t"))>0) { //tだけは含まない
    output = (tmp.substr(0,i) -0)*(1<<40);
  }else if ((i=tmp.indexOf("g"))>0) {
    output = (tmp.substr(0,i) -0)*(1<<30);
  }else if ((i=tmp.indexOf("m"))>0) {
    output = (tmp.substr(0,i) -0)*(1<<20);
  }else if ((i=tmp.indexOf("k"))>0) {
    output = (tmp.substr(0,i) -0)*(1<<10);
  }else {
    output = (tmp -0);
  }
  return output
}

async function postFile() {
  if (document.getElementById("file").files.length == 0){
      return 
  }
  document.getElementById("file").disabled = true;
  document.getElementById("post2").disabled = true;
  var formData = new FormData();
  var fileSizeSum = 0
  postFileFlag = true
  for(var i=0;i<document.getElementById("file").files.length;i++){
      fileSizeSum += document.getElementById("file").files[i].size
      if ((fileSizeSum > sizeread())&&(i!=0)) {
        await PostFIleUpload(formData)
        fileSizeSum = document.getElementById("file").files[i].size
        formData = new FormData()
      }
      formData.append("file", document.getElementById("file").files[i]);
  }
  postFileFlag = false
  await PostFIleUpload(formData)
}

function PostFIleUpload(formData) {
  const p = new Promise((resolve, reject) => {
    var url = HOSTURL + "/v1/upload"

    var request = new XMLHttpRequest();
    request.upload.addEventListener("progress", updateProgress, false);
    request.open("POST", url);
    request.send(formData);
  });

  return p;

}

function updateProgress(e) {
  if (e.lengthComputable) {
      var percent = e.loaded / e.total;
      document.getElementById("progress").value = percent * 100;
      if (percent == 1){
        document.getElementById("progress").value = 0;
          if (!postFileFlag){
            document.getElementById("file").disabled = false;
            document.getElementById("post2").disabled = false;
            document.getElementById("file").value = "";
            document.getElementById("fileck").innerHTML = "";
            document.getElementById("health").innerHTML = "";
          }
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
        TmpJdata=jata.Result
        document.getElementById(output).innerHTML = viewSearchTable(jata.Result, output)
        imageload();
        ckboxupdate();
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

function outputSortData(outid) {
  var tmp = TmpJdata
  if ((tmp == undefined||tmp == "")) {
    return
  }
  upflag = !upflag
  tmp.sort((a, b) => {
    if (upflag){
      return a.Zippass < b.Zippass ? -1 : 1;
    }else {
      return a.Zippass > b.Zippass ? -1 : 1;
    }
  });

  document.getElementById(outid).innerHTML = viewSearchTable(tmp, outid)
  imageload();
  ckboxupdate();
}

function viewSearchTable(jdata, id) {
  var output = ""
  var count = document.getElementById(id).clientWidth / 280
  if (count < 0) {
    count = 1
  }else {
    count = Math.floor(count)
  }
  tmpjdata = jdata
  for (var i=0;i<tmpjdata.length;i++){
    output += createViewCell(i)
    if (i%count==(count-1)){
      output += "<br>"
    }
  }
  return output
}


function serchDataTagSplit(tag){
  var output = ""
  var tmp = tag.split(",")
  for(var i=0;i<tmp.length;i++){
    if (tmp[i] == "") {
      continue;
    }
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
  var tmp = str
  for ( var i=1;i<str.length;i++){
      tmp = str.slice(-i) -0;
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
  var pdf_url = HOSTURL + "/v1/download/pdf/" + jtmp.Id
  var zip_url = HOSTURL + "/v1/download/zip/" + jtmp.Id
  var view_url = HOSTURL + "/view/" + jtmp.Id
  output += "<div class=\"serchdata\">"
  output += "<div>"+"<a href=\""+view_url+"\" target=\"_blank\">"+"<img class='cell' data-src=\"img/"+jtmp.Name+".jpg\" src=\"img/"+jtmp.Name+".jpg\">"+"</a>"+"</div>"
  output += "<div>" + serchDataTagSplit(jtmp.Tag)
  output +="</div>" + "<a class=\"button\" href=\""+pdf_url+"\">pdf download</a>" + "<a class=\"button\" href=\""+zip_url+"\">zip download</a>"
  output += "<div class=\"copyckbox\">" + jtmp.Id
  output +="</div>"
  output +="</div>"
  return output
}

function imageload() {
  let img_elements = document.querySelectorAll("img");
	for(let i=0; i<img_elements.length; i++) {

		// 画像読み込み完了したときの処理
		img_elements[i].addEventListener('load', (e)=> {
			console.log(e.target.alt + " load");
		});

		// 遅延読み込み
		img_elements[i].src = img_elements[i].getAttribute("data-src");
	}
}