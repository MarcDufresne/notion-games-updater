<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-6 text-white">History</h1>

    <div v-if="groupedGames.length > 0">
      <div v-for="group in groupedGames" :key="group.year" class="mb-8">
        <h2 class="text-2xl font-semibold mb-4 border-b-2 border-gray-700 pb-2 text-gray-200">
          {{ group.year || 'No Date' }}
        </h2>
        <TransitionGroup name="game-list" tag="div" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
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
      No games in history.
    </div>

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
import { computed, onMounted } from 'vue'
import { useGamesStore } from '../stores/games'
import GameCard from '../components/GameCard.vue'
import GameDetailsModal from '../components/GameDetailsModal.vue'
import { isDatePlayedSentinel } from '../lib/dateUtils'
import { useGameModal } from '../composables/useGameModal'

const gamesStore = useGamesStore()
const { isModalOpen, selectedGame, openModal, closeModal, handleStatusUpdate, handleDeleteGame, handleMatchUpdated } = useGameModal('history')

const groupedGames = computed(() => {
  const groups = {}

  gamesStore.history.forEach(game => {
    let year = 'No Date'
    if (game.date_played && !isDatePlayedSentinel(game.date_played)) {
      const date = new Date(game.date_played)
      year = date.getUTCFullYear()
    }

    if (!groups[year]) {
      groups[year] = []
    }
    groups[year].push(game)
  })

  return Object.keys(groups)
    .sort((a, b) => {
      if (a === 'No Date') return 1
      if (b === 'No Date') return -1
      return b - a
    })
    .map(year => ({
      year,
      games: groups[year].sort((a, b) => {
        // Sort games within each year by date_played descending
        const dateA = a.date_played && !isDatePlayedSentinel(a.date_played)
          ? new Date(a.date_played).getTime()
          : 0
        const dateB = b.date_played && !isDatePlayedSentinel(b.date_played)
          ? new Date(b.date_played).getTime()
          : 0
        return dateB - dateA
      })
    }))
})

onMounted(() => {
  gamesStore.fetchGames('history')
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
</style>
