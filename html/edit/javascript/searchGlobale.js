var searchurl = ["booknames","filelists","copyfile"];        //テーブルの名称
var booknames_row = ["Id","Name","Title","Writer","Burand","Booktype","Ext"];    //booknamesのテーブルID
var copyfile_row = ["Id","Zippass","Filesize","Copyflag"];  //copyfileのテーブルID
var filelists_row =["Id","Name","Pdfpass","Zippass","Tag"];   //filelistsのテーブルID
var table_list = [booknames_row,filelists_row,copyfile_row]; //テーブル選択値によるテーブルID指定
var selectdata = 0;      //テーブルの選択値

var HOSTURL = "";        //検索先のURLについて
