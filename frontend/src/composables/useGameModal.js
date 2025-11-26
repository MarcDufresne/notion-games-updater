import { ref, nextTick, triggerRef } from 'vue'
import { useGamesStore } from '../stores/games'

export function useGameModal(viewName) {
  const gamesStore = useGamesStore()
  const isModalOpen = ref(false)
  const selectedGame = ref({})

  function openModal(game) {
    selectedGame.value = { ...game }
    isModalOpen.value = true
  }

  function closeModal() {
    isModalOpen.value = false
    setTimeout(() => {
      if (!isModalOpen.value) {
        selectedGame.value = {}
      }
    }, 300)
  }

  async function handleStatusUpdate(gameId, newStatus, datePlayed = null) {
    const updatedGame = await gamesStore.updateStatus(gameId, newStatus, datePlayed)

    if (updatedGame) {
      const mergedGame = { ...selectedGame.value, ...updatedGame }
      selectedGame.value = mergedGame

      await nextTick()
      triggerRef(selectedGame)
    } else {
      if (selectedGame.value.id === gameId) {
        selectedGame.value = {
          ...selectedGame.value,
          status: newStatus,
          date_played: datePlayed || selectedGame.value.date_played,
          updated_at: new Date().toISOString()
        }
        await nextTick()
        triggerRef(selectedGame)
      }
    }
  }

  async function handleDeleteGame(gameId) {
    await gamesStore.deleteGame(gameId)
    closeModal()
  }

  function handleMatchUpdated(updatedGame) {
    if (updatedGame && updatedGame.id === selectedGame.value.id) {
      const mergedGame = { ...selectedGame.value, ...updatedGame }
      selectedGame.value = mergedGame
      nextTick()
      triggerRef(selectedGame)
    }
  }

  return {
    isModalOpen,
    selectedGame,
    openModal,
    closeModal,
    handleStatusUpdate,
    handleDeleteGame,
    handleMatchUpdated
  }
}
