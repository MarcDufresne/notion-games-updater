<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="isOpen && game"
        class="fixed inset-0 z-[60] flex items-center justify-center p-4 bg-black/80"
        @click.self="closeModal"
      >
        <div
          class="bg-gray-800 rounded-2xl shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-hidden border border-gray-700"
          @click.stop
        >
          <!-- Header -->
          <div class="flex items-start justify-between p-6 border-b border-gray-700 bg-gray-900/50">
            <div>
              <h2 class="text-2xl font-bold text-white mb-1">Fix Match</h2>
              <p class="text-gray-400 text-sm">Search IGDB and select the correct match for:</p>
              <p class="text-white font-semibold mt-1">{{ game.title }}</p>
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

          <!-- Search Box -->
          <div class="p-6 border-b border-gray-700">
            <div class="relative">
              <input
                v-model="searchText"
                type="text"
                placeholder="Search IGDB..."
                class="w-full px-4 py-3 pr-10 text-lg border border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-gray-700 text-white placeholder-gray-400"
              />
              <button
                v-if="searchText.trim().length > 0"
                @click="clearSearch"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-white transition-colors p-1 rounded hover:bg-gray-600"
                title="Clear search"
              >
                <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                </svg>
              </button>
            </div>
          </div>

          <!-- Results -->
          <div class="overflow-y-auto max-h-[50vh] p-6">
            <!-- Loading -->
            <div v-if="isSearching" class="text-center py-8">
              <div class="text-gray-400">Searching IGDB...</div>
            </div>

            <!-- Results List -->
            <div v-else-if="searchResults.length > 0" class="space-y-2">
              <button
                v-for="result in searchResults"
                :key="result.id"
                @click="selectMatch(result)"
                :disabled="isUpdating"
                class="w-full flex items-center p-3 hover:bg-gray-700 rounded-lg cursor-pointer border border-gray-700 hover:border-gray-600 transition-all text-left disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <img
                  v-if="result.cover_url"
                  :src="result.cover_url"
                  :alt="result.name"
                  class="w-12 h-16 object-cover rounded mr-3 shadow-md flex-shrink-0"
                />
                <div class="flex-shrink-0 w-12 h-16 bg-gray-700 rounded mr-3 flex items-center justify-center" v-else>
                  <span class="text-gray-500 text-xs">No Cover</span>
                </div>
                <div class="flex-1 min-w-0">
                  <div class="font-semibold text-white truncate">{{ result.name }}</div>
                  <div v-if="result.release_year" class="text-sm text-gray-400">
                    {{ result.release_year }}
                  </div>
                </div>
                <svg class="w-5 h-5 text-green-400 flex-shrink-0 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
              </button>
            </div>

            <!-- No Results -->
            <div v-else-if="hasSearched && searchResults.length === 0" class="text-center py-8 text-gray-400">
              No results found. Try a different search term.
            </div>

            <!-- Initial State -->
            <div v-else class="text-center py-8 text-gray-400">
              Type to search IGDB automatically...
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useDebounceFn } from '@vueuse/core'
import { api } from '../lib/api'
import { useGamesStore } from '../stores/games'

const gamesStore = useGamesStore()

const props = defineProps({
  isOpen: {
    type: Boolean,
    required: true
  },
  game: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close', 'match-updated'])

const searchText = ref('')
const searchResults = ref([])
const isSearching = ref(false)
const isUpdating = ref(false)
const hasSearched = ref(false)

// Debounced search function
const debouncedSearch = useDebounceFn(async (query) => {
  if (query.trim().length < 2) {
    searchResults.value = []
    isSearching.value = false
    hasSearched.value = false
    return
  }

  isSearching.value = true
  hasSearched.value = true
  try {
    const results = await api.searchGames(query.trim())
    searchResults.value = results
  } catch (error) {
    console.error('Search failed:', error)
    if (window.$toast) {
      window.$toast.error('Failed to search IGDB', 3000)
    }
  } finally {
    isSearching.value = false
  }
}, 300)

// Watch for search text changes
watch(searchText, (newValue) => {
  debouncedSearch(newValue)
})

watch(() => props.isOpen, (isOpen) => {
  if (isOpen && props.game) {
    // Pre-fill search with game title
    searchText.value = props.game.title
    searchResults.value = []
    hasSearched.value = false
    // The watch on searchText will trigger the search automatically
  }
})

function closeModal() {
  if (!isUpdating.value) {
    emit('close')
  }
}

function clearSearch() {
  searchText.value = ''
  searchResults.value = []
  hasSearched.value = false
}

async function selectMatch(result) {
  isUpdating.value = true
  try {
    const updatedGame = await gamesStore.updateGameMatch(props.game.id, result.id)
    if (window.$toast) {
      window.$toast.success(`Successfully matched "${props.game.title}" to "${result.name}"`, 5000)
    }
    emit('match-updated', updatedGame)
    emit('close')
  } catch (error) {
    console.error('Failed to update match:', error)
    if (window.$toast) {
      window.$toast.error(error.message, 5000)
    }
  } finally {
    isUpdating.value = false
  }
}
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-active .bg-gray-800,
.modal-leave-active .bg-gray-800 {
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .bg-gray-800,
.modal-leave-to .bg-gray-800 {
  transform: scale(0.95);
  opacity: 0;
}
</style>
