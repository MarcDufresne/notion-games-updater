<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-6 text-white">Release Calendar</h1>

    <div v-if="groupedGames.length > 0">
      <div v-for="group in groupedGames" :key="group.monthYear" class="mb-8">
        <h2 class="text-2xl font-semibold mb-4 border-b-2 border-gray-700 pb-2 text-gray-200">
          {{ group.monthYear }}
        </h2>
        <div class="space-y-3">
          <div
            v-for="game in group.games"
            :key="game.id"
            class="flex bg-gray-800 rounded-lg shadow-lg hover:shadow-xl transition-all border border-gray-700 hover:border-gray-600 overflow-hidden"
          >
            <!-- Full height cover art -->
            <div class="flex-shrink-0 self-stretch" style="width: 100px;">
              <img
                v-if="game.cover_url"
                :src="game.cover_url"
                :alt="game.title"
                class="w-full h-full object-cover"
              />
              <div v-else class="w-full h-full bg-gray-700 flex items-center justify-center border-r border-gray-600">
                <span class="text-gray-500 text-xs">No Cover</span>
              </div>
            </div>

            <!-- Content area -->
            <div class="flex-1 p-3 flex items-center gap-4">
              <div class="flex-1">
                <div class="font-semibold text-white">{{ game.title }}</div>
                <div v-if="game.release_date" class="text-sm text-gray-400">
                  {{ formatDate(game.release_date) }}
                </div>
                <div v-if="game.platforms && game.platforms.length" class="mt-1 relative" style="height: 1.5rem;">
                  <div class="absolute inset-0 overflow-hidden">
                    <div class="flex gap-1 flex-nowrap">
                      <span
                        v-for="platform in sortPlatforms(game.platforms)"
                        :key="platform"
                        :class="[
                          'inline-block text-xs px-2 py-0.5 rounded border flex-shrink-0 whitespace-nowrap',
                          getPlatformColor(platform).bg,
                          getPlatformColor(platform).text,
                          getPlatformColor(platform).border
                        ]"
                      >
                        {{ platform }}
                      </span>
                    </div>
                  </div>
                </div>
              </div>

              <div class="flex-shrink-0">
                <StatusPicker
                  :model-value="game.status"
                  @update:model-value="(newStatus) => handleStatusUpdate(game.id, newStatus)"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div v-else class="text-gray-400 text-center py-8 bg-gray-800/50 rounded-lg border border-gray-700">
      No upcoming releases in your backlog.
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useGamesStore } from '../stores/games'
import { getPlatformColor, sortPlatforms } from '../lib/platformColors'
import StatusPicker from '../components/StatusPicker.vue'

const gamesStore = useGamesStore()

const groupedGames = computed(() => {
  // Get calendar games (backend filters to last month onwards)
  const gamesWithDates = gamesStore.backlog.filter(g => g.release_date)

  // Group by month/year
  const groups = {}

  gamesWithDates.forEach(game => {
    const date = new Date(game.release_date)
    const monthYear = date.toLocaleDateString('en-US', { year: 'numeric', month: 'long' })

    if (!groups[monthYear]) {
      groups[monthYear] = {
        monthYear,
        sortKey: date.getTime(),
        games: []
      }
    }
    groups[monthYear].games.push(game)
  })

  // Sort by date and then sort games within each group
  return Object.values(groups)
    .sort((a, b) => a.sortKey - b.sortKey)
    .map(group => ({
      ...group,
      games: group.games.sort((a, b) => new Date(a.release_date) - new Date(b.release_date))
    }))
})

onMounted(() => {
  gamesStore.fetchGames('calendar')
})

async function handleStatusUpdate(gameId, newStatus) {
  await gamesStore.updateStatus(gameId, newStatus)
}

function formatDate(dateString) {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })
}
</script>

<style scoped>
.container {
  max-width: 800px;
}

h1 {
  color: #fff;
}

h2 {
  color: #e2e8f0;
}

.bg-gray-800 {
  background-color: #2d3748;
}

.bg-gray-700 {
  background-color: #4a5568;
}

.text-gray-400 {
  color: #cbd5e0;
}

.text-gray-500 {
  color: #edf2f7;
}

.border-gray-600 {
  border-color: #718096;
}

.shadow-lg {
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.1);
}

.hover\:shadow-xl:hover {
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.1);
}

.focus\:ring-2:focus {
  ring-width: 2px;
}

.focus\:ring-blue-500:focus {
  --tw-ring-color: rgb(37 99 235 / 50%);
}
</style>
