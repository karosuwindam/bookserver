function destory(id){
    myRet = confirm("destory id="+id+" OK??");
    if (myRet){
        var xhr = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
        xhr.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
        if(xhr.readyState == 4 && xhr.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
            var data = xhr.responseText;
            console.log(data);		          // 取得した ファイルの中身を表示
            var jsondata = JSON.parse(data)
            document.getElementById("answer").innerHTML = jsondata.result
        }
        };

        var url = "/v1/edit/"
        if (meta_suburl != ""){
        url += meta_suburl + "/"
        }
        url += id;
        xhr.open('DELETE', url, true);
        xhr.setRequestHeader('Authorization', 'bearer '+TOKEN);
        xhr.setRequestHeader('content-type', 'application/x-www-form-urlencoded;charset=UTF-8');
        xhr.send(null);
    }
}