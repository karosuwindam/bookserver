var HOSTURL = ""
var COUNTMAX = 0    //表示ファイルの最大数
var FILELIST = {}   //表示ファイルリスト
var NowPage     //現在の表示ページ
var modeonetwe = true;      //2 page view

function onTitleData(id) {

    var req = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    req.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
      if(req.readyState == 4 && req.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
          var data = req.responseText;
          var jata = JSON.parse(data);
          console.log(jata);
          applyTileData(jata.Result[0])
          createListGet("listl", jata.Result[0].Tag)
      }else if (req.readyState == 4 && req.status != 200){ 
          var data = req.responseText;
          var jata = JSON.parse(data);
          console.log(jata);		          // 取得した JSON ファイルの中身を表示
      }
    };
    var url = HOSTURL + "/v1/read/filelists/" + id;
    req.open("GET", url, true); // HTTPメソッドとアクセスするサーバーの　URL　を指定
    req.send(null);					    // 実際にサーバーへリクエストを送信
}

function applyTileData(data) {
    window.document.title = data.Name
}

function onZipList(id, page) {
    var req = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    req.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
      if(req.readyState == 4 && req.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
          var data = req.responseText;
          var jata = JSON.parse(data);
          console.log(jata);
          ViewZipFIle(id, jata.Result, page)
      }else if (req.readyState == 4 && req.status != 200){ 
          var data = req.responseText;
          var jata = JSON.parse(data);
          console.log(jata);		          // 取得した JSON ファイルの中身を表示
      }
    };
    var url = HOSTURL + "/v1/view/" + id;
    req.open("GET", url, true); // HTTPメソッドとアクセスするサーバーの　URL　を指定
    req.send(null);					    // 実際にサーバーへリクエストを送信
}

function ViewZipFIle(id, data, page) {
    FILELIST = jpgCount(data.Name)
    COUNTMAX = FILELIST.length

    var output = ""
    if (page < 1) {
        page = 1
    }else if (page > COUNTMAX+1){
        page = COUNTMAX+1
    }
    NowPage = page
    output += createViewPage(id,page)
    console.log()
    document.getElementById("output").innerHTML = output
    imageload(page)
}

function jpgCount(data) {
    var tmp = []
    for (var i=0;i<data.length;i++) {
        if (data[i].toLowerCase().indexOf(".jpg")>0){
            tmp.push(data[i])
        }
    }
    return tmp
}

function createViewPage(id, page) {
    var output = ""
    var url0, url1
    for (var i=0;i<=COUNTMAX;i++) {
        if (i<=0) {
            url0 = "img/blank.jpg"
        }else {
            url0 = "/v1/image/"+id+"/"+FILELIST[i-1]
        }
        if (i>=COUNTMAX) {
            url1 = "img/blank.jpg"
        }else {
            url1 = "/v1/image/"+id+"/"+FILELIST[i]
        }
        if (page -1 == i) {
            output += "<div class=\"pageI0\">" + "<img class=\"page0\" data-src=\""+url0+"\" src=\"\">" +"</div>"
            output += "<div class=\"pageI1\">" + "<img class=\"page1\" data-src=\""+url1+"\" src=\"\">" +"</div>"
        }else {
            output += "<div class=\"pageI0\">" + "<img class=\"page0\" style=\"display: none;\" data-src=\""+url0+"\" src=\"\">" +"</div>"
            output += "<div class=\"pageI1\">" + "<img class=\"page1\" style=\"display: none;\" data-src=\""+url1+"\" src=\"\">" +"</div>"
        }
    }
    document.getElementById("pageslider").value = page
    document.getElementById("pageslider").max = COUNTMAX
    nowPageView()
    return output

}

