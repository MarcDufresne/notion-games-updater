<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-6 text-white">Currently Playing</h1>

    <div v-if="gamesStore.playing.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <GameCard
        v-for="game in gamesStore.playing"
        :key="game.id"
        :game="game"
        @update-status="handleStatusUpdate"
      />
    </div>
    <div v-else class="text-gray-400 text-center py-8 bg-gray-800/50 rounded-lg border border-gray-700">
      No games currently being played.
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useGamesStore } from '../stores/games'
import GameCard from '../components/GameCard.vue'

const gamesStore = useGamesStore()

onMounted(() => {
  gamesStore.fetchGames('playing')
})

async function handleStatusUpdate(gameId, newStatus) {
  await gamesStore.updateStatus(gameId, newStatus)
}
</script>
