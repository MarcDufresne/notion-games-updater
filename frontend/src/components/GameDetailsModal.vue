<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="isOpen && game.id"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80"
        @click.self="closeModal"
      >
        <div
          class="bg-gray-800 rounded-2xl shadow-2xl max-w-4xl w-full max-h-[90vh] overflow-hidden border border-gray-700"
          @click.stop
        >
          <!-- Header with menu and close button -->
          <div class="flex items-start justify-between p-6 border-b border-gray-700 bg-gray-900/50">
            <h2 class="text-2xl font-bold text-white pr-8">{{ game.title }}</h2>
            <div class="flex items-center gap-2">
              <!-- 3-dot menu -->
              <div class="relative" ref="menuContainer">
                <button
                  @click="toggleMenu"
                  class="flex-shrink-0 text-gray-400 hover:text-white transition-colors p-2 hover:bg-gray-700 rounded-lg"
                  aria-label="Menu"
                >
                  <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M12 8c1.1 0 2-.9 2-2s-.9-2-2-2-2 .9-2 2 .9 2 2 2zm0 2c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2zm0 6c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2z"/>
                  </svg>
                </button>

                <!-- Dropdown menu - Teleported to body -->
                <Teleport to="body">
                  <div
                    v-if="isMenuOpen"
                    ref="menuDropdown"
                    :style="menuStyle"
                    class="fixed z-[60] bg-gray-800 border border-gray-600 rounded-lg shadow-2xl py-1 min-w-[200px]"
                  >
                    <!-- Fix Match -->
                    <button
                      @click="openFixMatchModalFromMenu"
                      class="w-full px-4 py-2 text-left text-sm text-blue-400 hover:bg-gray-700 transition-colors flex items-center gap-2"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                      </svg>
                      Fix Match
                    </button>

                    <!-- Delete Game -->
                    <button
                      @click="confirmDelete"
                      class="w-full px-4 py-2 text-left text-sm text-red-400 hover:bg-gray-700 transition-colors flex items-center gap-2"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                      </svg>
                      Delete Game
                    </button>
                  </div>
                </Teleport>
              </div>

              <button
                @click="closeModal"
                class="flex-shrink-0 text-gray-400 hover:text-white transition-colors p-2 hover:bg-gray-700 rounded-lg"
                aria-label="Close"
              >
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>

          <!-- Content -->
          <div class="overflow-y-auto max-h-[calc(90vh-80px)]">
            <div class="p-6">
              <div class="flex flex-col md:flex-row gap-6">
                <!-- Cover Image -->
                <div class="flex-shrink-0">
                  <div class="w-64 rounded-lg overflow-hidden shadow-xl border border-gray-700">
                    <img
                      v-if="game.cover_url"
                      :src="game.cover_url"
                      :alt="game.title"
                      class="w-full h-auto object-cover"
                    />
                    <div v-else class="w-64 h-85 bg-gray-700 flex items-center justify-center">
                      <span class="text-gray-500 py-36">No Cover</span>
                    </div>
                  </div>
                </div>

                <!-- Game Information -->
                <div class="flex-1 space-y-4">
                  <!-- Status -->
                  <div>
                    <h3 class="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-2">Status</h3>
                    <div class="flex items-center gap-2">
                      <StatusPicker
                        :model-value="game.status"
                        :date-played="game.date_played"
                        @update:model-value="(newStatus, datePlayed) => $emit('update-status', game.id, newStatus, datePlayed)"
                      />
                    </div>
                  </div>

                  <!-- Rating -->
                  <div v-if="game.rating">
                    <h3 class="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-2">Rating</h3>
                    <div class="flex items-center gap-2">
                      <div class="text-2xl font-bold text-blue-400">{{ game.rating }}</div>
                      <div class="text-gray-500">/100</div>
                    </div>
                  </div>

                  <!-- Genres -->
                  <div v-if="game.genres && game.genres.length">
                    <h3 class="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-2">Genres</h3>
                    <div class="flex flex-wrap gap-2">
                      <span
                        v-for="genre in game.genres"
                        :key="genre"
                        class="inline-block bg-blue-900/50 text-blue-300 text-sm px-3 py-1 rounded-full border border-blue-700/50"
                      >
                        {{ genre }}
                      </span>
                    </div>
                  </div>

                  <!-- Platforms -->
                  <div v-if="game.platforms && game.platforms.length">
                    <h3 class="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-2">Platforms</h3>
                    <div class="flex flex-wrap gap-2">
                      <span
                        v-for="platform in sortedPlatforms"
                        :key="platform"
                        :class="[
                          'inline-block text-sm px-3 py-1 rounded-full border',
                          getPlatformColor(platform).bg,
                          getPlatformColor(platform).text,
                          getPlatformColor(platform).border
                        ]"
                      >
                        {{ platform }}
                      </span>
                    </div>
                  </div>

                  <!-- Release Date -->
                  <div v-if="game.release_date">
                    <h3 class="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-2">Release Date</h3>
                    <div class="text-white">{{ formatReleaseDate(game.release_date) }}</div>
                  </div>

                  <!-- Steam Store Link -->
                  <div v-if="game.steam_url">
                    <h3 class="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-2">Steam Store</h3>
                    <a
                      :href="game.steam_url"
                      target="_blank"
                      rel="noopener noreferrer"
                      class="inline-flex items-center gap-2 text-blue-400 hover:text-blue-300 transition-colors"
                    >
                      <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                        <path d="M11.979 0C5.678 0 .511 4.86.022 11.037l6.432 2.658c.545-.371 1.203-.59 1.912-.59.063 0 .125.004.188.006l2.861-4.142V8.91c0-2.495 2.028-4.524 4.524-4.524 2.494 0 4.524 2.031 4.524 4.527s-2.03 4.525-4.524 4.525h-.105l-4.076 2.911c0 .052.004.105.004.159 0 1.875-1.515 3.396-3.39 3.396-1.635 0-3.016-1.173-3.331-2.727L.436 15.27C1.862 20.307 6.486 24 11.979 24c6.627 0 11.999-5.373 11.999-12S18.605 0 11.979 0zM7.54 18.21l-1.473-.61c.262.543.714.999 1.314 1.25 1.297.539 2.793-.076 3.332-1.375.263-.63.264-1.319.005-1.949s-.75-1.121-1.377-1.383c-.624-.26-1.29-.249-1.878-.03l1.523.63c.956.4 1.409 1.5 1.009 2.455-.397.957-1.497 1.41-2.454 1.012H7.54zm11.415-9.303c0-1.662-1.353-3.015-3.015-3.015-1.665 0-3.015 1.353-3.015 3.015 0 1.665 1.35 3.015 3.015 3.015 1.663 0 3.015-1.35 3.015-3.015zm-5.273-.005c0-1.252 1.013-2.266 2.265-2.266 1.249 0 2.266 1.014 2.266 2.266 0 1.251-1.017 2.265-2.266 2.265-1.253 0-2.265-1.014-2.265-2.265z"/>
                      </svg>
                      View on Steam
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                      </svg>
                    </a>
                  </div>

                  <!-- Official Website Link -->
                  <div v-if="game.official_url">
                    <h3 class="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-2">Official Website</h3>
                    <a
                      :href="game.official_url"
                      target="_blank"
                      rel="noopener noreferrer"
                      class="inline-flex items-center gap-2 text-blue-400 hover:text-blue-300 transition-colors"
                    >
                      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
                      </svg>
                      Visit Official Website
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                      </svg>
                    </a>
                  </div>

                  <!-- Date Played -->
                  <div v-if="isCompletedGame">
                    <h3 class="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-2">Date Played</h3>
                    <div v-if="!isEditingDate" class="flex items-center gap-2">
                      <div class="text-white">
                        {{ game.date_played && !isDatePlayedSentinel(game.date_played) ? formatDatePlayed(game.date_played) : 'Not set' }}
                      </div>
                      <button
                        @click="startEditingDate"
                        class="text-blue-400 hover:text-blue-300 transition-colors text-sm"
                        title="Edit date"
                      >
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                        </svg>
                      </button>
                      <button
                        v-if="game.date_played && !isDatePlayedSentinel(game.date_played)"
                        @click="clearDatePlayed"
                        class="text-red-400 hover:text-red-300 transition-colors text-sm"
                        title="Clear date"
                      >
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                      </button>
                    </div>
                    <div v-else class="flex items-center gap-2">
                      <input
                        ref="dateInput"
                        type="date"
                        v-model="editedDate"
                        :max="new Date().toISOString().split('T')[0]"
                        class="px-3 py-1.5 bg-gray-700 border border-gray-600 rounded text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                      />
                      <button
                        @click="saveDateEdit"
                        class="px-3 py-1.5 bg-blue-600 hover:bg-blue-700 text-white rounded text-sm transition-colors"
                      >
                        Save
                      </button>
                      <button
                        @click="cancelDateEdit"
                        class="px-3 py-1.5 bg-gray-700 hover:bg-gray-600 text-white rounded text-sm transition-colors"
                      >
                        Cancel
                      </button>
                    </div>
                  </div>

                  <!-- IGDB ID -->
                  <div class="pt-4 border-t border-gray-700">
                    <h3 class="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-2">IGDB ID</h3>
                    <div v-if="game.igdb_id > 0" class="text-white font-mono">
                      {{ game.igdb_id }}
                    </div>
                    <div v-else class="flex items-center gap-2">
                      <!-- Match Status Indicator -->
                      <div v-if="game.match_status === 'needs_review'" class="flex items-center gap-2 text-red-400 text-sm">
                        <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                          <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                        </svg>
                        <span>Needs Review (Duplicate Conflict)</span>
                      </div>
                      <div v-else-if="game.match_status === 'multiple'" class="flex items-center gap-2 text-yellow-400 text-sm">
                        <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                          <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                        </svg>
                        <span>Multiple Matches Found</span>
                      </div>
                      <div v-else-if="game.match_status === 'no_match'" class="flex items-center gap-2 text-gray-400 text-sm">
                        <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                        </svg>
                        <span>No Match Found</span>
                      </div>
                      <div v-else class="flex items-center gap-2 text-gray-400 text-sm">
                        <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
                        </svg>
                        <span>Not Matched</span>
                      </div>
                    </div>
                  </div>

                  <!-- Timestamps -->
                  <div class="pt-4 border-t border-gray-700">
                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 text-sm">
                      <div>
                        <span class="text-gray-400">Created:</span>
                        <span class="text-gray-300 ml-2">{{ formatDateTime(game.created_at) }}</span>
                      </div>
                      <div>
                        <span class="text-gray-400">Updated:</span>
                        <span class="text-gray-300 ml-2">{{ formatDateTime(game.updated_at) }}</span>
                      </div>
                    </div>
                  </div>

                  <!-- Sync Error (if any) -->
                  <div v-if="game.last_sync_error" class="pt-4 border-t border-gray-700">
                    <h3 class="text-sm font-semibold text-red-400 uppercase tracking-wider mb-2">Sync Error</h3>
                    <div class="text-red-300 text-sm bg-red-900/20 border border-red-700/50 rounded-lg p-3">
                      {{ game.last_sync_error }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Fix Match Modal -->
    <FixMatchModal
      :is-open="isFixMatchModalOpen"
      :game="game"
      @close="closeFixMatchModal"
      @match-updated="handleMatchUpdated"
    />

    <!-- Delete Confirmation Modal -->
    <Transition name="modal">
      <div
        v-if="showDeleteConfirm"
        class="fixed inset-0 z-[70] flex items-center justify-center p-4 bg-black/80"
        @click.self="cancelDelete"
      >
        <div class="bg-gray-800 rounded-lg shadow-2xl max-w-md w-full border border-red-700/50" @click.stop>
          <div class="p-6">
            <div class="flex items-start gap-4">
              <div class="flex-shrink-0">
                <svg class="w-6 h-6 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
              </div>
              <div class="flex-1">
                <h3 class="text-lg font-semibold text-white mb-2">Delete Game</h3>
                <p class="text-gray-300 mb-4">
                  Are you sure you want to delete <span class="font-semibold text-white">"{{ game.title }}"</span>? This action cannot be undone.
                </p>
                <div class="flex gap-3 justify-end">
                  <button
                    @click="cancelDelete"
                    class="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
                  >
                    Cancel
                  </button>
                  <button
                    @click="handleDelete"
                    class="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors flex items-center gap-2"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                    Delete
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { computed, watch, onUnmounted, ref, nextTick } from 'vue'
import { getPlatformColor, sortPlatforms } from '../lib/platformColors'
import { formatReleaseDate, formatDatePlayed, isDatePlayedSentinel } from '../lib/dateUtils'
import StatusPicker from './StatusPicker.vue'
import FixMatchModal from './FixMatchModal.vue'

const props = defineProps({
  isOpen: {
    type: Boolean,
    required: true
  },
  game: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close', 'update-status', 'delete-game', 'match-updated'])

const sortedPlatforms = computed(() => sortPlatforms(props.game.platforms || []))
const isCompletedGame = computed(() => {
  return ['Done', 'Abandoned', "Won't Play"].includes(props.game.status)
})

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

const isEditingDate = ref(false)
const editedDate = ref(null)
const dateInput = ref(null)

// Menu state
const isMenuOpen = ref(false)
const menuContainer = ref(null)
const menuDropdown = ref(null)
const menuStyle = ref({})

// Delete confirmation state
const showDeleteConfirm = ref(false)

// Fix Match modal state
const isFixMatchModalOpen = ref(false)

function closeModal() {
  isMenuOpen.value = false
  showDeleteConfirm.value = false
  isFixMatchModalOpen.value = false
  emit('close')
}

function openFixMatchModal() {
  isFixMatchModalOpen.value = true
}

function openFixMatchModalFromMenu() {
  isMenuOpen.value = false
  isFixMatchModalOpen.value = true
}

function closeFixMatchModal() {
  isFixMatchModalOpen.value = false
}

function handleMatchUpdated(updatedGame) {
  emit('match-updated', updatedGame)
  isFixMatchModalOpen.value = false
}

function toggleMenu() {
  isMenuOpen.value = !isMenuOpen.value
  if (isMenuOpen.value) {
    nextTick(() => {
      calculateMenuPosition()
    })
  }
}

function calculateMenuPosition() {
  if (!menuContainer.value) return

  const buttonRect = menuContainer.value.getBoundingClientRect()
  menuStyle.value = {
    top: `${buttonRect.bottom + 5}px`,
    right: `${window.innerWidth - buttonRect.right}px`
  }
}

function confirmDelete() {
  isMenuOpen.value = false
  showDeleteConfirm.value = true
}

function cancelDelete() {
  showDeleteConfirm.value = false
}

async function handleDelete() {
  emit('delete-game', props.game.id)
  showDeleteConfirm.value = false
}

function startEditingDate() {
  isEditingDate.value = true
  if (props.game.date_played && !isDatePlayedSentinel(props.game.date_played)) {
    const date = new Date(props.game.date_played)
    const year = date.getUTCFullYear()
    const month = String(date.getUTCMonth() + 1).padStart(2, '0')
    const day = String(date.getUTCDate()).padStart(2, '0')
    editedDate.value = `${year}-${month}-${day}`
  } else {
    editedDate.value = new Date().toISOString().split('T')[0]
  }
  nextTick(() => {
    if (dateInput.value) {
      dateInput.value.focus()
    }
  })
}

function cancelDateEdit() {
  isEditingDate.value = false
  editedDate.value = null
}

function saveDateEdit() {
  if (editedDate.value) {
    const dateObj = new Date(editedDate.value + 'T12:00:00Z')
    const isoDateString = dateObj.toISOString()

    emit('update-status', props.game.id, props.game.status, isoDateString)
  }
  isEditingDate.value = false
  editedDate.value = null
}

function clearDatePlayed() {
  const epochDate = new Date(0).toISOString()
  emit('update-status', props.game.id, props.game.status, epochDate)
}

function formatDateTime(dateString) {
  if (!dateString) return 'N/A'
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Close modal on Escape key and handle outside clicks
let escapeHandler = null
let clickHandler = null

watch(() => props.isOpen, (isOpen) => {
  if (isOpen) {
    // Add escape key listener
    escapeHandler = (e) => {
      if (e.key === 'Escape') {
        if (showDeleteConfirm.value) {
          cancelDelete()
        } else if (isEditingDate.value) {
          cancelDateEdit()
        } else if (isMenuOpen.value) {
          isMenuOpen.value = false
        } else {
          closeModal()
        }
      }
    }
    document.addEventListener('keydown', escapeHandler)

    // Add click handler to close menu when clicking outside
    clickHandler = (e) => {
      if (isMenuOpen.value && menuContainer.value && menuDropdown.value) {
        if (!menuContainer.value.contains(e.target) && !menuDropdown.value.contains(e.target)) {
          isMenuOpen.value = false
        }
      }
    }
    document.addEventListener('click', clickHandler, true) // Use capture phase
  } else {
    // Remove event listeners when modal closes
    if (escapeHandler) {
      document.removeEventListener('keydown', escapeHandler)
      escapeHandler = null
    }
    if (clickHandler) {
      document.removeEventListener('click', clickHandler, true) // Use capture phase
      clickHandler = null
    }
    // Reset states
    isEditingDate.value = false
    editedDate.value = null
    isMenuOpen.value = false
    showDeleteConfirm.value = false
  }
})

// Cleanup on component unmount
onUnmounted(() => {
  if (escapeHandler) {
    document.removeEventListener('keydown', escapeHandler)
  }
  if (clickHandler) {
    document.removeEventListener('click', clickHandler, true) // Use capture phase
  }
})
</script>

<style scoped>
/* Modal transition animations */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
  will-change: opacity;
}

.modal-enter-active > div,
.modal-leave-active > div {
  transition: transform 0.3s ease, opacity 0.3s ease;
  will-change: transform, opacity;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from > div,
.modal-leave-to > div {
  transform: scale(0.95) translateZ(0);
  opacity: 0;
}

/* Force GPU acceleration for the modal content */
.modal-enter-active > div > div,
.modal-leave-active > div > div {
  transform: translateZ(0);
  backface-visibility: hidden;
}
</style>
