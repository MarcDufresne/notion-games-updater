<template>
  <div id="app" class="min-h-screen bg-gray-900">
    <!-- Auth Gate -->
    <div v-if="!user" class="flex items-center justify-center min-h-screen">
      <div class="bg-gray-800 p-8 rounded-xl shadow-2xl text-center border border-gray-700">
        <div class="mb-6">
          <svg class="w-16 h-16 mx-auto text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 4a2 2 0 114 0v1a1 1 0 001 1h3a1 1 0 011 1v3a1 1 0 01-1 1h-1a2 2 0 100 4h1a1 1 0 011 1v3a1 1 0 01-1 1h-3a1 1 0 01-1-1v-1a2 2 0 10-4 0v1a1 1 0 01-1 1H7a1 1 0 01-1-1v-3a1 1 0 00-1-1H4a2 2 0 110-4h1a1 1 0 001-1V7a1 1 0 011-1h3a1 1 0 001-1V4z"/>
          </svg>
        </div>
        <h1 class="text-3xl font-bold mb-3 text-white">Game Tracker</h1>
        <p class="text-gray-400 mb-6">Sign in to manage your game library</p>
        <button
          @click="signInWithGoogle"
          class="bg-blue-600 text-white px-8 py-3 rounded-lg hover:bg-blue-700 transition-all duration-200 font-medium shadow-lg hover:shadow-xl transform hover:scale-105"
        >
          <span class="flex items-center gap-2 justify-center">
            <svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
              <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
              <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
              <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
              <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
            </svg>
            Sign in with Google
          </span>
        </button>
      </div>
    </div>

    <!-- Main App -->
    <div v-else>
      <!-- Header -->
      <header class="bg-gray-800 border-b border-gray-700 shadow-lg">
        <div class="container mx-auto px-4 py-4">
          <!-- Desktop Layout: Single Row -->
          <div class="hidden lg:flex items-center justify-between gap-6">
            <!-- Left: App Title -->
            <h1 class="text-2xl font-bold whitespace-nowrap text-white flex items-center gap-2">
              <svg class="w-8 h-8 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 4a2 2 0 114 0v1a1 1 0 001 1h3a1 1 0 011 1v3a1 1 0 01-1 1h-1a2 2 0 100 4h1a1 1 0 011 1v3a1 1 0 01-1 1h-3a1 1 0 01-1-1v-1a2 2 0 10-4 0v1a1 1 0 01-1 1H7a1 1 0 01-1-1v-3a1 1 0 00-1-1H4a2 2 0 110-4h1a1 1 0 001-1V7a1 1 0 011-1h3a1 1 0 001-1V4z"/>
              </svg>
              Game Tracker
            </h1>

            <!-- Center: Search Bar -->
            <div class="flex-1 max-w-2xl">
              <GameSearch @game-created="handleGameCreated" @open-existing-game="handleOpenExistingGame" />
            </div>

            <!-- Right: User Info -->
            <div class="flex items-center gap-4 whitespace-nowrap">
              <span class="text-sm text-gray-400">{{ user.email }}</span>
              <button
                @click="signOut"
                class="text-sm text-red-400 hover:text-red-300 transition-colors font-medium"
              >
                Sign Out
              </button>
            </div>
          </div>

          <!-- Mobile Layout: Stacked -->
          <div class="lg:hidden space-y-4">
            <!-- Top Row: Title and User Info -->
            <div class="flex items-center justify-between">
              <h1 class="text-xl font-bold text-white flex items-center gap-2">
                <svg class="w-6 h-6 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 4a2 2 0 114 0v1a1 1 0 001 1h3a1 1 0 011 1v3a1 1 0 01-1 1h-1a2 2 0 100 4h1a1 1 0 011 1v3a1 1 0 01-1 1h-3a1 1 0 01-1-1v-1a2 2 0 10-4 0v1a1 1 0 01-1 1H7a1 1 0 01-1-1v-3a1 1 0 00-1-1H4a2 2 0 110-4h1a1 1 0 001-1V7a1 1 0 011-1h3a1 1 0 001-1V4z"/>
                </svg>
                <span class="hidden sm:inline">Game Tracker</span>
              </h1>

              <button
                @click="signOut"
                class="text-sm text-red-400 hover:text-red-300 transition-colors font-medium"
              >
                Sign Out
              </button>
            </div>

            <!-- Bottom Row: Full-width Search Bar -->
            <div class="w-full">
              <GameSearch @game-created="handleGameCreated" @open-existing-game="handleOpenExistingGame" />
            </div>
          </div>
        </div>
      </header>

      <!-- Navigation Tabs -->
      <nav class="bg-gray-800 border-b border-gray-700 mb-6">
        <div class="container mx-auto px-4">
          <div class="flex gap-1 sm:gap-2 overflow-x-auto scrollbar-hide">
            <button
              v-for="view in views"
              :key="view.id"
              @click="currentView = view.id"
              :class="[
                'px-3 sm:px-6 py-2 sm:py-3 font-medium transition-all duration-200 border-b-2 whitespace-nowrap text-sm sm:text-base flex-shrink-0',
                currentView === view.id
                  ? 'border-blue-500 text-blue-400 bg-gray-700/50'
                  : 'border-transparent text-gray-400 hover:text-gray-200 hover:bg-gray-700/30'
              ]"
            >
              {{ view.name }}
            </button>
          </div>
        </div>
      </nav>

      <!-- View Content -->
      <main>
        <BacklogView v-if="currentView === 'backlog'" />
        <PlayingView v-else-if="currentView === 'playing'" />
        <HistoryView v-else-if="currentView === 'history'" />
        <CalendarView v-else-if="currentView === 'calendar'" />
        <AllView v-else-if="currentView === 'all'" />
      </main>

      <!-- Sticky Search Bar (appears on scroll up) -->
      <StickySearchBar @game-created="handleGameCreated" @open-existing-game="handleOpenExistingGame" />
    </div>

    <!-- Toast Notifications -->
    <Toast ref="toastRef" />

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
import { ref, onMounted } from 'vue'
import { GoogleAuthProvider, signInWithPopup, signOut as firebaseSignOut, onAuthStateChanged } from 'firebase/auth'
import { auth } from './lib/firebase'
import Toast from './components/Toast.vue'
import GameSearch from './components/GameSearch.vue'
import StickySearchBar from './components/StickySearchBar.vue'
import GameDetailsModal from './components/GameDetailsModal.vue'
import BacklogView from './views/BacklogView.vue'
import PlayingView from './views/PlayingView.vue'
import HistoryView from './views/HistoryView.vue'
import CalendarView from './views/CalendarView.vue'
import AllView from './views/AllView.vue'
import { useGameModal } from './composables/useGameModal'
import { useGamesStore } from './stores/games'

