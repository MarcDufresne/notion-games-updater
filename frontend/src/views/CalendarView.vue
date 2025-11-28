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
                <div v-if="game.release_date && !group.isTBD" :class="['text-sm flex items-center gap-1', isReleased(game.release_date) ? 'text-gray-300' : 'text-gray-500']">
                  <!-- Calendar icon for released games -->
                  <svg v-if="isReleased(game.release_date)" class="w-4 h-4 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z" clip-rule="evenodd" />
                  </svg>
                  <!-- Clock icon for upcoming releases -->
                  <svg v-else class="w-4 h-4 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd" />
                  </svg>
                  <span>{{ formatReleaseDate(game.release_date, { year: 'numeric', month: 'short', day: 'numeric' }) }}</span>
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
                  <!-- Gradient fade on right edge -->
                  <div class="absolute right-0 top-0 bottom-0 w-12 bg-gradient-to-l from-gray-800 to-transparent pointer-events-none"></div>
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
      @match-updated="handleMatchUpdated"
    />
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useGamesStore } from '../stores/games'
import { getPlatformColor, sortPlatforms } from '../lib/platformColors'
import { formatReleaseDate, isReleaseDateSentinel } from '../lib/dateUtils'
import StatusPicker from '../components/StatusPicker.vue'
import GameDetailsModal from '../components/GameDetailsModal.vue'
import { useGameModal } from '../composables/useGameModal'

const gamesStore = useGamesStore()
const { isModalOpen, selectedGame, openModal, closeModal, handleStatusUpdate, handleDeleteGame, handleMatchUpdated } = useGameModal('calendar')

const groupedGames = computed(() => {
  // Get today's date at midnight
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  // 30 days ago
  const thirtyDaysAgo = new Date(today)
  thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30)

  // Separate games with actual dates from TBD games
  const gamesWithDates = gamesStore.calendar.filter(g => g.release_date && !isReleaseDateSentinel(g.release_date))
  const tbdGames = gamesStore.calendar.filter(g => !g.release_date || isReleaseDateSentinel(g.release_date))

  // Separate recent releases (released in last 30 days) from upcoming games
  const recentReleases = []
  const upcomingGames = []

  gamesWithDates.forEach(game => {
    const date = new Date(game.release_date)
    const year = date.getUTCFullYear()
    const month = date.getUTCMonth()
    const day = date.getUTCDate()
    const releaseDateLocal = new Date(year, month, day)

    if (releaseDateLocal <= today && releaseDateLocal >= thirtyDaysAgo) {
      recentReleases.push(game)
    } else if (releaseDateLocal > today) {
      upcomingGames.push(game)
    }
  })

  // Group upcoming games by month/year
  const monthGroups = {}

  upcomingGames.forEach(game => {
    const date = new Date(game.release_date)
    // Extract UTC date components to avoid timezone issues
    const year = date.getUTCFullYear()
    const month = date.getUTCMonth()
    const day = date.getUTCDate()
    const utcDate = new Date(year, month, day)

    const monthYear = utcDate.toLocaleDateString('en-US', { year: 'numeric', month: 'long' })

    if (!monthGroups[monthYear]) {
      monthGroups[monthYear] = {
        monthYear,
        sortKey: utcDate.getTime(),
        games: [],
        isTBD: false
      }
    }
    monthGroups[monthYear].games.push(game)
  })

  // Sort month groups and games within each month
  const sortedGroups = Object.values(monthGroups)
    .sort((a, b) => a.sortKey - b.sortKey)
    .map(group => ({
      ...group,
      games: group.games.sort((a, b) => new Date(a.release_date) - new Date(b.release_date))
    }))

  // Add recent releases at the beginning if there are any
  if (recentReleases.length > 0) {
    sortedGroups.unshift({
      monthYear: 'Recent Releases',
      sortKey: -1,
      games: recentReleases.sort((a, b) => new Date(b.release_date) - new Date(a.release_date)), // Most recent first
      isTBD: false
    })
  }

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
