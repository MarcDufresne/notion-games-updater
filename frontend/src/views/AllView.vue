<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-6 text-white">All Games</h1>

    <!-- Controls Bar -->
    <div class="mb-6 bg-gray-800/50 rounded-lg border border-gray-700 p-4">
      <div class="flex flex-col lg:flex-row gap-6 items-start lg:items-center justify-between">
        <div class="flex flex-col sm:flex-row gap-6 flex-1">
          <!-- Sort Mode Button Group -->
          <div class="flex flex-col gap-2">
            <label class="text-gray-400 text-xs font-medium uppercase tracking-wider">Sort by</label>
            <div class="inline-flex rounded-lg border border-gray-600 overflow-hidden">
              <button
                @click="sortMode = 'release'"
                :class="[
                  'flex items-center gap-2 px-4 py-2 text-sm font-medium transition-all',
                  sortMode === 'release'
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-700 text-gray-300 hover:bg-gray-600 hover:text-white'
                ]"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                <span>Release</span>
              </button>
              <button
                @click="sortMode = 'name'"
                :class="[
                  'flex items-center gap-2 px-4 py-2 text-sm font-medium border-l border-gray-600 transition-all',
                  sortMode === 'name'
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-700 text-gray-300 hover:bg-gray-600 hover:text-white'
                ]"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12" />
                </svg>
                <span>Name</span>
              </button>
            </div>
          </div>

          <!-- Match Filter Button Group -->
          <div class="flex flex-col gap-2">
            <label class="text-gray-400 text-xs font-medium uppercase tracking-wider">Filter</label>
            <div class="inline-flex rounded-lg border border-gray-600 overflow-hidden">
              <button
                @click="matchFilter = 'all'"
                :class="[
                  'flex items-center gap-2 px-4 py-2 text-sm font-medium transition-all',
                  matchFilter === 'all'
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-700 text-gray-300 hover:bg-gray-600 hover:text-white'
                ]"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
                <span>All</span>
              </button>
              <button
                @click="matchFilter = 'matched'"
                :class="[
                  'flex items-center gap-2 px-4 py-2 text-sm font-medium border-l border-gray-600 transition-all',
                  matchFilter === 'matched'
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-700 text-gray-300 hover:bg-gray-600 hover:text-white'
                ]"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <span>Matched</span>
              </button>
              <button
                @click="matchFilter = 'unmatched'"
                :class="[
                  'flex items-center gap-2 px-4 py-2 text-sm font-medium border-l border-gray-600 transition-all',
                  matchFilter === 'unmatched'
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-700 text-gray-300 hover:bg-gray-600 hover:text-white'
                ]"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
                <span>Issues</span>
              </button>
            </div>
          </div>
        </div>

        <!-- Count Badge -->
        <div class="text-sm">
          <div class="text-gray-400 mb-1 text-xs font-medium uppercase tracking-wider">Results</div>
          <div class="text-white font-semibold">
            {{ filteredGames.length }} <span class="text-gray-500">/</span> {{ gamesStore.all.length }}
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content with Sidebar Layout -->
    <div class="flex gap-6 relative">
      <!-- Games List -->
      <div class="flex-1 min-w-0">
        <div v-if="groupedGames.length > 0">
          <div v-for="group in groupedGames" :key="group.section" :id="group.sectionId" class="mb-8 scroll-mt-24">
            <h2 class="text-2xl font-semibold mb-4 border-b-2 border-gray-700 pb-2 text-gray-200">
              {{ group.section }}
            </h2>
            <TransitionGroup name="game-list" tag="div" class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
              <GameCard
                v-for="game in group.games"
                :key="game.id"
                :game="game"
                @update-status="handleStatusUpdate"
                @card-click="openModal"
              />
            </TransitionGroup>
          </div>
        </div>
        <div v-else class="text-gray-400 text-center py-8 bg-gray-800/50 rounded-lg border border-gray-700">
          {{ emptyStateMessage }}
        </div>
      </div>

      <!-- Sidebar Navigation (Desktop) -->
      <div v-if="groupedGames.length > 0" class="hidden lg:block w-40 flex-shrink-0">
        <div class="sticky top-24 bg-gray-800 border border-gray-700 rounded-lg shadow-xl max-h-[calc(100vh-8rem)] overflow-y-auto">
          <div class="py-3">
            <nav>
              <template v-for="(group, index) in groupedGames" :key="group.section">
                <a
                  :href="`#${group.sectionId}`"
                  @click.prevent="scrollToSection(group.sectionId)"
                  :class="[
                    'block px-4 py-2.5 text-center transition-colors',
                    activeSectionId === group.sectionId
                      ? 'bg-blue-600 text-white font-medium'
                      : 'text-gray-300 hover:bg-gray-700 hover:text-white'
                  ]"
                >
                  {{ group.section }}
                </a>
                <div v-if="index < groupedGames.length - 1" class="border-t border-gray-700 mx-4"></div>
              </template>
            </nav>
          </div>
        </div>
      </div>
    </div>

    <!-- Sidebar Toggle Button (Mobile/Tablet only) -->
    <button
      v-if="groupedGames.length > 0"
      @click="isSidebarVisible = !isSidebarVisible"
      class="fixed bottom-6 right-6 z-50 bg-blue-600 hover:bg-blue-700 text-white rounded-full p-3 shadow-lg transition-all lg:hidden"
      :title="isSidebarVisible ? 'Hide sections' : 'Show sections'"
    >
      <svg v-if="!isSidebarVisible" class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
      </svg>
      <svg v-else class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
      </svg>
    </button>

    <!-- Sidebar Navigation (Mobile Overlay) -->
    <Transition name="sidebar">
      <div
        v-if="isSidebarVisible && groupedGames.length > 0"
        class="fixed top-24 right-4 z-40 bg-gray-800 border border-gray-700 rounded-lg shadow-xl max-h-[calc(100vh-8rem)] overflow-y-auto w-40 lg:hidden"
      >
        <div class="py-3">
          <nav>
            <template v-for="(group, index) in groupedGames" :key="group.section">
              <a
                :href="`#${group.sectionId}`"
                @click.prevent="scrollToSection(group.sectionId)"
                :class="[
                  'block px-4 py-2.5 text-center transition-colors',
                  activeSectionId === group.sectionId
                    ? 'bg-blue-600 text-white font-medium'
                    : 'text-gray-300 hover:bg-gray-700 hover:text-white'
                ]"
              >
                {{ group.section }}
              </a>
              <div v-if="index < groupedGames.length - 1" class="border-t border-gray-700 mx-4"></div>
            </template>
          </nav>
        </div>
      </div>
    </Transition>

    <!-- Sidebar Backdrop (Mobile/Tablet only) -->
    <Transition name="backdrop">
      <div
        v-if="isSidebarVisible && groupedGames.length > 0"
        @click="isSidebarVisible = false"
        class="fixed inset-0 bg-black/50 z-30 lg:hidden"
      ></div>
    </Transition>

    <!-- Game Details Modal -->
    <GameDetailsModal
      :is-open="isModalOpen"
      :game="selectedGame"
      @close="closeModal"
      @update-status="handleStatusUpdate"
      @delete-game="handleDeleteGame"
      @match-updated="handleMatchUpdated"
    />
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useGamesStore } from '../stores/games'
import GameCard from '../components/GameCard.vue'
import GameDetailsModal from '../components/GameDetailsModal.vue'
import { isReleaseDateSentinel } from '../lib/dateUtils'
import { useGameModal } from '../composables/useGameModal'

