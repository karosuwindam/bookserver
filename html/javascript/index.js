

var HOSTURL = "";        //検索先のURLについて


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
    // var jdata = {}
    // jdata["Name"] = filename
    // console.log(JSON.stringify(jdata));
    req.open("GET",url,true);
    // req.send(JSON.stringify(jdata));
    // req.send(JSON.stringify(jdata))
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
      // $("progress").attr("value", percent * 100);
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