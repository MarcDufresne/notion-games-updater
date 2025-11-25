<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-6 text-white">History</h1>

    <div v-if="groupedGames.length > 0">
      <div v-for="group in groupedGames" :key="group.year" class="mb-8">
        <h2 class="text-2xl font-semibold mb-4 border-b-2 border-gray-700 pb-2 text-gray-200">
          {{ group.year || 'No Date' }}
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <GameCard
            v-for="game in group.games"
            :key="game.id"
            :game="game"
            @update-status="handleStatusUpdate"
          />
        </div>
      </div>
    </div>
    <div v-else class="text-gray-400 text-center py-8 bg-gray-800/50 rounded-lg border border-gray-700">
      No games in history.
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useGamesStore } from '../stores/games'
import GameCard from '../components/GameCard.vue'

const gamesStore = useGamesStore()

const groupedGames = computed(() => {
  const groups = {}

  gamesStore.history.forEach(game => {
    let year = 'No Date'
    if (game.date_played) {
      const date = new Date(game.date_played)
      year = date.getFullYear()
    }

    if (!groups[year]) {
      groups[year] = []
    }
    groups[year].push(game)
  })

  return Object.keys(groups)
    .sort((a, b) => b - a)
    .map(year => ({
      year,
      games: groups[year]
    }))
})

onMounted(() => {
  gamesStore.fetchGames('history')
})

async function handleStatusUpdate(gameId, newStatus) {
  await gamesStore.updateStatus(gameId, newStatus)
}
</script>
