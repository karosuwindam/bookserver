var meta_suburl=""
var jsondata
var rowmax = 8
var rownum = 4
var nowserchpage = 1

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