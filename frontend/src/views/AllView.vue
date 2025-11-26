<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-6 text-white">All Games</h1>

    <div v-if="groupedGames.length > 0">
      <div v-for="group in groupedGames" :key="group.year" class="mb-8">
        <h2 class="text-2xl font-semibold mb-4 border-b-2 border-gray-700 pb-2 text-gray-200">
          {{ group.year }}
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
      No games in your library.
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
import { isReleaseDateSentinel } from '../lib/dateUtils'
import { useGameModal } from '../composables/useGameModal'

const gamesStore = useGamesStore()
const { isModalOpen, selectedGame, openModal, closeModal, handleStatusUpdate, handleDeleteGame, handleMatchUpdated } = useGameModal('all')

const groupedGames = computed(() => {
  const allGames = gamesStore.all
  const groups = {}

  allGames.forEach(game => {
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
      year,
      games: groups[year].sort((a, b) => {
        const aIsTBD = !a.release_date || isReleaseDateSentinel(a.release_date)
        const bIsTBD = !b.release_date || isReleaseDateSentinel(b.release_date)

        if (aIsTBD && bIsTBD) return a.title.localeCompare(b.title)
        if (aIsTBD) return 1
        if (bIsTBD) return -1

        return new Date(b.release_date) - new Date(a.release_date)
      })
    }))
})

onMounted(() => {
  gamesStore.fetchGames('all')
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
