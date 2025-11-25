import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '../lib/api'

export const useGamesStore = defineStore('games', () => {
  const backlog = ref([])
  const playing = ref([])
  const history = ref([])
  const loading = ref(false)
  const error = ref(null)

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
          // If no view specified, distribute games based on status
          backlog.value = games.filter(g => g.status === 'Backlog' || g.status === 'Break')
          playing.value = games.filter(g => g.status === 'Playing')
          history.value = games.filter(g => ['Done', 'Abandoned', "Won't Play"].includes(g.status))
      }
    } catch (e) {
      error.value = e.message
      console.error('Failed to fetch games:', e)
    } finally {
      loading.value = false
    }
  }

  async function createGame(gameData) {
    loading.value = true
    error.value = null
    try {
      const game = await api.createGame(gameData)

      // Add to appropriate list based on status
      if (game.status === 'Backlog' || game.status === 'Break') {
        backlog.value.push(game)
      } else if (game.status === 'Playing') {
        playing.value.push(game)
      } else {
        history.value.push(game)
      }

      return game
    } catch (e) {
      error.value = e.message
      console.error('Failed to create game:', e)
      throw e
    } finally {
      loading.value = false
    }
  }

  async function updateStatus(gameId, status) {
    loading.value = true
    error.value = null
    try {
      const updatedGame = await api.updateGameStatus(gameId, status)

      // Remove from all lists
      backlog.value = backlog.value.filter(g => g.id !== gameId)
      playing.value = playing.value.filter(g => g.id !== gameId)
      history.value = history.value.filter(g => g.id !== gameId)

      // Add to appropriate list
      if (status === 'Backlog' || status === 'Break') {
        backlog.value.push(updatedGame)
      } else if (status === 'Playing') {
        playing.value.push(updatedGame)
      } else {
        history.value.push(updatedGame)
      }

      return updatedGame
    } catch (e) {
      error.value = e.message
      console.error('Failed to update game status:', e)
      throw e
    } finally {
      loading.value = false
    }
  }

  async function searchIGDB(query) {
    if (!query || query.trim().length < 2) {
      return []
    }

    try {
      return await api.searchGames(query.trim())
    } catch (e) {
      console.error('Failed to search IGDB:', e)
      return []
    }
  }

  return {
    backlog,
    playing,
    history,
    loading,
    error,
    fetchGames,
    createGame,
    updateStatus,
    searchIGDB
  }
})
