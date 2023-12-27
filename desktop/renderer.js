window.onload = async () => {
  const result = await window.electronAPI.voltaserve.add(1, 1);
  document.getElementById("result").innerText = result;
};
