<template>
  <div>
    <h2 class="text-2xl font-bold mb-4">Playing Now</h2>
    <div v-if="store.loading" class="text-center">Loading...</div>
    <div v-else class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-4">
      <div v-for="game in store.playing" :key="game.id" class="bg-gray-800 rounded-lg overflow-hidden shadow-lg">
        <img :src="game.cover_url || 'https://via.placeholder.com/300x400'" alt="Cover" class="w-full h-48 object-cover">
        <div class="p-4">
          <h3 class="text-lg font-bold truncate">{{ game.title }}</h3>
          <p class="text-gray-400 text-sm">{{ game.platforms?.join(', ') }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useGameStore } from '../stores/games'

const store = useGameStore()

onMounted(() => {
  store.fetchGames('playing')
})
</script>
