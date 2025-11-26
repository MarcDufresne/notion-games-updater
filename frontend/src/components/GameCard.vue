<template>
  <div
    :class="[
      'bg-gray-800 rounded-lg shadow-xl hover:shadow-2xl border border-gray-700 hover:border-gray-600 transform hover:scale-105 overflow-hidden cursor-pointer',
      'transition-all duration-300',
      isUpdating ? 'opacity-50' : 'opacity-100'
    ]"
    @click="handleCardClick"
  >
    <div class="flex">
      <!-- Full height cover art on left - 528x704 aspect ratio (3:4) - fills card height -->
      <div class="flex-shrink-0 self-stretch" style="width: 200px;">
        <img
          v-if="game.cover_url"
          :src="game.cover_url"
          :alt="game.title"
          class="w-full h-full object-cover object-center"
        />
        <div v-else class="w-full h-full bg-gray-700 flex items-center justify-center border-r border-gray-600">
          <span class="text-gray-500 text-xs">No Cover</span>
        </div>
      </div>

      <!-- Content area -->
      <div class="flex-1 p-4 flex flex-col">
        <!-- Fixed 2-line title with match status indicator -->
        <div class="flex items-start gap-2 mb-2" style="min-height: 3.5rem;">
          <h3 class="text-lg font-semibold text-white line-clamp-2 flex-1">{{ game.title }}</h3>
          <!-- Match Status Indicator -->
          <div v-if="!game.igdb_id || game.igdb_id === 0" class="flex-shrink-0 mt-2" :title="getMatchStatusTooltip()">
            <svg v-if="game.match_status === 'needs_review'" class="w-4 h-4 text-red-400" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
            </svg>
            <svg v-else-if="game.match_status === 'multiple'" class="w-4 h-4 text-yellow-400" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
            </svg>
            <svg v-else-if="game.match_status === 'no_match'" class="w-4 h-4 text-gray-500" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
            </svg>
            <svg v-else class="w-4 h-4 text-gray-500" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
            </svg>
          </div>
        </div>

        <!-- Single line genres with overflow hidden and gradient fade -->
        <div v-if="game.genres && game.genres.length" class="mb-2 relative" style="height: 1.75rem;">
          <div class="absolute inset-0 overflow-hidden">
            <div class="flex gap-1 flex-nowrap">
              <span v-for="genre in game.genres" :key="genre" class="inline-block bg-blue-900/50 text-blue-300 text-xs px-2 py-1 rounded border border-blue-700/50 flex-shrink-0 whitespace-nowrap">
                {{ genre }}
              </span>
            </div>
          </div>
          <!-- Gradient fade on right edge -->
          <div class="absolute right-0 top-0 bottom-0 w-12 bg-gradient-to-l from-gray-800 to-transparent pointer-events-none"></div>
        </div>
        <div v-else class="mb-2" style="height: 1.75rem;"></div>

        <!-- Single line platforms with overflow hidden and gradient fade -->
        <div v-if="sortedPlatforms.length" class="mb-2 relative" style="height: 1.75rem;">
          <div class="absolute inset-0 overflow-hidden">
            <div class="flex gap-1 flex-nowrap">
              <span
                v-for="platform in sortedPlatforms"
                :key="platform"
                :class="[
                  'inline-block text-xs px-2 py-1 rounded border flex-shrink-0 whitespace-nowrap',
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
        <div v-else class="mb-2" style="height: 1.75rem;"></div>

        <div class="mb-2 relative">
          <span class="text-sm font-medium text-gray-400">Rating: </span>
          <span v-if="game.rating" class="font-semibold text-blue-400">{{ game.rating }}</span>
          <span v-else class="font-semibold text-gray-500">N/A</span>
        </div>

        <div v-if="game.release_date" :class="['text-sm mb-2', isReleased ? 'text-gray-300' : 'text-gray-500']">
          <span v-if="isReleased">Released: </span>
          <span v-else>Releases: </span>
          {{ formatReleaseDate(game.release_date) }}
        </div>
        <div v-else class="text-sm text-gray-400 mb-2" style="height: 1.25rem;"></div>

        <!-- Show played date for completed games, with placeholder for consistent height -->
        <div v-if="isCompletedGame && game.date_played && !isDatePlayedSentinel(game.date_played)" class="text-sm text-green-400 mb-2">
          Played: {{ formatDatePlayed(game.date_played) }}
        </div>
        <div v-else class="text-sm mb-2" style="height: 1.25rem;"></div>

        <!-- Push status picker to bottom -->
        <div class="mt-auto pt-2" @click.stop>
          <StatusPicker
            :model-value="game.status"
            :date-played="game.date_played"
            @update:model-value="handleStatusUpdate"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { getPlatformColor, sortPlatforms } from '../lib/platformColors'
import { formatReleaseDate, formatDatePlayed, isDatePlayedSentinel } from '../lib/dateUtils'
import StatusPicker from './StatusPicker.vue'

const COMPLETED_STATUSES = new Set(['Done', 'Abandoned', "Won't Play"])

const props = defineProps({
  game: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['update-status', 'card-click'])

const isUpdating = ref(false)

const sortedPlatforms = computed(() => sortPlatforms(props.game.platforms || []))

const isCompletedGame = computed(() => COMPLETED_STATUSES.has(props.game.status))

const isReleased = computed(() => {
  if (!props.game.release_date) return false
  const date = new Date(props.game.release_date)
  const year = date.getUTCFullYear()
  const month = date.getUTCMonth()
  const day = date.getUTCDate()
  const releaseDate = new Date(year, month, day)
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return releaseDate <= today
})

function handleCardClick() {
  emit('card-click', props.game)
}

function handleStatusUpdate(newStatus, datePlayed = null) {
  isUpdating.value = true
  emit('update-status', props.game.id, newStatus, datePlayed)
  // Keep the updating state for a bit to allow the card to fade out smoothly
  setTimeout(() => {
    isUpdating.value = false
  }, 300)
}

function getMatchStatusTooltip() {
  if (!props.game.igdb_id || props.game.igdb_id === 0) {
    if (props.game.match_status === 'multiple') {
      return 'Multiple IGDB matches found'
    } else if (props.game.match_status === 'no_match') {
      return 'No IGDB match found'
    } else if (props.game.match_status === 'needs_review') {
      return 'Match conflict - needs review'
    } else {
      return 'Not matched to IGDB'
    }
  }
  return ''
}
</script>
