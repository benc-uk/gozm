// ===============================================================
// GOZM - Go Z-Machine Engine
// WebAssembly frontend JavaScript code
// ===============================================================

const MAX_OUTBUFFER = 8000;

// WebAssembly Go interface
const go = new Go();
let prefs = {};
let outArea;
let inputBox;
let hist = [];
let histIndex = -1;

// Initiate loading a story file and starting the Go WASM module
// Called from JS only
async function openFile(filename, filedata) {
  const result = await WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject);
  if (!result) {
    alert("Failed to load WebAssembly module");
    return;
  }
  boot();

  document.querySelector("#fileMenu").style.display = "none";
  go.argv = [filename];

  // Start the Go program in the background
  const runPromise = go.run(result.instance);

  // Give Go a moment to set up its exported functions
  await new Promise((resolve) => setTimeout(resolve, 50));

  if (filedata) {
    // Pass the file data to the Go module
    console.log("Passing file data to WASM, size:", filedata.length);
    if (typeof receiveFileData === "function") {
      // This function is defined in the Go code,
      // BUT won't be visible until after go.run() is called!
      receiveFileData(filedata);
    } else {
      console.error("receiveFileData function not available");
    }
  }

  // Wait for the program to complete
  await runPromise;

  textOut("Program has exited. Load another file\n");
  inputBox.style.visibility = "hidden";
  inputBox.blur();
}

// Called from Go to send text to the screen
function textOut(text) {
  outArea.textContent += text;
  requestAnimationFrame(() => {
    outArea.scrollTop = outArea.scrollHeight;
  });

  // Trim output buffer if too large
  if (outArea.textContent.length > MAX_OUTBUFFER) {
    outArea.textContent = outArea.textContent.slice(-MAX_OUTBUFFER);
  }
}

// Called from Go to request user input
function requestInput(history) {
  hist = history || [];
  histIndex = hist.length;

  // scroll input box into view and focus
  inputBox.scrollIntoView();
  inputBox.focus();

  requestAnimationFrame(() => {
    outArea.scrollTop = outArea.scrollHeight;
  });
}

// Called from both Go and JS
function clearScreen() {
  outArea.textContent = "";
}

// Called from Go when file is loaded and program is running
function loadedFile(filename) {
  clearScreen();
  inputBox.style.visibility = "visible";
  requestInput();

  console.log("Loaded file:", filename);
  prefs.loadedFile = filename;
  localStorage.setItem("prefs", JSON.stringify(prefs));
}

// Reset the system to initial state
function reset() {
  prefs.loadedFile = null;
  localStorage.setItem("prefs", JSON.stringify(prefs));
  location.reload();
}

// Called from JS to fake a boot sequence
function boot() {
  clearScreen();
  inputBox.value = "";
  textOut("System booting...\n");
  textOut("64K dynamic memory available\n");
  textOut("I/O buffers flushed\n\n");
  textOut("WASM subsystem initializing... complete!\n");
  textOut(`GOZM v${version} © Ben Coleman 2025\n\n`);
  textOut("Open a file to begin\n");
}

function promptFile() {
  hideMenus();
  const input = document.createElement("input");
  input.type = "file";
  input.accept = ".z3";
  input.onchange = async (e) => {
    const file = e.target.files[0];
    if (!file) {
      return;
    }

    const arrayBuffer = await file.arrayBuffer();
    openFile("tempFile", new Uint8Array(arrayBuffer));
  };

  input.click();
}

