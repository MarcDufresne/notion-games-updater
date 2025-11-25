<template>
  <div class="relative w-full">
    <input
      v-model="searchText"
      type="text"
      placeholder="Search for a game..."
      class="w-full px-4 py-3 text-lg border border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-gray-700 text-white placeholder-gray-400"
      @focus="showDropdown = true"
    />

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

      <!-- Manual Entry Option -->
      <div
        v-if="searchText.trim().length > 0"
        class="p-3 hover:bg-gray-700 cursor-pointer border-t-2 border-gray-600 transition-colors"
        @click="createManual"
      >
        <div class="font-semibold text-blue-400">
          ✏️ Manually create an entry for "{{ searchText }}"
        </div>
      </div>

      <!-- No Results -->
      <div
        v-if="!isSearching && searchResults.length === 0 && searchText.trim().length > 0"
        class="p-4 text-center text-gray-400"
      >
        No results found
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
  }
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