//ページ移動設定
function viewPageChange(count) {
    var maxup = 0
    if (modeonetwe) {
        maxup = 1
    }
    if ((NowPage == COUNTMAX+maxup)&&(count>0)) {
        NowPage = 1
    }else if((NowPage == 1)&&(count<0)){
        NowPage = COUNTMAX+maxup
    }else {
        if (!modeonetwe){
            if (count>0) {
                NowPage++
            }else{
                NowPage--
            }
        }else{
            NowPage += count
        }
    }
    if (NowPage > COUNTMAX+maxup) {
        NowPage = COUNTMAX+maxup
    }else if (NowPage <= 0) {
        NowPage = 1
    }
    nowPage(NowPage)
}

function nowPage(nowpage){
    var img0_elements = document.getElementsByClassName("page0")
    var img1_elements = document.getElementsByClassName("page1")
    for (var i=0;i<img0_elements.length;i++) {
        if (i==(nowpage-1)) {
            img0_elements[i].style.display = ""
        }else {
            img0_elements[i].style.display = "none"
        }
    }
    for (var i=0;i<img1_elements.length;i++) {
        if (i==(nowpage-1)) {
            img1_elements[i].style.display = ""
        }else {
            img1_elements[i].style.display = "none"
        }
    }
    NowPage = nowpage-0;

    document.getElementById("pageslider").value=NowPage;
    nowPageView()

}

function imageload(page) {
    let img_elements = document.querySelectorAll("img");

    if (page > 1){

        img_elements[page-1].addEventListener('load', (e)=> {
            console.log(e.target.alt + " load");
        });
        img_elements[page-1].src = img_elements[page-1].getAttribute("data-src");
        if (page == COUNTMAX) {
            img_elements[0].addEventListener('load', (e)=> {
                console.log(e.target.alt + " load");
            });
            img_elements[0].src = img_elements[0].getAttribute("data-src");
        }else {

        img_elements[page].addEventListener('load', (e)=> {
            console.log(e.target.alt + " load");
        });
        img_elements[page].src = img_elements[page].getAttribute("data-src");
        }
    }
    for(let i=0; i<img_elements.length; i++) {

        // 画像読み込み完了したときの処理
        img_elements[i].addEventListener('load', (e)=> {
            console.log(e.target.alt + " load");
        });

        // 遅延読み込み
        img_elements[i].src = img_elements[i].getAttribute("data-src");
    }
  }

//表示ページを2ページと1ページを切り替える
function chPageOneTwe(mode) {
    var l_imgdiv = document.getElementsByClassName("pageI0")
    var r_imgdiv = document.getElementsByClassName("pageI1")
    switch (mode) {
        case 0:
            modeonetwe = true
            break;
        case 1:
            modeonetwe = false
            break;
        case 2:
            modeonetwe = !modeonetwe
            break;
        default:
            break;
    }
    if (modeonetwe){    //2page表示
        for (var i=0;i<r_imgdiv.length;i++){
            r_imgdiv[i].style.textAlign = ""
            r_imgdiv[i].style.width = ""
        }
        for (var i=0;i<l_imgdiv.length;i++){
            l_imgdiv[i].style.display = ""
        }
    }else { //1page 表示
        for (var i=0;i<r_imgdiv.length;i++){
            r_imgdiv[i].style.textAlign = "center"
            r_imgdiv[i].style.width = "100%"
        }
        for (var i=0;i<l_imgdiv.length;i++){
            l_imgdiv[i].style.display = "none"
        }
    }
    nowPageView()
}

function nowPageView() {
    var maxup = 0
    if (modeonetwe) {
        maxup = 1
    }
    document.getElementById("nowpage").innerHTML = (NowPage)+"/"+(COUNTMAX+maxup)+" page"
    document.getElementById("pageslider").max = COUNTMAX+maxup

}


function listonoff(){
    var list_div = document.getElementById("listl");
    list_flag = !list_flag;
    if(list_flag){
        list_div.style.display = "";
        document.getElementById("maxmin").style.display = "none";
    }else{
        list_div.style.display = "none";
        document.getElementById("maxmin").style.display = "";
    }
}

