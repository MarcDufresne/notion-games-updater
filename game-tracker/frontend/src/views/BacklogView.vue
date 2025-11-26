<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-6 text-white">Backlog</h1>

    <!-- Loading State -->
    <div v-if="gamesStore.loading" class="text-center py-8">
      <p class="text-gray-400">Loading games...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="gamesStore.error" class="text-center py-8">
      <p class="text-red-400">Error: {{ gamesStore.error }}</p>
    </div>

    <!-- Games Display -->
    <div v-else>
      <!-- Backlog Statistics -->
      <div class="mb-8">
        <h3 class="text-xl font-semibold mb-4 text-gray-300">Statistics</h3>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <!-- Total Backlog Card -->
          <div class="bg-gradient-to-br from-blue-900/50 to-blue-800/50 rounded-lg p-6 shadow-lg border border-blue-700/50 backdrop-blur-sm">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium text-blue-300 uppercase tracking-wide">Total Backlog</p>
                <p class="text-3xl font-bold text-blue-100 mt-2">{{ gamesStore.backlog.length }}</p>
              </div>
              <div class="bg-blue-700/50 rounded-full p-3">
                <svg class="w-8 h-8 text-blue-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
                </svg>
              </div>
            </div>
          </div>

          <!-- Break Games Card -->
          <div class="bg-gradient-to-br from-amber-900/50 to-amber-800/50 rounded-lg p-6 shadow-lg border border-amber-700/50 backdrop-blur-sm">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium text-amber-300 uppercase tracking-wide">On Break</p>
                <p class="text-3xl font-bold text-amber-100 mt-2">{{ breakGames.length }}</p>
              </div>
              <div class="bg-amber-700/50 rounded-full p-3">
                <svg class="w-8 h-8 text-amber-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
              </div>
            </div>
          </div>

          <!-- Up Next Card -->
          <div class="bg-gradient-to-br from-green-900/50 to-green-800/50 rounded-lg p-6 shadow-lg border border-green-700/50 backdrop-blur-sm">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium text-green-300 uppercase tracking-wide">Up Next</p>
                <p class="text-3xl font-bold text-green-100 mt-2">{{ upNextGames.length }}</p>
              </div>
              <div class="bg-green-700/50 rounded-full p-3">
                <svg class="w-8 h-8 text-green-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6"/>
                </svg>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Break Section -->
      <div v-if="breakGames.length > 0" class="mb-8">
        <h2 class="text-2xl font-semibold mb-4 text-gray-200">Break</h2>
        <TransitionGroup name="game-list" tag="div" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <GameCard
            v-for="game in breakGames"
            :key="game.id"
            :game="game"
            @update-status="handleStatusUpdate"
            @card-click="openModal"
          />
        </TransitionGroup>
      </div>

      <!-- Up Next Section -->
      <div>
        <h2 class="text-2xl font-semibold mb-4 text-gray-200">Up Next</h2>
        <TransitionGroup v-if="upNextGames.length > 0" name="game-list" tag="div" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <GameCard
            v-for="game in upNextGames"
            :key="game.id"
            :game="game"
            @update-status="handleStatusUpdate"
            @card-click="openModal"
          />
        </TransitionGroup>
        <div v-else class="text-gray-400 text-center py-8 bg-gray-800/50 rounded-lg border border-gray-700">
          No games in your backlog. Use the search above to add games!
        </div>
      </div>
    </div>

    <!-- Game Details Modal -->
    <GameDetailsModal
      :is-open="isModalOpen"
      :game="selectedGame"
      @close="closeModal"
      @update-status="handleStatusUpdate"
      @delete-game="handleDeleteGame"
    />
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useGamesStore } from '../stores/games'
import GameCard from '../components/GameCard.vue'
import GameDetailsModal from '../components/GameDetailsModal.vue'

const gamesStore = useGamesStore()
const isModalOpen = ref(false)
const selectedGameId = ref(null)

// Computed to always get fresh game data from store
const selectedGame = computed(() => {
  if (!selectedGameId.value) return {}
  return gamesStore.backlog.find(g => g.id === selectedGameId.value) || {}
})

const breakGames = computed(() => {
  const games = gamesStore.backlog.filter(g => g.status === 'Break')
  console.log('[BacklogView] breakGames computed:', games)
  return games
})

const upNextGames = computed(() => {
  const games = gamesStore.backlog.filter(g => g.status === 'Backlog')
  console.log('[BacklogView] upNextGames computed:', games)
  return games
})

// Watch backlog for changes
watch(() => gamesStore.backlog, (newVal) => {
  console.log('[BacklogView] backlog changed:', newVal)
  console.log('[BacklogView] backlog length:', newVal ? newVal.length : 0)
}, { deep: true })

onMounted(() => {
  console.log('[BacklogView] Component mounted, fetching games...')
  gamesStore.fetchGames('backlog')
})

function openModal(game) {
  selectedGameId.value = game.id
  isModalOpen.value = true
}

function closeModal() {
  isModalOpen.value = false
  selectedGameId.value = null
}

async function handleStatusUpdate(gameId, newStatus, datePlayed = null) {
  await gamesStore.updateStatus(gameId, newStatus, datePlayed)
}

async function handleDeleteGame(gameId) {
  await gamesStore.deleteGame(gameId)
  closeModal()
}
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
</style>
