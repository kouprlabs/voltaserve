window.onload = () => {
  window.electronAPI.voltaserve.getFileList("/example").then((result) => {
    document.getElementById("get-file-list").innerText = JSON.stringify(result);
  });
  window.electronAPI.voltaserve
    .uploadFile("/example/file.txt")
    .then((result) => {
      document.getElementById("upload-file").innerText = JSON.stringify(result);
    });
};
