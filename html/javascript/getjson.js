var meta_suburl=""
var jsondata
var rowmax = 8
var rownum = 4
var nowserchpage = 1
function destory(id){
  myRet = confirm("destory id="+id+" OK??");
  if (myRet){
    var xhr = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    xhr.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
      if(xhr.readyState == 4 && xhr.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
          var data = xhr.responseText;
          console.log(data);		          // 取得した ファイルの中身を表示
          document.getElementById("answer").innerHTML = "destory id=" + id + " OK"
          getJSON("/list","output");
      }
    };

    var url = "destory/"
    if (meta_suburl != ""){
      url += meta_suburl + "/"
    }
    url += id;
    xhr.open('POST', url, true);
    xhr.setRequestHeader('content-type', 'application/x-www-form-urlencoded;charset=UTF-8');
    xhr.send(null);
  }
}
function addfile(name,zippass,pdfpass,tag){
  var xhr = new XMLHttpRequest();
  xhr.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
    if(xhr.readyState == 4 && xhr.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
        var data = xhr.responseText;
        console.log(data);		          // 取得した ファイルの中身を表示
        document.getElementById("answer").innerHTML = "destory id=" + id + " OK"
    }
  };
  var url = "/new/filelist"
  xhr.open('POST',url,true)
  xhr.setRequestHeader('content-type', 'application/x-www-form-urlencoded;charset=UTF-8');
  var output = "name="+name+"&zippass="+zippass+"&pdfpass="+pdfpass+"&tag="+tag
  xhr.send(output);
}
function ck_copyfilebox(str,ckflag){
  var xhr = new XMLHttpRequest();
 
  xhr.open('POST', 'ckbox');
  xhr.setRequestHeader('content-type', 'application/x-www-form-urlencoded;charset=UTF-8');
  var output = "zippass=" + str + "&copyflag="+ckflag;
  xhr.send( output );
  //alert(output);

}

function serchBoxgetdata(str){
  var tmp = JSON.parse(str);
  for(var i=0;i<tmp.length;i++){
    ck_copyfilebox_ckj(tmp[i].Zippass,"data"+i)
  }
}

function ck_copyfilebox_ckj(str,ckdata){
  var xhr = new XMLHttpRequest();
  var URL = "/ckbox?" +"zippass=" + str;
  xhr.open('GET',URL,true);
  xhr.send( null );
  xhr.onreadystatechange = function(){
    if(xhr.readyState == 4){
      if(xhr.status == 200){
        var flag = false
        if (xhr.responseText=="1"){
          flag = true;
        }
        //<!-- レスポンスが返ってきたらテキストエリアに代入する -->
        document.getElementsByName(ckdata)[0].checked = flag;
      }
    }
  }
}

function ck_copyfilebox_ck(str,ckdata){
  var xhr = new XMLHttpRequest();
  var URL = "/ckbox?" +"zippass=" + str;
  var tmp = ckdata
  xhr.open('GET',URL,true);
  xhr.send( null );
  xhr.onreadystatechange = function(){
    if(xhr.readyState == 4){
      if(xhr.status == 200){
        var flag = false
        if (xhr.responseText=="1"){
          flag = true;
        }
        //<!-- レスポンスが返ってきたらテキストエリアに代入する -->
        document.getElementsByName(ckdata.name)[0].checked = flag;
      }
    }
  }
}

function listoutput(str){
  var output = ""
  table_title = ["pdfname","flag","zipname","flag","jpgname","flag"]
  var ary = JSON.parse(str)
  output += "<div>Time:"+ary.Time+"s</div><br>"
  var tmp = ary.Data
  output +="<table>"
  output += "<tr>"
  for (var i=0;i<table_title.length;i++){
      output += "<th>"+table_title[i]+"</th>"
  }
  output += "</tr>"
  for (var i=0;i < tmp.length;i++){
    output += "<tr>"
    output += "<td>"+tmp[i].Pdf.Name
    if (tmp[i].Pdf.Flag=="0"){
      output += " "+"<a href='/new/bookname'>"+"New"+"</a>"
    }else{
      
    }
    output +="</td>"
    output += "<td>"+tmp[i].Pdf.Flag+"</td>"
    output += "<td>"+tmp[i].Zip.Name
    if ((tmp[i].Pdf.Flag=="1")&&(tmp[i].Zip.Flag=="0")&&(tmp[i].Jpg.Flag=="1")&&(tmp[i].Data.name!="")){
      output += " " 
      // output += "<form action='/new/filelist' method='post'>"
      // output += "<input type='hidden' name='name' value="+tmp[i].Data.name+">"
      // output += "<input type='hidden' name='pdfpass' value="+tmp[i].Data.pdfpass+">"
      // output += "<input type='hidden' name='zippass' value="+tmp[i].Data.Zippass+">"
      // output += "<input type='hidden' name='tag' value="+tmp[i].Data.tag+">"
      // output += "<input type='submit' value='send'>"
      // output += "</form>"
      output += "<input type='button' value='send' onclick=\""
      output += "addfile('"+tmp[i].Data.name+"','"+tmp[i].Data.Zippass+"','"+tmp[i].Data.pdfpass+"','"+tmp[i].Data.tag+"')"
      output += ";this.disabled=true;return false\"> none"
    }else if (tmp[i].Zip.Flag=="0"){
      output += " none"
    }
    output +="</td>"
    output += "<td>"+tmp[i].Zip.Flag+"</td>"
    output += "<td>"+tmp[i].Jpg.Name+"</td>"
    output += "<td>"+tmp[i].Jpg.Flag+"</td>"
    output += "</tr>"
  }
  output += "</table>"
  return output
}

