function doPost(e) {
  try {
    // Parsing data JSON yang dikirim dari Backend Golang
    var data = JSON.parse(e.postData.contents);
    var filename = data.filename;
    var fileData = data.fileData; // base64 encoded
    var mimeType = data.mimeType;
    
    // Ganti dengan ID Folder Google Drive Anda
    // Anda bisa mendapatkan ID dari URL folder: https://drive.google.com/drive/folders/ID_FOLDER_INI
    var folderId = data.folderId || "ID_FOLDER_DRIVE_ANDA_DISINI"; 

    var folder = DriveApp.getFolderById(folderId);
    
    // Decode data base64 menjadi file
    var decodedData = Utilities.base64Decode(fileData);
    var blob = Utilities.newBlob(decodedData, mimeType, filename);
    
    // Buat file di folder tujuan
    var file = folder.createFile(blob);
    
    // Set file agar bisa diakses public (dilihat oleh siapa saja yang memiliki link)
    file.setSharing(DriveApp.Access.ANYONE_WITH_LINK, DriveApp.Permission.VIEW);
    
    // Kembalikan response JSON berisi URL dan ID file
    return ContentService.createTextOutput(JSON.stringify({
      status: "success",
      fileId: file.getId(),
      fileUrl: file.getUrl()
    })).setMimeType(ContentService.MimeType.JSON);
    
  } catch (error) {
    // Tangani jika terjadi error
    return ContentService.createTextOutput(JSON.stringify({
      status: "error",
      message: error.toString()
    })).setMimeType(ContentService.MimeType.JSON);
  }
}
