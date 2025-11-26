// Platform category definitions
const PLATFORM_CATEGORIES = {
  PC: {
    priority: 1,
    bg: 'bg-gray-700',
    text: 'text-gray-300',
    border: 'border-gray-600',
    matcher: (platformLower) =>
      platformLower === 'pc' ||
      platformLower.includes('windows') ||
      platformLower.includes('linux') ||
      platformLower.includes('mac') ||
      platformLower === 'steam'
  },
  SONY: {
    priority: 2,
    bg: 'bg-blue-900/50',
    text: 'text-blue-300',
    border: 'border-blue-700/50',
    matcher: (platformLower) =>
      platformLower.includes('playstation') ||
      platformLower.includes('ps') ||
      platformLower === 'ps5' ||
      platformLower === 'ps4' ||
      platformLower === 'ps3' ||
      platformLower === 'ps2' ||
      platformLower === 'psvita' ||
      platformLower === 'psp'
  },
  NINTENDO: {
    priority: 3,
    bg: 'bg-red-900/50',
    text: 'text-red-300',
    border: 'border-red-700/50',
    matcher: (platformLower) =>
      platformLower.includes('nintendo') ||
      platformLower.includes('switch') ||
      platformLower.includes('wii') ||
      platformLower.includes('gamecube') ||
      platformLower === 'nes' ||
      platformLower === 'snes' ||
      platformLower.includes('3ds') ||
      platformLower.includes('ds')
  },
  XBOX: {
    priority: 4,
    bg: 'bg-green-900/50',
    text: 'text-green-300',
    border: 'border-green-700/50',
    matcher: (platformLower) =>
      platformLower.includes('xbox') ||
      platformLower.includes('xb') ||
      platformLower === 'xone' ||
      platformLower === 'series x|s' ||
      platformLower === 'series s' ||
      platformLower === 'series x'
  },
  OTHER: {
    priority: 5,
    bg: 'bg-purple-900/50',
    text: 'text-purple-300',
    border: 'border-purple-700/50',
    matcher: () => true // Catch-all
  }
}

// Get platform category info
function getPlatformInfo(platform) {
  const platformLower = platform.toLowerCase()

  for (const [name, info] of Object.entries(PLATFORM_CATEGORIES)) {
    if (info.matcher(platformLower)) {
      return { name, ...info }
    }
  }

  return { name: 'OTHER', ...PLATFORM_CATEGORIES.OTHER }
}

// Platform color coding utility
export function getPlatformColor(platform) {
  const info = getPlatformInfo(platform)
  return {
    bg: info.bg,
    text: info.text,
    border: info.border
  }
}

// Sort platforms: PC, Sony, Nintendo, Xbox, Other
export function sortPlatforms(platforms) {
  if (!platforms || !Array.isArray(platforms)) {
    return []
  }

  return [...platforms].sort((a, b) => {
    const infoA = getPlatformInfo(a)
    const infoB = getPlatformInfo(b)

    if (infoA.priority !== infoB.priority) {
      return infoA.priority - infoB.priority
    }

    // Within same category, sort alphabetically
    return a.localeCompare(b)
  })
}
