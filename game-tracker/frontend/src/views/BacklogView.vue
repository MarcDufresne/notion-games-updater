<template>
  <div>
    <h2 class="text-2xl font-bold mb-4">Backlog</h2>
    <div v-if="store.loading" class="text-center">Loading...</div>
    <div v-else>
        <div class="mb-8">
            <h3 class="text-xl font-semibold mb-2 text-yellow-400">On Break</h3>
            <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-4">
                <div v-for="game in onBreak" :key="game.id" class="bg-gray-800 rounded-lg overflow-hidden shadow-lg">
                    <img :src="game.cover_url || 'https://via.placeholder.com/300x400'" alt="Cover" class="w-full h-48 object-cover">
                    <div class="p-4">
                    <h3 class="text-lg font-bold truncate">{{ game.title }}</h3>
                    <p class="text-gray-400 text-sm">{{ game.platforms?.join(', ') }}</p>
                    </div>
                </div>
            </div>
        </div>

        <div>
            <h3 class="text-xl font-semibold mb-2 text-blue-400">Up Next</h3>
            <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-4">
                <div v-for="game in upNext" :key="game.id" class="bg-gray-800 rounded-lg overflow-hidden shadow-lg">
                    <img :src="game.cover_url || 'https://via.placeholder.com/300x400'" alt="Cover" class="w-full h-48 object-cover">
                    <div class="p-4">
                    <h3 class="text-lg font-bold truncate">{{ game.title }}</h3>
                    <p class="text-gray-400 text-sm">{{ game.platforms?.join(', ') }}</p>
                    </div>
                </div>
            </div>
        </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, computed } from 'vue'
import { useGameStore } from '../stores/games'

const store = useGameStore()

const onBreak = computed(() => store.backlog.filter(g => g.status === 'Break'))
const upNext = computed(() => store.backlog.filter(g => g.status === 'Backlog'))

onMounted(() => {
  store.fetchGames('backlog')
})
</script>
