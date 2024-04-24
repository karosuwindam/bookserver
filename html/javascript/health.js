
var HEATHURL = "/health"

var statusTimeer = 0

function healthckeck(outid){
    clearInterval(statusTimeer)
    var req = new XMLHttpRequest();
    req.onreadystatechange = function(){
        if(req.readyState == 4 && req.status == 200){
        var data=req.responseText;
        var tmp = JSON.parse(data)
        console.log(tmp)
        document.getElementById(outid).innerHTML = healthcheckout(tmp)
        var clear = function(){
            document.getElementById(outid).innerHTML = ""
        };
        statusTimeer = setInterval(clear,1000)
        };
    }
    var url = HOSTURL + HEATHURL
    req.open("GET",url,true);
    req.send(null);  
}

function healthcheckout(json) {
    var output = "<ul>"
    var convert = json.Controller.Convert
    output += "<li>処理前</li>"
    output += "<ul>"
    for (var i=0;i<convert.Startfile.length;i++){
        output += "<li>"+convert.Startfile[i]+"</li>"
    }
    output += "</ul>"
    output += "<li>処理済み</li>"
    output += "<ul>"
    for (var i=0;i<convert.Endfile.length;i++){
        output += "<li>"+convert.Endfile[i]+"</li>"
    }
    output += "</ul>"
    output += "</ul>"
    return output
}

function isSmartPhone() {
      const ua = navigator.userAgent;
      var width = window.innerWidth
      var hight = window.innerHeight
    if (navigator.userAgent.match(/iPhone|Android.+Mobile/)) {
      return true;
      
    } else if ((width<hight))   {
        return true
    } else {
      return false;
    }
  }
  