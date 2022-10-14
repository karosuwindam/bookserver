var bookname_row = ["id","name","title","writer","brand","booktype","ext"];
var copyfile_row = ["id","zippass","filesize","copyflag"];
var filelists_row =["id","name","pdfpass","zippass"];

function bookname_edit(ary){
    var output =[]
    for (var i=0;i<ary.length;i++){
        var tmp =[]
        tmp.push(ary[i].id,ary[i].name,ary[i].title,ary[i].Writer,ary[i].brand,ary[i].ext)
        output.push(tmp)
    }
    return output
}
function copyfile_edit(ary){
    var output =[]
    for (var i=0;i<ary.length;i++){
        var tmp =[]
        tmp.push(ary[i].id,ary[i].Zippass,ary[i].Filesize,ary[i].copyflag)
        output.push(tmp)
    }
    return output
}
function filelist_edit(ary){
    var output =[]
    for (var i=0;i<ary.length;i++){
        var tmp =[]
        tmp.push(ary[i].id,ary[i].name,ary[i].pdfpass,ary[i].Zippass,ary[i].tag)
        output.push(tmp)
    }
    return output
}
function outputTable(ary){
    var output
    var sumdata = 0
    var outf =false
    output = "<table>"
    for(var i=0;i<ary.length;i++){
        output += "<tr style='color:#FFFFFF'>";
        var tmp = ary[i]
        for(var j=0;j<tmp.length;j++){
            output += "<td>"+tmp[j]+"</td>";
            outf = true;
        }
        if(selectdata == 2){
            sumdata += tmp[2]-0
            output += "<td><input type='checkbox' "
            if (tmp[3] == "1"){
              output += "checked='checked' "
            }
            output += "onclick=\"ck_copyfilebox(\'"+tmp[1]+"\',this.checked)\""
            output += ">"+"</td>"
        }
        output += "</tr>"
    }
    output += "</table>"
    if((selectdata == 2)){
        output += "<br>"
        if (sumdata > 1024*1024*1024){
            output += "<div>"+"sumsize:"+Math.round(sumdata/1024/1024/1024*1000)/1000 +"G"+"</div>"
        }else if (sumdata > 1024*1024){
            output += "<div>"+"sumsize:"+Math.round(sumdata/1024/1024*1000)/1000 +"M"+"</div>"
        }else if (sumdata > 1024){
            output += "<div>"+"sumsize:"+Math.round(sumdata/1024*1000)/1000 +"K"+"</div>"
        }else{
            output += "<div>"+"sumsize:"+sumdata+"</div>"
        }
    }
    if (!outf){
        output = ""
    }
    return output
}