function createListGet(output,tag) {
    TMP_tag = tag                       //タグ名を一時保存
    var req = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    req.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
      if(req.readyState == 4 && req.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
          var data = req.responseText;
          var jata = JSON.parse(data);
          console.log(jata);
          document.getElementById(output).innerHTML = createListData(output,jata.Result,0)
      }else if (req.readyState == 4 && req.status != 200){ 
          var data = req.responseText;
          var jata = JSON.parse(data);
          console.log(jata);		          // 取得した JSON ファイルの中身を表示
      }
    };
    var jsondata = {};
    var tmp = tag.split(",")
    jsondata["Table"] = "filelists";
    var keyword = tmp[0]

    for ( var i=1;i<tmp[0].length;i++){
        keyword = tmp[0].slice(-i) -0;
        if (isNaN(keyword)){
            keyword = tmp[0].slice(0,tmp[0].length-i+1)
            break
        }
    }
    jsondata["Keyword"] = keyword;
    var url = HOSTURL + "/v1/search";
    req.open("POST", url, true); // HTTPメソッドとアクセスするサーバーの　URL　を指定
    req.send(JSON.stringify(jsondata));					    // 実際にサーバーへリクエストを送信
}

var TMP_tag
var TmpListData
var upflag = false

function createTmpListGet(outid,num) {
    var req = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    req.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
      if(req.readyState == 4 && req.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
          var data = req.responseText;
          var jata = JSON.parse(data);
          console.log(jata);
          TmpListData = jata.Result
          document.getElementById(outid).innerHTML = createListData(outid,jata.Result,num)
      }else if (req.readyState == 4 && req.status != 200){ 
          var data = req.responseText;
          var jata = JSON.parse(data);
          console.log(jata);		          // 取得した JSON ファイルの中身を表示
      }
    };
    var jsondata = {};
    var tmp = TMP_tag.split(",")
    jsondata["Table"] = "filelists";
    var keyword = tmp[num]

    for ( var i=1;i<tmp[num].length;i++){
        keyword = tmp[num].slice(-i) -0;
        if (isNaN(keyword)){
            keyword = tmp[num].slice(0,tmp[num].length-i+1)
            break
        }
    }
    jsondata["Keyword"] = keyword;
    var url = HOSTURL + "/v1/search";
    req.open("POST", url, true); // HTTPメソッドとアクセスするサーバーの　URL　を指定
    req.send(JSON.stringify(jsondata));					    // 実際にサーバーへリクエストを送信
}

function createListData(outid,data,num) {
    var output = ""
    output += "<div>"
    var tmp = TMP_tag.split(",")
    for (var i=0;i<tmp.length;i++) {
        if (i==num){
            output += "<a href=javascript:void(0); class=\"tab-button-a\" onclick='createTmpListGet(\""+outid+"\","+i+")'>"+tmp[i]+"</a>"
        }else{
            output += "<a href=javascript:void(0); class=\"tab-button\" onclick='createTmpListGet(\""+outid+"\","+i+")'>"+tmp[i]+"</a>"
        }
    }
    output += "<a href=javascript:void(0); class=\"tab-button\" onclick='sortListData(\""+outid+"\","+num+")'>"+"sort"+"</a>"
    output += "</div>"
    for (var i=0;i<data.length;i++) {
        var url = "/view/" + data[i].Id
        output += "<div class=\"list\">"+"<a href=\""+url+"\">"+data[i].Zippass+"</a></div>"
    }
    return output
}

function sortListData(outid,num) {
    var tmp = TmpListData
    upflag = !upflag
    tmp.sort((a, b) => {
      if (upflag){
        return a.Zippass < b.Zippass ? -1 : 1;
      }else {
        return a.Zippass > b.Zippass ? -1 : 1;
      }
    });
    document.getElementById(outid).innerHTML = createListData(outid,tmp,num)
}