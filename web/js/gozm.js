// ===============================================================
// GOZM - Go Z-Machine Engine
// WebAssembly frontend JavaScript code
// ===============================================================

import { version } from './version.js'
import { initMenus } from './menus.js'
import { initInput, requestInput, removeInputDisplay } from './input.js'

const MAX_OUTBUFFER = 8000

// WebAssembly Go interface
const go = new Go()
let prefs = {}
let outArea
let modal

// Two way bridge between Go and JS
window.bridge = {
  textOut: textOut,
  requestInput: requestInput,
  loadedFile: loadedFile,
  playSound: playSound,
  // These are stubs to be replaced by Go when the module is running
  save: null,
  load: null,
  getInfo: null,
  inputSend: null,
  receiveFileData: null,
}

// When DOM is loaded, initialize everything, this is our entry point
window.addEventListener('DOMContentLoaded', async () => {
  outArea = document.querySelector('pre')
  modal = document.getElementById('modal')
  const hiddenInput = document.getElementById('hiddenInput')
  initMenus()

  // Initialize input handling with submit callback
  initInput(outArea, hiddenInput, (text) => {
    // textOut(text + '\n')
    bridge.inputSend(text)
  })

  // detect resizes to adjust scroll
  window.addEventListener('resize', () => {
    requestAnimationFrame(() => {
      outArea.scrollTop = outArea.scrollHeight
    })
  })

  // Initialize preferences from localStorage
  try {
    const savedTheme = localStorage.getItem('prefs')
    if (!savedTheme) {
      throw new Error('No saved preferences, using defaults')
    }
    prefs = JSON.parse(savedTheme)
  } catch (e) {
    prefs = {
      theme: 'RetroGlow',
      loadedFile: null,
    }

    localStorage.setItem('prefs', JSON.stringify(prefs))
  }

  setTheme(prefs.theme || 'RetroGlow')

  if (prefs.loadedFile) {
    openFile(prefs.loadedFile) // Auto-load last file
  } else {
    boot()
  }
})

// Initiate loading a story file and starting the Go WASM module
export async function openFile(filename, filedata) {
  const result = await WebAssembly.instantiateStreaming(fetch('main.wasm'), go.importObject)
  if (!result) {
    alert('Failed to load WebAssembly module')
    return
  }
  boot()

  // Start the Go program in the background
  go.argv = [filename]
  const runPromise = go.run(result.instance)

  if (filedata && bridge.receiveFileData) {
    // Pass the file data to the Go module
    console.log('Passing file data to WASM, size:', filedata.length)
    bridge.receiveFileData(filedata)
  }

  // Wait for the program to complete, this will block until exit
  await runPromise

  textOut('Program has exited. Load another file\n')
}

// Called from Go to send text to the screen
export function textOut(text) {
  // Temporarily remove cursor elements before modifying textContent
  removeInputDisplay()

  outArea.textContent += text

  // Trim output buffer if too large
  if (outArea.textContent.length > MAX_OUTBUFFER) {
    outArea.textContent = outArea.textContent.slice(-MAX_OUTBUFFER)
  }

  requestAnimationFrame(() => {
    outArea.scrollTop = outArea.scrollHeight
  })
}

// Called from both Go and JS
export function clearScreen() {
  outArea.textContent = ''
}
// Reset the system to initial state
export function reset() {
  prefs.loadedFile = null
  localStorage.setItem('prefs', JSON.stringify(prefs))
  location.reload()
}

// Start a file open dialog to select a Z-Machine story file
export function promptFile() {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.z3'
  input.onchange = async (e) => {
    const file = e.target.files[0]
    if (!file) {
      return
    }

    const arrayBuffer = await file.arrayBuffer()
    openFile('tempFile', new Uint8Array(arrayBuffer))
  }

  input.click()
}

// Change theme of the fake terminal UI
export function setTheme(theme) {
  outArea.className = `theme${theme}`
  modal.className = `theme${theme}`

  prefs.theme = theme
  localStorage.setItem('prefs', JSON.stringify(prefs))
}

// Called from Go when file is loaded and program is running
function loadedFile(filename) {
  clearScreen()
  requestInput()

  console.log('Loaded file:', filename)
  prefs.loadedFile = filename
  localStorage.setItem('prefs', JSON.stringify(prefs))
}

// Called from JS to fake a boot sequence
function boot() {
  clearScreen()
  textOut('System booting...\n')
  textOut('64K dynamic memory available\n')
  textOut('I/O buffers flushed\n\n')
  textOut('WASM subsystem initializing... complete!\n')
  textOut(`GOZM v${version} © Ben Coleman 2025\n\n`)
  textOut('Open a file to begin\n')
}

function playSound(soundID, effect, vol) {
  console.log('STUB! Play sound requested:', soundID, effect, vol)
}

export function showModal(message) {
  const modal = document.getElementById('modal')
  modal.firstChild.textContent = message
  modal.style.display = 'block'
}
