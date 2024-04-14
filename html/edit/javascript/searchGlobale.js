var searchurl = ["booknames","filelists","copyfiles"];        //テーブルの名称
var booknames_row = ["Id","Name","Title","Writer","Burand","Booktype","Ext"];    //booknamesのテーブルID
var booknames_row_type = ["int","string","string","string","string","string","string"];    
var copyfile_row = ["Id","Zippass","Filesize","Copyflag"];  //copyfileのテーブルID
var copyfile_row_type = ["int","string","int","int"];  
var filelists_row =["Id","Name","Pdfpass","Zippass","Tag"];   //filelistsのテーブルID
var filelists_row_type =["int","string","string","string","string"];   
var listdata_row = ["Name","Pdf", "PdfFlag", "Zip", "ZipFlag"]
var listdata_row_type = ["string","string", "string", "string", "string"]
var table_list = [booknames_row,filelists_row,copyfile_row,listdata_row]; //テーブル選択値によるテーブルID指定
var table_list_type = [booknames_row_type,filelists_row_type,copyfile_row_type,listdata_row_type]
var selectdata = 0;      //テーブルの選択値

var HOSTURL = "";        //検索先のURLについて