// When DOM is loaded, initialize everything, this is our entry point
window.addEventListener("DOMContentLoaded", async () => {
  inputBox = document.querySelector("input");
  outArea = document.querySelector("pre");
  const fileMenuButton = document.querySelector("#fileButton");
  const sysMenuButton = document.querySelector("#sysButton");
  const prefMenuButton = document.querySelector("#prefsButton");
  const infoMenuButton = document.querySelector("#infoButton");

  // hide input box initially
  inputBox.style.visibility = "hidden";

  fileMenuButton.onclick = function () {
    showMenu("fileMenu");
  };

  sysMenuButton.onclick = function () {
    showMenu("sysMenu");
  };

  prefMenuButton.onclick = function () {
    showMenu("prefsMenu");
  };

  infoMenuButton.onclick = function () {
    showMenu("infoMenu");
  };

  inputBox.onkeydown = function (e) {
    if (e.key === "Enter") {
      inputSend(inputBox.value);
      inputBox.value = "";
    }

    if (e.key === "ArrowUp") {
      if (histIndex > 0) {
        histIndex--;
        inputBox.value = hist[histIndex];
      }
      e.preventDefault();
    }

    if (e.key === "ArrowDown") {
      if (histIndex < hist.length - 1) {
        histIndex++;
        inputBox.value = hist[histIndex];
      } else {
        histIndex = hist.length;
        inputBox.value = "";
      }
      e.preventDefault();
    }
  };

  inputBox.onfocus = function () {
    requestAnimationFrame(() => {
      outArea.scrollTop = outArea.scrollHeight;
    });
  };

  // detect rezize to adjust scroll
  window.addEventListener("resize", () => {
    requestAnimationFrame(() => {
      outArea.scrollTop = outArea.scrollHeight;
    });
  });

  // get saved theme from localStorage
  try {
    const savedTheme = localStorage.getItem("prefs");
    if (!savedTheme) {
      throw new Error("No saved preferences, using defaults");
    }
    prefs = JSON.parse(savedTheme);
  } catch (e) {
    prefs = {
      theme: "RetroGlow",
      loadedFile: null,
    };

    localStorage.setItem("prefs", JSON.stringify(prefs));
  }

  setTheme(prefs.theme || "RetroGlow");

  if (prefs.loadedFile) {
    // Auto-load last file
    openFile(prefs.loadedFile);
  } else {
    boot();
  }
});

// Menu system
function showMenu(id) {
  hideMenus();
  menuDiv = document.querySelector("#" + id);
  if (menuDiv) {
    menuDiv.style.display = "block";
  }

  // Disable all with class disableNoFile if no file loaded
  const noFileLoaded = typeof load !== "function";
  const disableItems = document.querySelectorAll(".disableNoFile");
  disableItems.forEach((item) => {
    item.style.pointerEvents = noFileLoaded ? "none" : "auto";
    item.style.color = noFileLoaded ? "#888888" : "#000000";
  });
}

// Menu system
function hideMenus() {
  document.querySelector("#fileMenu").style.display = "none";
  document.querySelector("#sysMenu").style.display = "none";
  document.querySelector("#prefsMenu").style.display = "none";
  document.querySelector("#infoMenu").style.display = "none";
}

// Change theme of the fake terminal UI
function setTheme(theme) {
  outArea.className = `theme${theme}`;
  inputBox.className = `theme${theme}`;

  prefs.theme = theme;
  localStorage.setItem("prefs", JSON.stringify(prefs));

  hideMenus();
}

// Print about information
function printAbout() {
  clearScreen();
  hideMenus();
  textOut(`About:\nGo Z-Machine Engine (GOZM) v${version}\n`);
  textOut("A Z-Machine interpreter written in Go, compiled to WebAssembly.\n");
  textOut("© Ben Coleman 2025\n");
  textOut("Github: https://github.com/benc-uk/gozm\n\n");
  textOut("See the GitHub repository for full license, acknowledgments etc.\n");
}

// Print help information
function printHelp() {
  clearScreen();
  hideMenus();
  textOut("Help:\n");
  textOut("Use the File menu to open a Z-Machine story file (.z3),\n");
  textOut("  or open a supplied adventure file.\n");
  textOut("Type your commands in the input box and press Enter to send.\n");
  textOut("Use the System menu to save/load game state.\n");
  textOut("Use the Preferences menu to change themes.\n");
  textOut("");
  textOut("System commands:\n");
  textOut("  /quit - Exit the game\n");
  textOut("  /restart - Restart the game\n");
  textOut("  /save - Save the game\n");
  textOut("  /load - Load a saved game\n");
}
