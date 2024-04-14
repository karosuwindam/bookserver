
var HEATHURL = "/health"

function healthckeck(outid){
    var req = new XMLHttpRequest();
    req.onreadystatechange = function(){
        if(req.readyState == 4 && req.status == 200){
        var data=req.responseText;
        var tmp = JSON.parse(data)
        console.log(tmp)
        document.getElementById(outid).innerHTML = healthcheckout(tmp)
        };
    }
    var url = HOSTURL + HEATHURL
    req.open("GET",url,true);
    req.send(null);  
}

function healthcheckout(json) {
    var output = ""
    var convert = json.Controller.Convert
    if (convert.Status) {
        output += "<li>処理前</li>"
        output += "<li>"
        for (var i=0;i<convert.Startfile.length;i++){
            output += "<li>"+convert.Startfile[i]+"</li>"
        }
        output += "</li>"
        output += "<li>処理済み</li>"
        output += "<li>"
        for (var i=0;i<convert.Endfile.length;i++){
            output += "<li>"+convert.Endfile[i]+"</li>"
        }
        output += "<li>"
    }
    return output
}