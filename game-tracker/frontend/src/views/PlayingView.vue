<template>
  <div>
    <h1 class="text-2xl font-bold text-gray-900 mb-6">Currently Playing</h1>
    
    <div v-if="loading" class="text-center py-12">
      <p class="text-gray-500">Loading games...</p>
    </div>
    
    <div v-else-if="error" class="text-center py-12">
      <p class="text-red-600">{{ error }}</p>
    </div>
    
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="game in gamesStore.playing" :key="game.id" class="bg-white rounded-lg shadow-lg p-6">
        <div v-if="game.cover_url" class="mb-4">
          <img :src="game.cover_url" :alt="game.title" class="w-full h-48 object-cover rounded" />
        </div>
        <h3 class="text-lg font-semibold text-gray-900 mb-2">{{ game.title }}</h3>
        <p class="text-sm text-gray-600 mb-1">Rating: {{ game.rating }}</p>
        <p class="text-sm text-gray-600 mb-4">
          Genres: {{ game.genres?.join(', ') || 'N/A' }}
        </p>
        <div class="flex space-x-2">
          <button
            @click="markAsDone(game.id)"
            class="flex-1 py-2 px-4 bg-green-600 text-white text-sm rounded hover:bg-green-700"
          >
            Done
          </button>
          <button
            @click="markAsAbandoned(game.id)"
            class="flex-1 py-2 px-4 bg-gray-600 text-white text-sm rounded hover:bg-gray-700"
          >
            Abandon
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useGamesStore } from '@/stores/games'

const gamesStore = useGamesStore()
const loading = ref(false)
const error = ref(null)

onMounted(async () => {
  loading.value = true
  try {
    await gamesStore.fetchGames('playing')
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
})

async function markAsDone(gameId) {
  try {
    await gamesStore.updateGameStatus(gameId, 'Done', new Date().toISOString())
  } catch (err) {
    error.value = err.message
  }
}

async function markAsAbandoned(gameId) {
  try {
    await gamesStore.updateGameStatus(gameId, 'Abandoned', new Date().toISOString())
  } catch (err) {
    error.value = err.message
  }
}
</script>
