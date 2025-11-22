import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '@/services/api'

export const useGamesStore = defineStore('games', () => {
  // State
  const backlog = ref([])
  const playing = ref([])
  const history = ref([])
  const loading = ref(false)
  const error = ref(null)

  // Computed
  const backlogGames = computed(() => backlog.value.filter(g => g.status === 'Backlog'))
  const breakGames = computed(() => backlog.value.filter(g => g.status === 'Break'))

  // Actions
  async function fetchGames(view) {
    loading.value = true
    error.value = null
    
    try {
      const games = await api.getGames(view)
      
      switch (view) {
        case 'backlog':
          backlog.value = games
          break
        case 'playing':
          playing.value = games
          break
        case 'history':
          history.value = games
          break
        default:
          // If no view specified, organize by status
          backlog.value = games.filter(g => g.status === 'Backlog' || g.status === 'Break')
          playing.value = games.filter(g => g.status === 'Playing')
          history.value = games.filter(g => 
            g.status === 'Done' || g.status === 'Abandoned' || g.status === "Won't Play"
          )
      }
    } catch (err) {
      error.value = err.message
      console.error('Failed to fetch games:', err)
    } finally {
      loading.value = false
    }
  }

  async function updateGameStatus(gameId, newStatus, datePlayed = null) {
    loading.value = true
    error.value = null
    
    try {
      const updatedGame = await api.updateGameStatus(gameId, newStatus, datePlayed)
      
      // Remove game from current list
      backlog.value = backlog.value.filter(g => g.id !== gameId)
      playing.value = playing.value.filter(g => g.id !== gameId)
      history.value = history.value.filter(g => g.id !== gameId)
      
      // Add to appropriate list
      if (newStatus === 'Backlog' || newStatus === 'Break') {
        backlog.value.push(updatedGame)
      } else if (newStatus === 'Playing') {
        playing.value.push(updatedGame)
      } else {
        history.value.push(updatedGame)
      }
      
      return updatedGame
    } catch (err) {
      error.value = err.message
      console.error('Failed to update game status:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function saveGame(game) {
    loading.value = true
    error.value = null
    
    try {
      const savedGame = await api.saveGame(game)
      
      // Add or update in appropriate list
      if (savedGame.status === 'Backlog' || savedGame.status === 'Break') {
        const index = backlog.value.findIndex(g => g.id === savedGame.id)
        if (index >= 0) {
          backlog.value[index] = savedGame
        } else {
          backlog.value.push(savedGame)
        }
      } else if (savedGame.status === 'Playing') {
        const index = playing.value.findIndex(g => g.id === savedGame.id)
        if (index >= 0) {
          playing.value[index] = savedGame
        } else {
          playing.value.push(savedGame)
        }
      } else {
        const index = history.value.findIndex(g => g.id === savedGame.id)
        if (index >= 0) {
          history.value[index] = savedGame
        } else {
          history.value.push(savedGame)
        }
      }
      
      return savedGame
    } catch (err) {
      error.value = err.message
      console.error('Failed to save game:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  function clearError() {
    error.value = null
  }

  return {
    // State
    backlog,
    playing,
    history,
    loading,
    error,
    // Computed
    backlogGames,
    breakGames,
    // Actions
    fetchGames,
    updateGameStatus,
    saveGame,
    clearError
  }
})
