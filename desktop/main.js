const {
  app,
  BrowserWindow,
  Tray,
  Menu,
  nativeImage,
  nativeTheme,
  ipcMain,
} = require("electron");
const path = require("node:path");
const Registry = require("winreg");
const positioner = require("electron-traywindow-positioner");
const ffi = require("ffi-napi");

const isWindows = process.platform === "win32";
const isMacOS = process.platform === "darwin";

const voltaserveLib = ffi.Library("build/Release/voltaserve", {
  add: ["int", ["int", "int"]],
});

let tray;
let window;

app.whenReady().then(() => {
  ipcMain.handle("voltaserve:add", () => voltaserveLib.add(1, 1));
  createTray();
  createWindow();
  app.on("activate", () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createTray();
    }
  });
});

const createTray = async () => {
  tray = new Tray(await getIcon());

  if (isWindows) {
    tray.on("right-click", showWindow);
    tray.on("double-click", showWindow);
    tray.on("click", showWindow);

    nativeTheme.on("updated", async () => {
      tray.setImage(await getIcon());
    });
  }

  const contextMenu = Menu.buildFromTemplate([
    {
      label: "Open",
      click: () => showWindow(),
    },
    {
      label: "Quit",
      click: () => app.exit(),
    },
  ]);
  tray.setToolTip("Voltaserve");
  tray.setContextMenu(contextMenu);
};

const createWindow = () => {
  window = new BrowserWindow({
    width: 450,
    height: 600,
    frame: true,
    show: false,
    fullscreenable: false,
    resizable: false,
    minimizable: false,
    autoHideMenuBar: true,
    webPreferences: {
      backgroundThrottling: false,
      preload: path.join(__dirname, "preload.js"),
    },
  });
  window.setTitle("Voltaserve");
  window.loadURL(`file://${path.join(__dirname, "index.html")}`);
  window.on("close", (event) => {
    event.preventDefault();
    window.hide();
  });
};

const showWindow = () => {
  positioner.position(window, tray.getBounds());
  window.show();
  window.focus();
};

const getIcon = async () => {
  const isDark = await isDarkTheme();
  return nativeImage
    .createFromPath(
      isDark || isMacOS ? "assets/icon-dark.png" : "assets/icon.png"
    )
    .resize({ height: 16, width: 16 });
};

const isDarkTheme = async () => {
  if (isWindows) {
    return await isWindowsDarkTheme();
  } else {
    return nativeTheme.shouldUseDarkColors;
  }
};

const isWindowsDarkTheme = () => {
  try {
    const regKey = new Registry({
      hive: Registry.HKCU,
      key: "\\Software\\Microsoft\\Windows\\CurrentVersion\\Themes\\Personalize",
    });
    return new Promise((resolve, reject) => {
      regKey.get("SystemUsesLightTheme", (err, result) => {
        if (err) {
          reject(err);
        } else {
          resolve(result.value === "0x0");
        }
      });
    });
  } catch (error) {
    console.log(error);
    return false;
  }
};