const user = ref(null)
const currentView = ref('backlog')
const toastRef = ref(null)

// Initialize game modal
const { isModalOpen, selectedGame, openModal, closeModal, handleStatusUpdate, handleDeleteGame, handleMatchUpdated } = useGameModal()

// Initialize games store
const gamesStore = useGamesStore()

// Make toast available globally
if (typeof window !== 'undefined') {
  window.$toast = {
    success: (message, duration) => toastRef.value?.addToast(message, 'success', duration),
    error: (message, duration) => toastRef.value?.addToast(message, 'error', duration),
    warning: (message, duration) => toastRef.value?.addToast(message, 'warning', duration),
    info: (message, duration) => toastRef.value?.addToast(message, 'info', duration)
  }
}

const views = [
  { id: 'backlog', name: 'Backlog' },
  { id: 'playing', name: 'Playing' },
  { id: 'history', name: 'History' },
  { id: 'calendar', name: 'Calendar' },
  { id: 'all', name: 'All' }
]

onMounted(() => {
  onAuthStateChanged(auth, (firebaseUser) => {
    user.value = firebaseUser
    if (firebaseUser) {
      // Load all games for library checking
      gamesStore.fetchGames('all')
    }
  })
})

async function signInWithGoogle() {
  const provider = new GoogleAuthProvider()
  try {
    await signInWithPopup(auth, provider)
  } catch (error) {
    console.error('Sign in error:', error)
    alert('Failed to sign in')
  }
}

async function signOut() {
  try {
    await firebaseSignOut(auth)
  } catch (error) {
    console.error('Sign out error:', error)
  }
}

function handleGameCreated() {
  // Optionally show a success message or refresh views
  console.log('Game created successfully')
}

function handleOpenExistingGame(game) {
  // Open the modal with the existing game
  openModal(game)
}
</script>
