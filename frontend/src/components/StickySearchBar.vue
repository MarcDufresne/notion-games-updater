<template>
  <Transition name="slide-down">
    <div
      v-if="showSticky"
      class="fixed top-0 left-0 right-0 z-50 bg-gray-800/95 backdrop-blur-sm border-b border-gray-700 shadow-lg"
    >
      <div class="container mx-auto px-4 py-3">
        <GameSearch @game-created="handleGameCreated" @open-existing-game="handleOpenExistingGame" />
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import GameSearch from './GameSearch.vue'

const emit = defineEmits(['game-created', 'open-existing-game'])

const showSticky = ref(false)
let lastScrollY = 0
let ticking = false

const SCROLL_THRESHOLD = 200 // Show sticky bar after scrolling down this many pixels
const SCROLL_UP_THRESHOLD = 15 // Minimum scroll up distance to trigger show

function handleGameCreated(game) {
  emit('game-created', game)
}

function handleOpenExistingGame(game) {
  emit('open-existing-game', game)
}

function updateStickyBar() {
  const currentScrollY = window.scrollY

  // Only show if we've scrolled down past threshold
  if (currentScrollY < SCROLL_THRESHOLD) {
    showSticky.value = false
  } else {
    // Check if user is scrolling up
    const isScrollingUp = currentScrollY < lastScrollY
    const scrollDelta = Math.abs(currentScrollY - lastScrollY)

    if (isScrollingUp && scrollDelta > SCROLL_UP_THRESHOLD) {
      showSticky.value = true
    } else if (!isScrollingUp && scrollDelta > SCROLL_UP_THRESHOLD) {
      showSticky.value = false
    }
  }

  lastScrollY = currentScrollY
  ticking = false
}

function onScroll() {
  if (!ticking) {
    window.requestAnimationFrame(updateStickyBar)
    ticking = true
  }
}

onMounted(() => {
  window.addEventListener('scroll', onScroll, { passive: true })
  lastScrollY = window.scrollY
})

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
})
</script>

<style scoped>
.slide-down-enter-active,
.slide-down-leave-active {
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.slide-down-enter-from {
  transform: translateY(-100%);
  opacity: 0;
}

.slide-down-leave-to {
  transform: translateY(-100%);
  opacity: 0;
}
</style>
