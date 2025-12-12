// ===============================================================
// GOZM - Go Z-Machine Engine
// WebAssembly frontend JavaScript code
// ===============================================================

import { version } from './version.js'
import { initMenus } from './menus.js'

const MAX_OUTBUFFER = 8000

// WebAssembly Go interface
const go = new Go()
let prefs = {}
let outArea
let inputBox
let hist = []
let histIndex = -1

// Two way bridge between Go and JS
window.bridge = {
  textOut: textOut,
  requestInput: requestInput,
  loadedFile: loadedFile,
  playSound: playSound,
  // These are stubs to be replaced by Go when the module is running
  save: null,
  load: null,
  printInfo: null,
  inputSend: null,
  receiveFileData: null,
}

// When DOM is loaded, initialize everything, this is our entry point
window.addEventListener('DOMContentLoaded', async () => {
  inputBox = document.querySelector('input')
  outArea = document.querySelector('pre')
  initMenus()
  inputBox.style.display = 'none'

  inputBox.onkeydown = function (e) {
    if (e.key === 'Enter') {
      bridge.inputSend(inputBox.value)
      inputBox.value = ''
    }

    if (e.key === 'ArrowUp') {
      if (histIndex > 0) {
        histIndex--
        inputBox.value = hist[histIndex]
      }
      e.preventDefault()
    }

    if (e.key === 'ArrowDown') {
      if (histIndex < hist.length - 1) {
        histIndex++
        inputBox.value = hist[histIndex]
      } else {
        histIndex = hist.length
        inputBox.value = ''
      }
      e.preventDefault()
    }
  }

  inputBox.onfocus = function () {
    requestAnimationFrame(() => {
      outArea.scrollTop = outArea.scrollHeight
    })
  }

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
  inputBox.style.visibility = 'hidden'
  inputBox.blur()
}

// Called from Go to send text to the screen
export function textOut(text) {
  outArea.textContent += text
  requestAnimationFrame(() => {
    outArea.scrollTop = outArea.scrollHeight
  })

  // Trim output buffer if too large
  if (outArea.textContent.length > MAX_OUTBUFFER) {
    outArea.textContent = outArea.textContent.slice(-MAX_OUTBUFFER)
  }
}

// Called from both Go and JS
export function clearScreen() {
  outArea.textContent = ''
}

// Called from Go to request user input
function requestInput(history) {
  hist = history || []
  histIndex = hist.length

  // scroll input box into view and focus
  inputBox.scrollIntoView()
  inputBox.focus()

  requestAnimationFrame(() => {
    outArea.scrollTop = outArea.scrollHeight
  })
}

// Reset the system to initial state
export function reset() {
  prefs.loadedFile = null
  localStorage.setItem('prefs', JSON.stringify(prefs))
  location.reload()
}

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
  inputBox.className = `theme${theme}`

  prefs.theme = theme
  localStorage.setItem('prefs', JSON.stringify(prefs))
}

// Called from Go when file is loaded and program is running
function loadedFile(filename) {
  clearScreen()
  inputBox.style.display = 'block'
  requestInput()

  console.log('Loaded file:', filename)
  prefs.loadedFile = filename
  localStorage.setItem('prefs', JSON.stringify(prefs))
}

// Called from JS to fake a boot sequence
function boot() {
  clearScreen()
  inputBox.value = ''
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
