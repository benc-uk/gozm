// ===============================================================
// GOZM - Go Z-Machine Engine
// WebAssembly frontend JavaScript code
// ===============================================================

let outArea;
let inputBox;
const MAX_OUTBUFFER = 8000;

// WebAssembly Go interface
const go = new Go();

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

function requestInput() {
  // scroll input box into view and focus
  inputBox.scrollIntoView();
  inputBox.focus();
  requestAnimationFrame(() => {
    outArea.scrollTop = outArea.scrollHeight;
  });
}

function clearScreen() {
  outArea.textContent = "";
}

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

function loadedFile() {
  clearScreen();
  inputBox.style.visibility = "visible";
}

function boot() {
  clearScreen();
  inputBox.value = "";
  textOut("System booting...\n");
  textOut("64K dynamic memory available.\n");
  textOut("I/O buffers flushed.\n\n");
  textOut("WASM subsystem initializing... complete.\n");
  textOut(`Go Z-Machine Engine\nGOZM v${version} © Ben Coleman 2025\n\n`);
  textOut("Open a file to begin.\n");
}

function hideMenus() {
  document.querySelector("#fileMenu").style.display = "none";
  document.querySelector("#sysMenu").style.display = "none";
  document.querySelector("#prefsMenu").style.display = "none";
  document.querySelector("#infoMenu").style.display = "none";
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

// When DOM is loaded
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
  };
  inputBox.onfocus = function () {
    requestAnimationFrame(() => {
      outArea.scrollTop = outArea.scrollHeight;
    });
  };

  boot();

  // check query params for story file URL
  const urlParams = new URLSearchParams(window.location.search);
  const storyUrl = urlParams.get("story");
  if (storyUrl) {
    try {
      const response = await fetch(storyUrl);
      if (response.ok) {
        const arrayBuffer = await response.arrayBuffer();
        openFile(storyUrl, new Uint8Array(arrayBuffer));
      } else {
        console.error(`Failed to load story from URL: ${storyUrl}`);
      }
    } catch (error) {
      console.error(`Error loading story from URL: ${storyUrl}`, error);
    }
  }
});

function reset() {
  window.location.reload();
}

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

function setTheme(theme) {
  outArea.className = `theme${theme}`;
  inputBox.className = `theme${theme}`;

  hideMenus();
}

function printAbout() {
  clearScreen();
  hideMenus();
  textOut(`About:\nGo Z-Machine Engine (GOZM) v${version}\n`);
  textOut("A Z-Machine interpreter written in Go, compiled to WebAssembly.\n");
  textOut("© Ben Coleman 2025\n");
  textOut("Github: https://github.com/benc-uk/gozm\n");
}

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
