
var viewEditFormF = false;
function viewEditForm(output) {
    if (viewEditFormF) {
        document.getElementById("edit").style.display = "";
        createAddedForm(output);
    }else {
        document.getElementById("edit").style.display = "none";
    }
}
function viewaddform(output) {
    viewEditFormF = !viewEditFormF
    viewEditForm(output)
}

function closeViewForm(output) {
    viewEditFormF = false
    viewEditForm(output)
}

function createAddedForm(output) {
    var data = ""
    var list = table_list[selectdata]
    var tablename = searchurl[selectdata]
    data += tablename + "<br>"
    data += "<table>"+"<tr>"+"<th>Key名</th>"+"<th>値</th>"+"<th>nowinput</th>"+"<th>button</th>"+"</tr>"
    for (var i=0;i<list.length;i++) {
        if (list[i] == "Id") {
            continue
        }
        data += "<tr>"
        data += "<td>" + list[i] + "</td>"
        data += "<td>"
        data += "<input type=\"text\" id=\""+tablename+"_"+list[i]+"\" onkeyup=\"serchckbox(this.value,'"+list[i]+"');return\">"
        data += "</td>"
        data += "<td>"+createlistbox(list[i])+"</td>"
        data += "<td>"+"<input type=\"button\" value=\"copy\" onclick=\"copybutton('"+"listbox_"+list[i]+"','"+tablename+"_"+list[i]+"');\">"+"</td>"
        data += "</tr>"
    }
    data += "</table>"
    data += "<input type=\"button\" value=\"add\" onclick=\"sendAddForm();closeViewForm('"+output+"')\">"
    document.getElementById(output).innerHTML = data
}

function copybutton(inputid,uptputid) {
    var input = document.getElementById(inputid)
    var output = document.getElementById(uptputid)
    output.value = input.value
}

function createlistbox(name) {
    var output = "<select name=\"\" id=\""+"listbox_"+name+"\">"
    for (var i=0;i<tablelist.length;i++) {
        if (tablelist[i][name]!=""){
            output += "<option value=\""+tablelist[i][name]+"\">"+tablelist[i][name]+"</option>"
        }
    }
    output += "</select>"
    return output
}

function serchckbox(key,name){
    var data = document.getElementById("listbox_"+name)
    var output = ""
    if (key=="") {
        for (var i=0;i<tablelist.length;i++) {
            if (tablelist[i][name]!=""){
                output += "<option value=\""+tablelist[i][name]+"\">"+tablelist[i][name]+"</option>"
            }
        }
        data.innerHTML = output
    }else {
        for (var i=0;i<tablelist.length;i++) {
            if (tablelist[i][name].substr(0,key.length).toLowerCase()==key.toLowerCase()){
                if (tablelist[i][name]!=""){
                    output += "<option value=\""+tablelist[i][name]+"\">"+tablelist[i][name]+"</option>"
                }
            }
        }
        data.innerHTML = output
    }
}

function sendAddForm() {
    var list = table_list[selectdata]
    var type = table_list_type[selectdata]
    var tablename = searchurl[selectdata]
    var jsondata = {};
    for (var i=0;i<list.length;i++) {
        if (list[i] == "Id") {
            continue
        }
        jsondata[list[i]] = document.getElementById(tablename+"_"+list[i]).value
        if (type[i] != "string") {
            jsondata[list[i]] = jsondata[list[i]] -0
        }
    }
    var url = HOSTURL + "/v1/add/"+tablename

    var req = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    req.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
        if(req.readyState == 4 && req.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
            var data = req.responseText;
            var jata = JSON.parse(data);
            console.log(jata);		          // 取得した JSON ファイルの中身を表示
            getTablelist("output");
        }else if (req.readyState == 4 && req.status != 200){ 
            var data = req.responseText;
            var jata = JSON.parse(data);
            console.log(jata);		          // 取得した JSON ファイルの中身を表示
        }
    };
    req.open("POST", url, true); // HTTPメソッドとアクセスするサーバーの　URL　を指定
    req.send(JSON.stringify(jsondata));					    // 実際にサーバーへリクエストを送信
}

