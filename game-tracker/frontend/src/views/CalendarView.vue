<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-6 text-white">Release Calendar</h1>

    <div v-if="groupedGames.length > 0">
      <div v-for="group in groupedGames" :key="group.monthYear" class="mb-8">
        <h2 class="text-2xl font-semibold mb-4 border-b-2 border-gray-700 pb-2 text-gray-200">
          {{ group.monthYear }}
        </h2>
        <TransitionGroup name="game-list" tag="div" class="space-y-3">
          <div
            v-for="game in group.games"
            :key="game.id"
            class="flex bg-gray-800 rounded-lg shadow-lg hover:shadow-xl transition-all border border-gray-700 hover:border-gray-600 overflow-hidden cursor-pointer"
            @click="openModal(game)"
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
                <div v-if="game.release_date && !group.isTBD" :class="['text-sm', isReleased(game.release_date) ? 'text-gray-300' : 'text-gray-500']">
                  <span v-if="isReleased(game.release_date)">Released: </span>
                  <span v-else>Releases: </span>
                  {{ formatReleaseDate(game.release_date, { year: 'numeric', month: 'short', day: 'numeric' }) }}
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

              <div class="flex-shrink-0" @click.stop>
                <StatusPicker
                  :model-value="game.status"
                  @update:model-value="(newStatus) => handleStatusUpdate(game.id, newStatus)"
                />
              </div>
            </div>
          </div>
        </TransitionGroup>
      </div>
    </div>
    <div v-else class="text-gray-400 text-center py-8 bg-gray-800/50 rounded-lg border border-gray-700">
      No upcoming releases in your backlog.
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
import { computed, onMounted, ref } from 'vue'
import { useGamesStore } from '../stores/games'
import { getPlatformColor, sortPlatforms } from '../lib/platformColors'
import { formatReleaseDate, isReleaseDateSentinel } from '../lib/dateUtils'
import StatusPicker from '../components/StatusPicker.vue'
import GameDetailsModal from '../components/GameDetailsModal.vue'

const gamesStore = useGamesStore()
const isModalOpen = ref(false)
const selectedGameId = ref(null)

// Computed to always get fresh game data from store
const selectedGame = computed(() => {
  if (!selectedGameId.value) return {}
  return gamesStore.calendar.find(g => g.id === selectedGameId.value) || {}
})

const groupedGames = computed(() => {
  // Separate games with actual dates from TBD games
  const gamesWithDates = gamesStore.calendar.filter(g => g.release_date && !isReleaseDateSentinel(g.release_date))
  const tbdGames = gamesStore.calendar.filter(g => !g.release_date || isReleaseDateSentinel(g.release_date))

  // Group games with dates by month/year
  const groups = {}

  gamesWithDates.forEach(game => {
    const date = new Date(game.release_date)
    // Extract UTC date components to avoid timezone issues
    const year = date.getUTCFullYear()
    const month = date.getUTCMonth()
    const day = date.getUTCDate()
    const utcDate = new Date(year, month, day)

    const monthYear = utcDate.toLocaleDateString('en-US', { year: 'numeric', month: 'long' })

    if (!groups[monthYear]) {
      groups[monthYear] = {
        monthYear,
        sortKey: utcDate.getTime(),
        games: [],
        isTBD: false
      }
    }
    groups[monthYear].games.push(game)
  })

  // Sort dated groups by date
  const sortedGroups = Object.values(groups)
    .sort((a, b) => a.sortKey - b.sortKey)
    .map(group => ({
      ...group,
      games: group.games.sort((a, b) => new Date(a.release_date) - new Date(b.release_date))
    }))

  // Add TBD games as a separate group at the end if there are any
  if (tbdGames.length > 0) {
    sortedGroups.push({
      monthYear: 'To Be Determined',
      sortKey: Infinity,
      games: tbdGames.sort((a, b) => a.title.localeCompare(b.title)), // Sort TBD games alphabetically
      isTBD: true
    })
  }

  return sortedGroups
})

onMounted(() => {
  gamesStore.fetchGames('calendar')
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

function isReleased(releaseDate) {
  if (!releaseDate) return false
  const date = new Date(releaseDate)
  const year = date.getUTCFullYear()
  const month = date.getUTCMonth()
  const day = date.getUTCDate()
  const releaseDateLocal = new Date(year, month, day)
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return releaseDateLocal <= today
}
</script>

<style scoped>
.container {
  max-width: 800px;
}

/* Smooth transitions for game cards */
.game-list-move,
.game-list-enter-active,
.game-list-leave-active {
  transition: all 0.5s ease;
}

.game-list-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.game-list-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

/* Ensure leaving items are positioned absolutely to allow smooth movement */
.game-list-leave-active {
  position: absolute;
  width: 100%;
}
</style>
