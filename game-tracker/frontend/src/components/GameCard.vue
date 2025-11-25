<template>
  <div class="bg-gray-800 rounded-lg shadow-xl hover:shadow-2xl transition-all duration-300 border border-gray-700 hover:border-gray-600 transform hover:scale-105 overflow-hidden">
    <div class="flex">
      <!-- Full height cover art on left - 528x704 aspect ratio (3:4) - fills card height -->
      <div class="flex-shrink-0 self-stretch" style="width: 165px;">
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
      <div class="flex-1 p-4 flex flex-col">
        <!-- Fixed 2-line title -->
        <h3 class="text-lg font-semibold mb-2 text-white line-clamp-2" style="min-height: 3.5rem;">{{ game.title }}</h3>

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

        <div v-if="game.release_date" class="text-sm text-gray-400 mb-2">
          Release: {{ formatDate(game.release_date) }}
        </div>
        <div v-else class="text-sm text-gray-400 mb-2" style="height: 1.25rem;"></div>

        <!-- Push status picker to bottom -->
        <div class="mt-auto pt-2">
          <StatusPicker
            :model-value="game.status"
            @update:model-value="(newStatus) => $emit('update-status', game.id, newStatus)"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { getPlatformColor, sortPlatforms } from '../lib/platformColors'
import StatusPicker from './StatusPicker.vue'

const props = defineProps({
  game: {
    type: Object,
    required: true
  }
})

defineEmits(['update-status'])

const sortedPlatforms = computed(() => sortPlatforms(props.game.platforms || []))

function formatDate(dateString) {
  const date = new Date(dateString)
  return date.toLocaleDateString()
}
</script>
