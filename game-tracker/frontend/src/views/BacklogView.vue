<template>
  <div>
    <h1 class="text-2xl font-bold text-gray-900 mb-6">Backlog</h1>
    
    <div v-if="loading" class="text-center py-12">
      <p class="text-gray-500">Loading games...</p>
    </div>
    
    <div v-else-if="error" class="text-center py-12">
      <p class="text-red-600">{{ error }}</p>
    </div>
    
    <div v-else>
      <div v-if="gamesStore.breakGames.length > 0" class="mb-8">
        <h2 class="text-xl font-semibold text-gray-800 mb-4">On Break</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div v-for="game in gamesStore.breakGames" :key="game.id" class="bg-white rounded-lg shadow p-4">
            <h3 class="font-semibold text-gray-900">{{ game.title }}</h3>
            <p class="text-sm text-gray-600">Rating: {{ game.rating }}</p>
            <button
              @click="moveToPlaying(game.id)"
              class="mt-2 text-sm text-blue-600 hover:text-blue-800"
            >
              Start Playing
            </button>
          </div>
        </div>
      </div>
      
      <div>
        <h2 class="text-xl font-semibold text-gray-800 mb-4">Up Next</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div v-for="game in gamesStore.backlogGames" :key="game.id" class="bg-white rounded-lg shadow p-4">
            <h3 class="font-semibold text-gray-900">{{ game.title }}</h3>
            <p class="text-sm text-gray-600">Rating: {{ game.rating }}</p>
            <button
              @click="moveToPlaying(game.id)"
              class="mt-2 text-sm text-blue-600 hover:text-blue-800"
            >
              Start Playing
            </button>
          </div>
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
    await gamesStore.fetchGames('backlog')
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
})

async function moveToPlaying(gameId) {
  try {
    await gamesStore.updateGameStatus(gameId, 'Playing')
  } catch (err) {
    error.value = err.message
  }
}
</script>