function destory(id) {
    myRet = confirm("destory id="+id+" OK??");
    var tablename = searchurl[selectdata]
    if (myRet) {
        var url = HOSTURL + "/v1/edit/" + tablename + "/" + id;

        var req = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
        req.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
            if(req.readyState == 4 && req.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
                var data = req.responseText;
                var jata = JSON.parse(data);
                console.log(jata);		          // 取得した JSON ファイルの中身を表示
                getTablelist("output");
            }else if (req.readyState == 4 && req.status != 200){ 
                var data = req.responseText;
                var jata = JSON.parse(data);
                console.log(jata);		          // 取得した JSON ファイルの中身を表示
            }
        };
        req.open("DELETE", url, true); // HTTPメソッドとアクセスするサーバーの　URL　を指定
        req.setRequestHeader('content-type', 'application/x-www-form-urlencoded;charset=UTF-8');
        req.send(null);					    // 実際にサーバーへリクエストを送信
    }
}


function vieweditforms(output,id) {
    closeViewForm(output)
    var tablename = searchurl[selectdata]
    var url = HOSTURL + "/v1/read/" + tablename + "/" + id

    var req = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    req.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
        if(req.readyState == 4 && req.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
            var data = req.responseText;
            var jata = JSON.parse(data);
            console.log(jata);		          // 取得した JSON ファイルの中身を表示
            createEditform(output,jata.Result)
        }else if (req.readyState == 4 && req.status != 200){ 
            var data = req.responseText;
            var jata = JSON.parse(data);
            console.log(jata);		          // 取得した JSON ファイルの中身を表示
        }
    };
    req.open("GET", url, true); // HTTPメソッドとアクセスするサーバーの　URL　を指定
    req.send(null);					    // 実際にサーバーへリクエストを送信
}

function createEditform(output,jdata) {
    viewaddform(output)
    var data = ""
    var list = table_list[selectdata]
    var tablename = searchurl[selectdata]
    var id = 0
    data += tablename + "<br>"
    data += "<table>"+"<tr>"+"<th>Key名</th>"+"<th>値</th>"+"</tr>"
    for (var i=0;i<list.length;i++) {
        if (list[i] == "Id") {
            data += "<tr>"
            data += "<td>" + list[i] + "</td>"
            data += "<td>"
            id = jdata[list[i]];
            data += jdata[list[i]]
            data += "</td>"
            data += "</tr>"
            continue
        }
        data += "<tr>"
        data += "<td>" + list[i] + "</td>"
        data += "<td>"
        data += "<input type=\"text\" id=\""+tablename+"_"+list[i]+"\" value=\""+jdata[list[i]]+"\">"
        data += "</td>"
        data += "</tr>"
    }
    data += "</table>"
    data += "<input type=\"button\" value=\"edit\" onclick=\"sendEditForm("+id+");closeViewForm('"+output+"')\">"
    data += "<input type=\"button\" value=\"clsoe\" onclick=\"closeViewForm('"+output+"')\">"
    document.getElementById(output).innerHTML = data
}

function sendEditForm(id){
    var list = table_list[selectdata]
    var tablename = searchurl[selectdata]
    var tabletype = table_list_type[selectdata]
    var jsondata = {};
    for (var i=0;i<list.length;i++) {
        if (list[i] == "Id") {
            jsondata["Id"] = id;
            continue
        }
        jsondata[list[i]] = document.getElementById(tablename+"_"+list[i]).value
        if (tabletype[i] != "string") {
            jsondata[list[i]] = jsondata[list[i]] -0;
        }
    }
    var url = HOSTURL + "/v1/edit/"+tablename+"/"+id

    var req = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    req.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
        if(req.readyState == 4 && req.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
            var data = req.responseText;
            var jata = JSON.parse(data);
            console.log(jata);		          // 取得した JSON ファイルの中身を表示
            getTablelist("output");
        }else if (req.readyState == 4 && req.status != 200){ 
            var data = req.responseText;
            var jata = JSON.parse(data);
            console.log(jata);		          // 取得した JSON ファイルの中身を表示
        }
    };
    req.open("POST", url, true); // HTTPメソッドとアクセスするサーバーの　URL　を指定
    req.send(JSON.stringify(jsondata));					    // 実際にサーバーへリクエストを送信

}