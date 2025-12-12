// ===============================================================
// GOZM - Input handling module
// Terminal-style input with mobile keyboard support
// ===============================================================

const isMobile = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)

let outArea
let hiddenInput
let hist = []
let histIndex = -1
let inputBuffer = ''
let inputActive = false
let onSubmit = null // Callback when input is submitted

// Initialize input handling
export function initInput(outputElement, hiddenInputElement, submitCallback) {
  outArea = outputElement
  hiddenInput = hiddenInputElement
  onSubmit = submitCallback

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
}

// Called to request user input
export function requestInput(history) {
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
    submitInput()
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
    submitInput()
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

// Submit the current input
function submitInput() {
  inputActive = false
  removeInputDisplay()
  const submitted = inputBuffer
  inputBuffer = ''
  if (onSubmit) {
    onSubmit(submitted)
  }
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
export function removeInputDisplay() {
  const cursor = document.getElementById('termCursor')
  if (cursor) cursor.remove()
  const inputSpan = document.getElementById('termInput')
  if (inputSpan) inputSpan.remove()
}
