<template>
  <div>
    <h2 class="text-2xl font-bold mb-4">History</h2>
    <div v-if="store.loading" class="text-center">Loading...</div>
    <div v-else class="space-y-4">
        <div v-for="(group, year) in groupedHistory" :key="year">
            <h3 class="text-xl font-semibold mb-2 border-b border-gray-700 pb-1">{{ year }}</h3>
            <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-4">
                <div v-for="game in group" :key="game.id" class="bg-gray-800 rounded-lg overflow-hidden shadow-lg opacity-80 hover:opacity-100 transition-opacity">
                    <img :src="game.cover_url || 'https://via.placeholder.com/300x400'" alt="Cover" class="w-full h-48 object-cover grayscale hover:grayscale-0 transition-all">
                    <div class="p-4">
                        <h3 class="text-lg font-bold truncate">{{ game.title }}</h3>
                        <div class="flex justify-between items-center mt-2">
                            <span class="text-xs px-2 py-1 rounded bg-gray-700" :class="statusColor(game.status)">{{ game.status }}</span>
                            <span class="text-gray-400 text-sm">{{ formatDate(game.date_played) }}</span>
                        </div>
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

onMounted(() => {
  store.fetchGames('history')
})

const groupedHistory = computed(() => {
    const groups = {}
    store.history.forEach(game => {
        const year = game.date_played ? new Date(game.date_played).getFullYear() : 'Unknown'
        if (!groups[year]) groups[year] = []
        groups[year].push(game)
    })
    // Sort years descending
    return Object.keys(groups).sort((a, b) => b - a).reduce((acc, key) => {
        acc[key] = groups[key]
        return acc
    }, {})
})

function statusColor(status) {
    switch(status) {
        case 'Done': return 'text-green-400';
        case 'Abandoned': return 'text-red-400';
        case "Won't Play": return 'text-gray-500';
        default: return 'text-gray-400';
    }
}

function formatDate(dateStr) {
    if (!dateStr) return ''
    return new Date(dateStr).toLocaleDateString()
}
</script>