function jsonOutput(str){
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
// nowserchpage
function serchpageout(tmp){
  var num = tmp.length
  var output=""
  for(var i=0;i<num/rowmax;i++){
    if (i>0){
      if (nowserchpage == (i+1)){
        output += " "+(i+1);
      }else{
        output += " "+"<a href=\"javascript:void(0);\" onclick='"+"chData("+(i+1)+");"+"'>"+(i+1)+"</a>";
      }
    }else{
      if (nowserchpage == (i+1)){
        output += (i+1);
      }else{
        output += "<a href='#' onclick='"+"chData("+(i+1)+");"+"'>"+(i+1)+"</a>";        
      }
    }
  }
  document.getElementById("page").innerHTML = output;
}
function outputSerchData(tmp,num){
  var output = ""
  for(var i=rowmax*(num-1);(i<tmp.length);i++){
    output += "<div class=\"serchdata\">"
    output += "<a href='"+"/view/"+tmp[i].id+"' target=\"_blank\">"
    output += "<img width='250px' src='jpg/"+tmp[i].name+".jpg"+"' title='"+tmp[i].tag+"'>"
    output +="</a><br>\n"
    output += serchDataTagSplit(tmp[i].tag)
    output += "<br>"
    output += "<a href='"+"/download/zip/"+tmp[i].id+"'>"+ "zip download" +"</a>"
    output += " <a href='"+"/download/pdf/"+tmp[i].id+"'>"+ "pdf download" +"</a>"
    output += "<input type='checkbox' "
    output += "onclick=\"ck_copyfilebox(\'"+tmp[i].Zippass+"\',this.checked)\" "
    output += "name=\"data"+i+"\""
    output += " id=\"ckbox"+i+"\""
    output += ">"
    //ck_copyfilebox_ck()
    output += "<input type='button' "
    output += "onclick=\""+"ck_copyfilebox_ck('"+tmp[i].Zippass+"',this)"+""+"\" "
    output += "name=\"data"+i+"\""
    output += ">"
    output += "</div>\n"
    if (i%rownum==(rownum-1)){
    output += "<br>"}
    if (i>=rowmax*(num)-1){
      break;
    }
  }
  return output;
}
function serchDataGet(str){
  var tmp = JSON.parse(str);
  jsondata = tmp;
  serchpageout(jsondata)
  return outputSerchData(jsondata,nowserchpage);
}


function chData(num){
  nowserchpage = num -0
  serchpageout(jsondata);
  var tmp = outputSerchData(jsondata,nowserchpage);
  document.getElementById("output").innerHTML = tmp;
}
function serchgetJSON(output){
  var keyword = document.getElementById("keyword").value
  var req = new XMLHttpRequest();
  req.onreadystatechange = function(){
    if(req.readyState == 4 && req.status == 200){
      nowserchpage = 1
      var data=req.responseText;
      document.getElementById(output).innerHTML = serchDataGet(data);
      serchBoxgetdata(data);
    }
  };
  req.open("GET","/serch/filelist/"+keyword,false);
  req.send(null);
}
function fileckdata(str){
  var output = "not file"
  var tmp = JSON.parse(str)
  if (str != "{}"){
    output = tmp[0].title
    if (tmp[1].flag == 1){
      output += " 既存ファイルあり"
    }else{
      output += " file is not"
    }
  }
  return output
}
function formdataJSON(inputElement){
  var filelist = inputElement.files;
  var filename = filelist[0].name
  tmp = filename.substr(0,filename.length-4)
  var req = new XMLHttpRequest();
  req.onreadystatechange = function(){
    if(req.readyState == 4 && req.status == 200){
      var data=req.responseText;
      document.getElementById("fileck").innerHTML = fileckdata(data);
    }
  };
  req.open("GET","/mach/"+tmp,false);
  req.send(null)
}