const gamesStore = useGamesStore()
const { isModalOpen, selectedGame, openModal, closeModal, handleStatusUpdate, handleDeleteGame, handleMatchUpdated } = useGameModal('all')

// Reactive state for controls
const sortMode = ref('release')
const matchFilter = ref('all')
const isSidebarVisible = ref(typeof window !== 'undefined' ? window.innerWidth >= 1024 : false)
const activeSectionId = ref('')

// Filter games based on match status
const filteredGames = computed(() => {
  const allGames = gamesStore.all

  if (matchFilter.value === 'matched') {
    return allGames.filter(game => game.igdb_id && game.igdb_id > 0)
  } else if (matchFilter.value === 'unmatched') {
    return allGames.filter(game => {
      // Include games without IGDB ID or with problematic match status
      const hasNoIgdbId = !game.igdb_id || game.igdb_id === 0
      const hasMatchIssues = ['multiple', 'no_match', 'needs_review'].includes(game.match_status)
      return hasNoIgdbId || hasMatchIssues
    })
  }

  return allGames
})

// Group and sort games based on current mode
const groupedGames = computed(() => {
  const games = filteredGames.value
  const groups = {}

  if (sortMode.value === 'name') {
    // Group by first letter (A-Z)
    games.forEach(game => {
      const firstLetter = game.title.charAt(0).toUpperCase()
      const letter = /[A-Z]/.test(firstLetter) ? firstLetter : '#'

      if (!groups[letter]) {
        groups[letter] = []
      }
      groups[letter].push(game)
    })

    // Sort groups alphabetically, with # at the end
    return Object.keys(groups)
      .sort((a, b) => {
        if (a === '#') return 1
        if (b === '#') return -1
        return a.localeCompare(b)
      })
      .map(letter => ({
        section: letter,
        sectionId: `section-${letter.toLowerCase()}`,
        games: groups[letter].sort((a, b) => a.title.localeCompare(b.title))
      }))
  } else {
    // Group by release year (existing logic)
    games.forEach(game => {
      let year = 'TBD'
      if (game.release_date && !isReleaseDateSentinel(game.release_date)) {
        const date = new Date(game.release_date)
        year = date.getUTCFullYear()
      }

      if (!groups[year]) {
        groups[year] = []
      }
      groups[year].push(game)
    })

    return Object.keys(groups)
      .sort((a, b) => {
        if (a === 'TBD') return 1
        if (b === 'TBD') return -1
        return b - a
      })
      .map(year => ({
        section: year,
        sectionId: `section-${year.toString().toLowerCase()}`,
        games: groups[year].sort((a, b) => {
          const aIsTBD = !a.release_date || isReleaseDateSentinel(a.release_date)
          const bIsTBD = !b.release_date || isReleaseDateSentinel(b.release_date)

          if (aIsTBD && bIsTBD) return a.title.localeCompare(b.title)
          if (aIsTBD) return 1
          if (bIsTBD) return -1

          return new Date(b.release_date) - new Date(a.release_date)
        })
      }))
  }
})

