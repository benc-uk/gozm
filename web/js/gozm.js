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
let hiddenInput
let hist = []
let histIndex = -1
let inputBuffer = ''
let inputActive = false
const isMobile = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)

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
  outArea = document.querySelector('pre')
  hiddenInput = document.getElementById('hiddenInput')
  initMenus()

  // Capture keyboard input on the document for terminal-style input
  document.addEventListener('keydown', handleKeyDown)

  // Handle paste events
  document.addEventListener('paste', handlePaste)

  // Click on output area to focus - shows keyboard on mobile
  outArea.addEventListener('click', () => {
    if (inputActive && isMobile) {
      hiddenInput.focus()
    } else {
      outArea.focus()
    }
  })

  // Handle hidden input for mobile
  if (hiddenInput) {
    hiddenInput.addEventListener('input', handleMobileInput)
    hiddenInput.addEventListener('keydown', handleMobileKeyDown)
  }

  // Make outArea focusable
  outArea.tabIndex = 0

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

  inputActive = false
  removeInputDisplay()
  textOut('Program has exited. Load another file\n')
}

// Called from Go to send text to the screen
export function textOut(text) {
  // Temporarily remove cursor elements before modifying textContent
  const cursor = document.getElementById('termCursor')
  const inputSpan = document.getElementById('termInput')
  if (cursor) cursor.remove()
  if (inputSpan) inputSpan.remove()

  outArea.textContent += text

  // Trim output buffer if too large
  if (outArea.textContent.length > MAX_OUTBUFFER) {
    outArea.textContent = outArea.textContent.slice(-MAX_OUTBUFFER)
  }

  // Restore cursor if input is active
  if (inputActive) {
    updateInputDisplay()
  }

  requestAnimationFrame(() => {
    outArea.scrollTop = outArea.scrollHeight
  })
}

// Called from both Go and JS
export function clearScreen() {
  outArea.textContent = ''
}

// Called from Go to request user input
function requestInput(history) {
  hist = history || []
  histIndex = hist.length
  inputBuffer = ''
  inputActive = true

  // Show cursor
  updateInputDisplay()
  
  // Focus appropriate element based on platform
  if (isMobile && hiddenInput) {
    hiddenInput.value = ''
    hiddenInput.focus()
  } else {
    outArea.focus()
  }

  requestAnimationFrame(() => {
    outArea.scrollTop = outArea.scrollHeight
  })
}

// Handle input from hidden input on mobile
function handleMobileInput(e) {
  if (!inputActive) return
  inputBuffer = e.target.value
  updateInputDisplay()
}

// Handle special keys from hidden input on mobile
function handleMobileKeyDown(e) {
  if (!inputActive) return

  if (e.key === 'Enter') {
    e.preventDefault()
    inputActive = false
    removeInputDisplay()
    textOut(inputBuffer + '\n')
    bridge.inputSend(inputBuffer)
    inputBuffer = ''
    hiddenInput.value = ''
    hiddenInput.blur()
    return
  }
}

// Handle keyboard input for terminal-style typing
function handleKeyDown(e) {
  if (!inputActive) return

  // Don't capture if user is in a menu or other input
  if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') return

  if (e.key === 'Enter') {
    e.preventDefault()
    inputActive = false
    // Remove cursor and add newline
    removeInputDisplay()
    textOut(inputBuffer + '\n')
    bridge.inputSend(inputBuffer)
    inputBuffer = ''
    return
  }

  if (e.key === 'Backspace') {
    e.preventDefault()
    if (inputBuffer.length > 0) {
      inputBuffer = inputBuffer.slice(0, -1)
      updateInputDisplay()
    }
    return
  }

  if (e.key === 'ArrowUp') {
    e.preventDefault()
    if (histIndex > 0) {
      histIndex--
      inputBuffer = hist[histIndex] || ''
      updateInputDisplay()
    }
    return
  }

  if (e.key === 'ArrowDown') {
    e.preventDefault()
    if (histIndex < hist.length - 1) {
      histIndex++
      inputBuffer = hist[histIndex] || ''
    } else {
      histIndex = hist.length
      inputBuffer = ''
    }
    updateInputDisplay()
    return
  }

  // Ignore other control keys
  if (e.key.length > 1 || e.ctrlKey || e.metaKey || e.altKey) {
    return
  }

  // Add printable character
  e.preventDefault()
  inputBuffer += e.key
  updateInputDisplay()
}

// Handle paste events
function handlePaste(e) {
  if (!inputActive) return
  if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') return

  e.preventDefault()
  const text = e.clipboardData.getData('text')
  // Only take first line and filter to printable chars
  const firstLine = text.split('\n')[0].replace(/[^\x20-\x7E]/g, '')
  inputBuffer += firstLine
  updateInputDisplay()
}

// Update the display to show current input with cursor
function updateInputDisplay() {
  // Remove any existing input display
  const cursor = document.getElementById('termCursor')
  if (cursor) {
    cursor.remove()
  }

  // Remove previously displayed input text
  const inputSpan = document.getElementById('termInput')
  if (inputSpan) {
    inputSpan.remove()
  }

  // Create new input span with cursor
  const newInputSpan = document.createElement('span')
  newInputSpan.id = 'termInput'
  newInputSpan.textContent = inputBuffer

  const newCursor = document.createElement('span')
  newCursor.id = 'termCursor'
  newCursor.className = 'cursor'
  newCursor.textContent = '█'

  outArea.appendChild(newInputSpan)
  outArea.appendChild(newCursor)

  requestAnimationFrame(() => {
    outArea.scrollTop = outArea.scrollHeight
  })
}

// Remove the input display (cursor and typed text)
function removeInputDisplay() {
  const cursor = document.getElementById('termCursor')
  if (cursor) cursor.remove()
  const inputSpan = document.getElementById('termInput')
  if (inputSpan) inputSpan.remove()
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
  inputBuffer = ''
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
