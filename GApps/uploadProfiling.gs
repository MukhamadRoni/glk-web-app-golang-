function doPost(e) {
  try {
    // ID Folder Google Drive Profiling yang Anda sediakan
    var folderId = "1JX9K0aAf01nWK1rcMu3AlMk2HCMR8gQM";
    var folder = DriveApp.getFolderById(folderId);

    // Parsing data JSON yang dikirim dari Backend Golang
    var data = JSON.parse(e.postData.contents);
    var filename = data.filename;
    var fileData = data.fileData; // base64 encoded
    var mimeType = data.mimeType;

    // Decode base64 menjadi byte array
    var decoded = Utilities.base64Decode(fileData);
    var blob = Utilities.newBlob(decoded, mimeType, filename);

    // Simpan file ke folder spesifik Profiling
    var file = folder.createFile(blob);

    // Set file agar bisa dilihat oleh siapa saja yang memiliki link (optional)
    file.setSharing(DriveApp.Access.ANYONE_WITH_LINK, DriveApp.Permission.VIEW);

    // Return informasi file yang berhasil diupload
    return ContentService.createTextOutput(JSON.stringify({
      status: "success",
      fileId: file.getId(),
      fileUrl: file.getUrl(),
      message: "File Profiling Skill berhasil diunggah"
    })).setMimeType(ContentService.MimeType.JSON);

  } catch (error) {
    // Tangani jika terjadi error
    return ContentService.createTextOutput(JSON.stringify({
      status: "error",
      message: error.toString()
    })).setMimeType(ContentService.MimeType.JSON);
  }
}