// Empty state message based on filter
const emptyStateMessage = computed(() => {
  if (matchFilter.value === 'matched') {
    return 'No matched games found.'
  } else if (matchFilter.value === 'unmatched') {
    return 'No unmatched or problematic games found.'
  }
  return 'No games in your library.'
})

// Scroll to section smoothly
function scrollToSection(sectionId) {
  const element = document.getElementById(sectionId)
  if (element) {
    element.scrollIntoView({ behavior: 'smooth', block: 'start' })
    // Hide sidebar on mobile after clicking
    if (window.innerWidth < 1024) {
      isSidebarVisible.value = false
    }
  }
}

// Track active section with IntersectionObserver
let observer = null

function setupObserver() {
  // Disconnect existing observer if any
  if (observer) {
    observer.disconnect()
  }

  // Set up intersection observer for section tracking
  observer = new IntersectionObserver(
    (entries) => {
      // Find the entry with the highest intersection ratio that's visible
      const visibleEntries = entries.filter(entry => entry.isIntersecting)
      if (visibleEntries.length > 0) {
        // Sort by intersection ratio and position, prioritize top sections
        const topEntry = visibleEntries.reduce((prev, curr) => {
          if (curr.boundingClientRect.top < 200 && curr.boundingClientRect.top >= 0) {
            return curr
          }
          if (prev.boundingClientRect.top < 200 && prev.boundingClientRect.top >= 0) {
            return prev
          }
          return curr.intersectionRatio > prev.intersectionRatio ? curr : prev
        })
        activeSectionId.value = topEntry.target.id
      }
    },
    {
      threshold: [0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1],
      rootMargin: '-96px 0px -60% 0px' // Account for sticky header
    }
  )

  // Observe all section elements
  setTimeout(() => {
    groupedGames.value.forEach(group => {
      const element = document.getElementById(group.sectionId)
      if (element) {
        observer.observe(element)
      }
    })
  }, 100)
}

// Watch for changes in groupedGames to re-setup observer
watch(groupedGames, () => {
  setupObserver()
}, { flush: 'post' })

onMounted(() => {
  gamesStore.fetchGames('all')
  setupObserver()
})

onUnmounted(() => {
  if (observer) {
    observer.disconnect()
  }
})
</script>

<style scoped>
/* Smooth transitions for game cards */
.game-list-move,
.game-list-enter-active,
.game-list-leave-active {
  transition: all 0.5s ease;
}

.game-list-enter-from {
  opacity: 0;
  transform: translateY(-30px);
}

.game-list-leave-to {
  opacity: 0;
  transform: translateY(30px);
}

/* Ensure leaving items are positioned absolutely to allow smooth movement */
.game-list-leave-active {
  position: absolute;
}

/* Sidebar transitions */
.sidebar-enter-active,
.sidebar-leave-active {
  transition: all 0.3s ease;
}

.sidebar-enter-from,
.sidebar-leave-to {
  opacity: 0;
  transform: translateX(20px);
}

/* Backdrop transitions */
.backdrop-enter-active,
.backdrop-leave-active {
  transition: opacity 0.3s ease;
}

.backdrop-enter-from,
.backdrop-leave-to {
  opacity: 0;
}
</style>
