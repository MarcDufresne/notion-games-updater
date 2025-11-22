<template>
  <div>
    <h1 class="text-2xl font-bold text-gray-900 mb-6">Gaming History</h1>
    
    <div v-if="loading" class="text-center py-12">
      <p class="text-gray-500">Loading games...</p>
    </div>
    
    <div v-else-if="error" class="text-center py-12">
      <p class="text-red-600">{{ error }}</p>
    </div>
    
    <div v-else>
      <div v-for="(games, year) in gamesByYear" :key="year" class="mb-8">
        <h2 class="text-xl font-semibold text-gray-800 mb-4">{{ year }}</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <div v-for="game in games" :key="game.id" class="bg-white rounded-lg shadow p-4">
            <div v-if="game.cover_url" class="mb-2">
              <img :src="game.cover_url" :alt="game.title" class="w-full h-32 object-cover rounded" />
            </div>
            <h3 class="font-semibold text-gray-900 text-sm">{{ game.title }}</h3>
            <p class="text-xs text-gray-600">{{ game.status }}</p>
            <p class="text-xs text-gray-500">
              {{ formatDate(game.date_played) }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import { useGamesStore } from '@/stores/games'

const gamesStore = useGamesStore()
const loading = ref(false)
const error = ref(null)

const gamesByYear = computed(() => {
  const grouped = {}
  
  for (const game of gamesStore.history) {
    if (game.date_played) {
      const year = new Date(game.date_played).getFullYear()
      if (!grouped[year]) {
        grouped[year] = []
      }
      grouped[year].push(game)
    }
  }
  
  // Sort years descending
  const sortedYears = Object.keys(grouped).sort((a, b) => b - a)
  const result = {}
  for (const year of sortedYears) {
    result[year] = grouped[year]
  }
  
  return result
})

onMounted(async () => {
  loading.value = true
  try {
    await gamesStore.fetchGames('history')
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
})

function formatDate(dateString) {
  if (!dateString) return 'N/A'
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}
</script>
