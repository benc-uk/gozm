let outArea;
let inputBox;

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

async function openFile(filename) {
  const go = new Go();
  const result = await WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject);
  if (!result) {
    alert("Failed to load WebAssembly module");
    return;
  }

  boot();
  document.querySelector("#fileMenu").style.display = "none";
  go.argv = [filename];
  await go.run(result.instance);

  textOut("Program has exited. Load another file\n");
  inputBox.style.visibility = "hidden";
  inputBox.blur();
}

function loadedFile() {
  inputBox.style.visibility = "visible";
}

function boot() {
  clearScreen();
  inputBox.value = "";
  textOut("Power on sequence...\n");
  textOut("System booting...\n");
  textOut("64K dynamic memory available.\n");
  textOut("I/O buffers flushed.\n\n");
  textOut("Loading: WASM subsystem... complete.\n");
  textOut("GOZM v0.1.0 - Go Z-Machine Engine (c) Ben Coleman 2025\n\n");
  textOut("Open a file to begin.\n");
}

function hideMenus() {
  document.querySelector("#fileMenu").style.display = "none";
}

// When DOM is loaded
window.addEventListener("DOMContentLoaded", async () => {
  inputBox = document.querySelector("input");
  outArea = document.querySelector("pre");
  const fileMenuButton = document.querySelector("#fileButton");

  // hide input box initially
  inputBox.style.visibility = "hidden";

  fileMenuButton.onclick = function () {
    const menu = document.querySelector("#fileMenu");
    if (menu.style.display === "block") {
      menu.style.display = "none";
    } else {
      menu.style.display = "block";
    }
  };

  inputBox.onkeydown = function (e) {
    if (e.key === "Enter") {
      inputSend(inputBox.value);
      inputBox.value = "";
    }
  };

  boot();
});
