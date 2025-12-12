// ===============================================================
// GOZM - Go Z-Machine Engine
// Menus JavaScript code
// ===============================================================

import { openFile, promptFile, setTheme, reset, showModal } from './gozm.js'
import { version } from './version.js'

let fileMenu, sysMenu, prefsMenu, infoMenu

// Initialize menus and menu items
export function initMenus() {
  fileMenu = document.querySelector('#fileMenu')
  sysMenu = document.querySelector('#sysMenu')
  prefsMenu = document.querySelector('#prefsMenu')
  infoMenu = document.querySelector('#infoMenu')

  addMenuButton('#fileButton', fileMenu)
  addMenuButton('#sysButton', sysMenu)
  addMenuButton('#prefsButton', prefsMenu)
  addMenuButton('#infoButton', infoMenu)

  addMenuItem(fileMenu, 'Open...', promptFile)
  addMenuSeparator(fileMenu)
  addMenuItem(fileMenu, 'Mini Zork', async () => await openFile('minizork.z3'))
  addMenuItem(fileMenu, 'Moonglow', async () => await openFile('moonglow.z3'))
  addMenuItem(fileMenu, 'Catseye', async () => await openFile('catseye.z3'))
  addMenuItem(fileMenu, 'Adventure', async () => await openFile('advent.z3'))
  addMenuItem(fileMenu, 'Buccaneers Cache', async () => await openFile('buccaneers_cache.z3'))
  addMenuItem(fileMenu, 'Duck Me', async () => await openFile('duckme.z3'))
  addMenuSeparator(fileMenu)
  addMenuItem(fileMenu, 'Zork I', async () => await openFile('zork1-r119-s880429.z3'))
  addMenuItem(fileMenu, 'Zork II', async () => await openFile('zork2-r63-s860811.z3'))
  addMenuItem(fileMenu, 'Zork III', async () => await openFile('zork3-r25-s860811.z3'))
  addMenuItem(fileMenu, 'Enchanter', async () => await openFile('enchanter-r24-s851118.z3'))
  addMenuItem(fileMenu, "Hitchhiker's Guide", async () => await openFile('hitchhiker-r60-s861002.z3'))
  addMenuItem(fileMenu, 'The Lurking Horror', async () => await openFile('lurkinghorror-r221-s870918.z3'))
  addMenuItem(fileMenu, 'Planetfall', async () => await openFile('planetfall-r39-s880501.z3'))
  addMenuItem(fileMenu, 'Wishbringer', async () => await openFile('wishbringer-r69-s850920.z3'))

  // System Menu
  //prettier-ignore
  addMenuItem(sysMenu, 'Save', () => { bridge.save() }, true)
  //prettier-ignore
  addMenuItem(sysMenu, 'Restore', () => { bridge.load() }, true)
  addMenuSeparator(sysMenu)
  addMenuItem(sysMenu, 'Reset System', () => {
    reset()
  })

  // Preferences Menu
  addMenuItem(prefsMenu, 'Theme: Retro Green', () => setTheme('RetroGlow'))
  addMenuItem(prefsMenu, 'Theme: Retro Amber', () => setTheme('RetroGlowAmber'))
  addMenuItem(prefsMenu, 'Theme: Retro No Glow', () => setTheme('RetroPlain'))
  addMenuItem(prefsMenu, 'Theme: Simple', () => setTheme('Simple'))

  // Info Menu
  //prettier-ignore
  addMenuItem(infoMenu, 'System Info', () => printInfo(), true)
  addMenuItem(infoMenu, 'About', () => printAbout())
  addMenuItem(infoMenu, 'Help', () => printHelp())

  // Close menus when clicking/touching the main content area
  const column = document.querySelector('.column')
  column.addEventListener('click', () => hideMenus())
  column.addEventListener('touchstart', () => hideMenus(), { passive: true })
}

// Show a specific menu and hide others
function showMenu(menuDiv) {
  hideMenus()
  menuDiv.style.display = 'block'

  // Disable all with class disableNoFile if no file loaded
  const noFileLoaded = typeof bridge.load !== 'function'
  const disableItems = document.querySelectorAll('.disableNoFile')
  disableItems.forEach((item) => {
    item.style.pointerEvents = noFileLoaded ? 'none' : 'auto'
    item.style.color = noFileLoaded ? '#888888' : '#000000'
  })
}

// Hide all menus
function hideMenus() {
  fileMenu.style.display = 'none'
  sysMenu.style.display = 'none'
  prefsMenu.style.display = 'none'
  infoMenu.style.display = 'none'
}

// Add a menu item to a menu
function addMenuItem(menu, label, onClick, needsFile = false) {
  const item = document.createElement('div')
  item.className = 'menuItem'
  item.textContent = label
  if (needsFile) {
    item.classList.add('disableNoFile')
  }

  const handler = (e) => {
    e.preventDefault()
    e.stopPropagation()
    onClick()
    hideMenus()
  }

  // Handle touch for mobile
  item.addEventListener('touchend', handler, { passive: false })
  // Handle click for desktop
  item.addEventListener('click', handler)

  menu.appendChild(item)
}

// Add a separator to a menu
function addMenuSeparator(menu) {
  menu.appendChild(document.createElement('hr'))
}

// Add click/touch handler to a menu button
function addMenuButton(selector, menu) {
  const button = document.querySelector(selector)

  // Handle touch events for mobile - prevent double firing
  button.addEventListener(
    'touchend',
    (e) => {
      e.preventDefault()
      e.stopPropagation()
      toggleMenu(menu)
    },
    { passive: false }
  )

  // Handle click for desktop
  button.addEventListener('click', (e) => {
    e.preventDefault()
    e.stopPropagation()
    toggleMenu(menu)
  })
}

// Toggle a menu open/closed
function toggleMenu(menu) {
  const isVisible = menu.style.display === 'block'
  hideMenus()
  if (!isVisible) {
    showMenu(menu)
  }
}

// Print about information
function printAbout() {
  showModal(
    `GOZM v${version}\n\nA Z-Machine interpreter written in Go, compiled to WebAssembly.
    \n© Ben Coleman 2025\nGithub: https://github.com/benc-uk/gozm\n\nSee the GitHub repository for full license, acknowledgments etc.`
  )
}

// Print help information
function printHelp() {
  showModal(
    `Help:
Use the File menu to open a Z-Machine story file (.z3) or open a supplied adventure file.
Type your commands in the input box and press Enter to send.
Use the System menu to save/load game state.
Use the Preferences menu to change themes.

System commands:
  /quit - Exit the game
  /restart - Restart the game
  /save - Save the game
  /load - Load a saved game`
  )
}

function printInfo() {
  showModal(bridge.getInfo())
}
