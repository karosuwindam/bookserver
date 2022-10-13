var TOKEN = "";

function readstoragesession() {
  var session = localStorage.getItem("hello-session");
  var tokendata = localStorage.getItem("tokendata")
  if (session != "") {
    document.cookie = "hello-session="+session;

  }
  if (tokendata != "") {
    TOKEN = tokendata
  }

}

function writestoragesession(){
  var cookies = document.cookie; //全てのcookieを取り出して
  var cookiesArray = cookies.split(';'); // ;で分割し配列に
  if (cookies == "") {
    localStorage.setItem('hello-session', '');
    return
  }

  for(var c of cookiesArray){ //一つ一つ取り出して
    var cArray = c.split('='); //さらに=で分割して配列に
    if( cArray[0] == 'hello-session'){ // 取り出したいkeyと合致したら
        // console.log(cArray);  // [key,value] 
        var out = cArray[1]
        for(var i=2;i<cArray.length;i++){
          out += "="
        }
        localStorage.setItem('hello-session', out);

    }
  }
  localStorage.setItem('tokendata', TOKEN);

}

function login(name,pass) {
    var xhr = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    xhr.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
      if(xhr.readyState == 4 && xhr.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
          var data = xhr.responseText;
          console.log(data);		          // 取得した ファイルの中身を表示
          var tmp = JSON.parse(data)
          TOKEN = tmp.result.token
          writestoragesession()
      }else if(xhr.readyState == 4 && xhr.status != 200){
          var data = xhr.responseText;
          console.log(data);		          // 取得した ファイルの中身を表示
          if (TOKEN != "") {
            logout();
          }
  
      }
    };
    xhr.open('POST','/login');
    xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    var request = "name=" + name + "&pass=" + pass;
    xhr.send(request);
}

function loginget(){
    var xhr = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    xhr.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
      if(xhr.readyState == 4 && xhr.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
          var data = xhr.responseText;
          console.log(data);		          // 取得した ファイルの中身を表示
          var tmp = JSON.parse(data)
          writestoragesession()
          TOKEN = tmp.result.token
      }else if(xhr.readyState == 4 && xhr.status != 200){
        var data = xhr.responseText;
        console.log(data);		          // 取得した ファイルの中身を表示
        if (TOKEN != "") {
          logout();
        }
      }
    };
    xhr.open('GET','/login');
    xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    xhr.send(null);

}

function logout() {
    var xhr = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    xhr.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
      if(xhr.readyState == 4 && xhr.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
          var data = xhr.responseText;
          console.log(data);		          // 取得した ファイルの中身を表示
          var tmp = JSON.parse(data)
          TOKEN = ""
          document.cookie = "hello-session=; max-age=0";
          writestoragesession();
      }
    };
    xhr.open('POST','/logout');
    xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    xhr.send(null);

}