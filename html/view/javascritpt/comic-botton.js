var list_flag = false;

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
var maxmin = false;
function maxminonoff(){
    if( ! enabledFullScreen() ){
        alert("フルスクリーンに対応していません");
        return(false);
    }
    if (!maxmin){
        goFullScreen();
        document.getElementById("maxmin").innerHTML = "Min";
    }else{
        cancelFullScreen();
        document.getElementById("maxmin").innerHTML = "Max";
    }
    maxmin = !maxmin;
}
/**
 * フルスクリーンが利用できるか
 *
 * @return {boolean}
 */
function enabledFullScreen(){
    return(
      document.fullscreenEnabled || document.mozFullScreenEnabled || document.documentElement.webkitRequestFullScreen || document.msFullscreenEnabled
    );
  }
  
/**
 * フルスクリーンにする
 *
 * @param {object} [element]
 */
function goFullScreen(element=null){
    const doc = window.document;
    const docEl = (element === null)?  doc.documentElement:element;
    let requestFullScreen = docEl.requestFullscreen || docEl.mozRequestFullScreen || docEl.webkitRequestFullScreen || docEl.msRequestFullscreen;
    requestFullScreen.call(docEl);
  }
  
  /**
   * フルスクリーンをやめる
   */
  function cancelFullScreen(){
    const doc = window.document;
    const cancelFullScreen = doc.exitFullscreen || doc.mozCancelFullScreen || doc.webkitExitFullscreen || doc.msExitFullscreen;
    cancelFullScreen.call(doc);
  }
function jsonOutput(str){
    var output = ""
    var tmp = JSON.parse(str)
    for (var i=0; i< tmp.length;i++){
        var data = tmp[i].tag.split(",")
      output += "<div class=\"list\">"+"<a href='/view/"+tmp[i].id+"'>"+data[0]+"</a></div>"
    }
    return output
}
function getlist(str) {
    var req = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    req.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
      if(req.readyState == 4 && req.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
          var data = req.responseText;
          console.log(data);		          // 取得した JSON ファイルの中身を表示
          document.getElementById('jsonout').innerHTML = jsonOutput(data);
      }
    };
    req.open("GET", "/serch/filelist/"+str, false); // HTTPメソッドとアクセスするサーバーの　URL　を指定
    req.send(null);					    // 実際にサーバーへリクエストを送信
  }
function listdatawrite(){
    var listout = document.getElementById("listl");
    var str = ""
    str += "<div class=\"list\">" +"<a href=\"/\">index</a>" +"</div>";
    str += "<div id=\"jsonout\">"
    for(var i=0;i<20;i++){
        str += "<div class=\"list\">"+ createlink("#","data"+i) + "</div>"
    }
    str += "</div>"
    str += "<div class=\"list\">"+ "<a href=javascript:void(0); onclick=\"helpviewonoff();return false\">help</a><br>"+ "</div>"
    str += "<div class=\"list\">"+ "<a href=javascript:void(0); onclick=\"pageviewonoff();return false\">page</a><br>"+ "</div>"
    str += "<div class=\"list\">"+ "<a href=javascript:void(0); onclick=\"infoviewonoff();return false\">info</a><br>"+ "</div>"
    str += "<div class=\"list\">"+ "<a href=javascript:void(0); onclick=\"borderviewonoff();return false\">border</a><br>"+ "</div>"

    listout.innerHTML = str;
}

function createlink(link,name){
    var str = ""
    str += "<a class=\"list\" href=\"" +link + "\">"
    str += name
    str += "</a>"
    return str
}