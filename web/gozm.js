// ===============================================================
// GOZM - Go Z-Machine Engine
// WebAssembly frontend JavaScript code
// ===============================================================

let outArea;
let inputBox;

// WebAssembly Go interface
const go = new Go();

function textOut(text) {
  outArea.textContent += text;
  requestAnimationFrame(() => {
    outArea.scrollTop = outArea.scrollHeight;
  });
}

function requestInput() {
  inputBox.focus();
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
    // hide options that are not available
    const saveOption = document.querySelector("#saveOption");
    if (typeof save === "function") {
      saveOption.classList.remove("disabled");
    } else {
      saveOption.classList.add("disabled");
    }
    const loadOption = document.querySelector("#loadOption");
    if (typeof load === "function") {
      loadOption.classList.remove("disabled");
    } else {
      loadOption.classList.add("disabled");
    }
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

  boot();
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
}

function setTheme(theme) {
  outArea.className = `theme${theme}`;
  inputBox.className = `theme${theme}`;

  hideMenus();
}
