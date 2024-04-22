var lefte =[0];
var mide = [1,3,4];
var righte = [2];
var mouseover_f = true;     //mouns over flag
var view_message = ["Next","1 page next","back","1/2 page change","page view"]

function onclickevent(num)
{
    switch(num){
        case 0: //2 page up
            viewPageChange(2)
            break;
        case 1: //1 page up
            viewPageChange(1)
            break;
        case 2: //2 page down
            viewPageChange(-2)
            break;
        case 3: //page change 1 to 2
            chPageOneTwe(2)
            break;
        case 4: //page bear view on/off
            pageviewonoff();
            break;
        default:
            break;
    }
}


document.onkeydown = function(event){
    var keyEvent = event||window.event;
    console.log("KEYEVENT:"+keyEvent.keyCode);
    switch (keyEvent.keyCode){
        case 37:			//←
            viewPageChange(+2);
            break;
        case 38:			//↑
            viewPageChange(+1);
            break;
        case 39:			//→
            viewPageChange(-2);
            break;
        case 40:			//↓
            viewPageChange(-1);
            break;
        case 72:            //h
            helpviewonoff();
            break;
        case 73:        //i
            infoviewonoff();
            break;
        case 77:        //m
            mouseover_f = !mouseover_f;
            break;
        case 80:        //p
            pageviewonoff()
            break;
        case 32:       //space
            onclickevent(3);
            break;
        default:
            break;
    }
}

function initload(){
    var data = ["leftT","midT","rightT","midM","midB"];
    for(var i=0;i<data.length;i++){
        var tmp = document.getElementsByClassName(data[i]);
        tmp[0].style.border = "5px solid  #000000";
        tmp[0].innerHTML = view_message[i];
    }//border: inset 10px #ff0000;
    mousescroll_int();
    // nowpage(pagenow);
}
var onloadclose = function(){     //起動時を非表示項目
    var info = document.getElementsByClassName("help");
    var page = document.getElementsByClassName("nowpage");
    var sidber = document.getElementById("pageslider");
    info[0].style.display = "none"
    page[0].style.display = "none"
    sidber.style.display = "none"

    var data = ["leftT","midT","rightT","midM","midB"];
    for(var i=0;i<data.length;i++){
        var tmp = document.getElementsByClassName(data[i]);
        tmp[0].style.border = "";
        tmp[0].innerHTML = "";
    }
}
function helpviewonoff(){
    var help = document.getElementsByClassName("help");
    if (help[0].style.display == ""){
        help[0].style.display = "none"
    }else{
        help[0].style.display = ""
    }
}


/*
 * スワイプイベント設定
 */
var swapflag = false
function setSwipe() {
	let startX;		// タッチ開始 x座標
	let startY;		// タッチ開始 y座標
	let moveX;	// スワイプ中の x座標
	let moveY;	// スワイプ中の y座標
	let dist = 30;	// スワイプを感知する最低距離（ピクセル単位）
	var sflag=false;
	
	// タッチ開始時： xy座標を取得
	window.addEventListener("touchstart", function(e) {
		e.preventDefault();
		startX = e.touches[0].pageX;
		startY = e.touches[0].pageY;
		sflag=true;
	});
	
	// スワイプ中： xy座標を取得
	window.addEventListener("touchmove", function(e) {
		e.preventDefault();
		moveX = e.changedTouches[0].pageX;
		moveY = e.changedTouches[0].pageY;
		sflag=false;
	});
	
	
	// タッチ終了時： スワイプした距離から左右どちらにスワイプしたかを判定する/距離が短い場合何もしない
	window.addEventListener("touchend", function(e) {
		//スワイプ移動しない場合の処理
		if ((sflag)||(swapflag)||(list_flag)){
			moveX = startX;
			moveY = startY;
		}

        if (startX > moveX && startX > moveX + dist) {		// 右から左にスワイプ
            // console.log("left");
            viewPageChange(-2);
		}
		else if (startX < moveX && startX + dist < moveX) {	// 左から右にスワイプ
            // console.log("right");
            viewPageChange(+2);
		}
		swapflag = false;
	});
}
window.addEventListener("load", function(){

	// スワイプイベント設定
	setSwipe();
});

//pageバーのON/OFF
function pageviewonoff(){
    var info = document.getElementsByClassName("nowpage");
    var sidber = document.getElementById("pageslider");
    if (info[0].style.display == ""){
        info[0].style.display = "none"
        sidber.style.display = "none"
    }else{
        info[0].style.display = ""
        sidber.style.display = ""
    }
}


function mousescroll_int(){
    //Firefox
    if(window.addEventListener){
        window.addEventListener('DOMMouseScroll', function(e){
            // alert(e.detail);
            if((document.getElementsByClassName("listdata")[0].style.display == "none")){
            // &&(document.getElementsByClassName("output")[0].style.display == "none")){
                if (e.detail < 0){
                    viewPageChange(-2);
                }else{
                    viewPageChange(+2);
                }
            }
        }, false);
    }
    
    //IE
    if(document.attachEvent){
        document.attachEvent('onmousewheel', function(e){
            // alert(e.wheelDelta);
            if((document.getElementsByClassName("listdata")[0].style.display == "none")){
            // &&(document.getElementsByClassName("output")[0].style.display == "none")){
                if (e.wheelDelta > 0){
                    viewPageChange(-2);
                }else{
                    viewPageChange(+2);
                }
            }
        });
    }
    
    //Chrome
    window.onmousewheel = function(e){
        // alert(e.wheelDelta);
        if((document.getElementsByClassName("listdata")[0].style.display == "none")){
        // &&(document.getElementsByClassName("output")[0].style.display == "none")){
            if (e.wheelDelta > 0){
                viewPageChange(-2);
            }else{
                viewPageChange(+2);
            }
        }
    }

}