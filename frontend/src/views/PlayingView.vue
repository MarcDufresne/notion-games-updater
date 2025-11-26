<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-6 text-white">Currently Playing</h1>

    <TransitionGroup
      v-if="gamesStore.playing.length > 0"
      name="game-list"
      tag="div"
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"
    >
      <GameCard
        v-for="game in gamesStore.playing"
        :key="game.id"
        :game="game"
        @update-status="handleStatusUpdate"
        @card-click="openModal"
      />
    </TransitionGroup>
    <div v-else class="text-gray-400 text-center py-8 bg-gray-800/50 rounded-lg border border-gray-700">
      No games currently being played.
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
import { onMounted } from 'vue'
import { useGamesStore } from '../stores/games'
import GameCard from '../components/GameCard.vue'
import GameDetailsModal from '../components/GameDetailsModal.vue'
import { useGameModal } from '../composables/useGameModal'

const gamesStore = useGamesStore()
const { isModalOpen, selectedGame, openModal, closeModal, handleStatusUpdate, handleDeleteGame, handleMatchUpdated } = useGameModal('playing')

onMounted(() => {
  gamesStore.fetchGames('playing')
})
</script>

<style scoped>
/* Smooth transitions for game cards */
.game-list-move,
.game-list-enter-active,
.game-list-leave-active {
  transition: all 0.5s ease;
}

.game-list-enter-from {
  opacity: 0;
  transform: translateY(-30px);
}

.game-list-leave-to {
  opacity: 0;
  transform: translateY(30px);
}

/* Ensure leaving items are positioned absolutely to allow smooth movement */
.game-list-leave-active {
  position: absolute;
}
</style>
