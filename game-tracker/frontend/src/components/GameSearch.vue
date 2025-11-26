<template>
  <div class="relative w-full">
    <input
      v-model="searchText"
      type="text"
      placeholder="Search for a game..."
      class="w-full px-4 py-3 pr-12 text-lg border border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-gray-700 text-white placeholder-gray-400"
      @focus="showDropdown = true"
    />

    <!-- Clear Button -->
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

    <div
      v-if="showDropdown && (searchResults.length > 0 || searchText.trim().length > 0)"
      class="absolute z-10 w-full mt-2 bg-gray-800 border border-gray-600 rounded-lg shadow-2xl max-h-96 overflow-y-auto"
    >
      <!-- Loading State -->
      <div v-if="isSearching" class="p-4 text-center text-gray-400">
        Searching...
      </div>

      <!-- Search Results -->
      <div
        v-for="result in searchResults"
        :key="result.id"
        class="flex items-center p-3 hover:bg-gray-700 cursor-pointer border-b border-gray-700 transition-colors"
        @click="selectGame(result)"
      >
        <img
          v-if="result.cover_url"
          :src="result.cover_url"
          :alt="result.name"
          class="w-12 h-16 object-cover rounded mr-3 shadow-md"
        />
        <div class="flex-1">
          <div class="font-semibold text-white">{{ result.name }}</div>
          <div v-if="result.release_year" class="text-sm text-gray-400">
            {{ result.release_year }}
          </div>
        </div>
      </div>

      <!-- No Results Message -->
      <div
        v-if="!isSearching && searchResults.length === 0 && searchText.trim().length > 0"
        class="p-4 text-center text-gray-400 border-b border-gray-700"
      >
        No results found
      </div>

      <!-- Manual Entry Option - Always at the end -->
      <div
        v-if="searchText.trim().length > 0 && !isSearching"
        class="p-3 hover:bg-gray-700 cursor-pointer border-t-2 border-gray-600 transition-colors"
        @click="createManual"
      >
        <div class="flex items-center gap-2">
          <svg class="w-5 h-5 text-blue-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          <span class="font-semibold text-blue-400">Manually create an entry for "{{ searchText }}"</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useDebounceFn } from '@vueuse/core'
import { useGamesStore } from '../stores/games'

const gamesStore = useGamesStore()

const searchText = ref('')
const searchResults = ref([])
const showDropdown = ref(false)
const isSearching = ref(false)

const emit = defineEmits(['game-created'])

// Debounced search function
const debouncedSearch = useDebounceFn(async (query) => {
  if (query.trim().length < 2) {
    searchResults.value = []
    isSearching.value = false
    return
  }

  isSearching.value = true
  try {
    const results = await gamesStore.searchIGDB(query)
    searchResults.value = results || []
  } catch (error) {
    console.error('Search failed:', error)
    searchResults.value = []
  } finally {
    isSearching.value = false
  }
}, 300)

// Watch for search text changes
watch(searchText, (newValue) => {
  debouncedSearch(newValue)
})

async function selectGame(result) {
  try {
    const game = await gamesStore.createGame({
      title: result.name,
      igdb_id: result.id,
      status: 'Backlog'
    })

    searchText.value = ''
    searchResults.value = []
    showDropdown.value = false
    emit('game-created', game)
  } catch (error) {
    console.error('Failed to create game:', error)
    // Show toast notification
    if (window.$toast) {
      window.$toast.error(error.message, 5000)
    }
    showDropdown.value = false
  }
}

async function createManual() {
  try {
    const game = await gamesStore.createGame({
      title: searchText.value,
      status: 'Backlog'
    })

    searchText.value = ''
    searchResults.value = []
    showDropdown.value = false
    emit('game-created', game)
  } catch (error) {
    console.error('Failed to create game:', error)
    // Show toast notification
    if (window.$toast) {
      window.$toast.error(error.message, 5000)
    }
    showDropdown.value = false
  }
}

function clearSearch() {
  searchText.value = ''
  searchResults.value = []
  showDropdown.value = false
}

// Close dropdown when clicking outside
function handleClickOutside(event) {
  if (!event.target.closest('.relative')) {
    showDropdown.value = false
  }
}

// Add event listener for clicking outside
if (typeof window !== 'undefined') {
  window.addEventListener('click', handleClickOutside)
}
</script>
