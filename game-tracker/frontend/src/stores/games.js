import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '../lib/api'

export const useGamesStore = defineStore('games', () => {
  const backlog = ref([])
  const playing = ref([])
  const history = ref([])
  const calendar = ref([])
  const all = ref([])
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
        case 'calendar':
          calendar.value = games
          break
        case 'all':
          all.value = games
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
    // Don't set loading to true - this causes jarring UI
    // Don't set error.value - let calling components handle errors with toasts
    try {
      const game = await api.createGame(gameData)

      // Add to appropriate list based on status with proper sorting
      if (game.status === 'Backlog' || game.status === 'Break') {
        // Backlog is sorted by release_date ASC
        backlog.value.push(game)
        backlog.value.sort((a, b) => {
          const dateA = a.release_date ? new Date(a.release_date) : new Date(0)
          const dateB = b.release_date ? new Date(b.release_date) : new Date(0)
          return dateA - dateB
        })

        // Also add to calendar if it has a release date in the relevant range
        if (game.release_date) {
          const releaseDate = new Date(game.release_date)
          const oneMonthAgo = new Date()
          oneMonthAgo.setMonth(oneMonthAgo.getMonth() - 1)
          if (releaseDate >= oneMonthAgo) {
            calendar.value.push(game)
            calendar.value.sort((a, b) => {
              const dateA = new Date(a.release_date)
              const dateB = new Date(b.release_date)
              return dateA - dateB
            })
          }
        }
      } else if (game.status === 'Playing') {
        // Playing is sorted by updated_at DESC
        playing.value.unshift(game) // Add to beginning since it's newest
      } else {
        // History is sorted by date_played DESC
        history.value.unshift(game) // Add to beginning since it's newest
      }

      // Also add to 'all' array sorted by release_date DESC (newest first)
      all.value.push(game)
      all.value.sort((a, b) => {
        const dateA = a.release_date ? new Date(a.release_date) : new Date(0)
        const dateB = b.release_date ? new Date(b.release_date) : new Date(0)
        return dateB - dateA // Descending
      })

      return game
    } catch (e) {
      // Don't set error.value - calling component will show toast
      console.error('Failed to create game:', e)
      throw e
    }
  }

  async function updateStatus(gameId, status, datePlayed = null) {
    // Don't set loading to true - this causes jarring UI
    error.value = null
    try {
      const updatedGame = await api.updateGameStatus(gameId, status, datePlayed)

      // Remove from all lists
      backlog.value = backlog.value.filter(g => g.id !== gameId)
      playing.value = playing.value.filter(g => g.id !== gameId)
      history.value = history.value.filter(g => g.id !== gameId)
      calendar.value = calendar.value.filter(g => g.id !== gameId)

      // Update in 'all' array (in place to maintain sort order)
      const allIndex = all.value.findIndex(g => g.id === gameId)
      if (allIndex !== -1) {
        all.value[allIndex] = updatedGame
      }

      // Add to appropriate list with proper sorting
      if (status === 'Backlog' || status === 'Break') {
        // Backlog is sorted by release_date ASC
        backlog.value.push(updatedGame)
        backlog.value.sort((a, b) => {
          const dateA = a.release_date ? new Date(a.release_date) : new Date(0)
          const dateB = b.release_date ? new Date(b.release_date) : new Date(0)
          return dateA - dateB
        })

        // Also add to calendar if it has a release date in the relevant range
        if (updatedGame.release_date) {
          const releaseDate = new Date(updatedGame.release_date)
          const oneMonthAgo = new Date()
          oneMonthAgo.setMonth(oneMonthAgo.getMonth() - 1)
          if (releaseDate >= oneMonthAgo) {
            calendar.value.push(updatedGame)
            calendar.value.sort((a, b) => {
              const dateA = new Date(a.release_date)
              const dateB = new Date(b.release_date)
              return dateA - dateB
            })
          }
        }
      } else if (status === 'Playing') {
        // Playing is sorted by updated_at DESC (newest first)
        playing.value.unshift(updatedGame)
      } else {
        // History is sorted by date_played DESC (newest first)
        history.value.unshift(updatedGame)
      }

      return updatedGame
    } catch (e) {
      error.value = e.message
      console.error('Failed to update game status:', e)
      throw e
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

  async function deleteGame(gameId) {
    error.value = null
    try {
      await api.deleteGame(gameId)

      // Remove from all lists
      backlog.value = backlog.value.filter(g => g.id !== gameId)
      playing.value = playing.value.filter(g => g.id !== gameId)
      history.value = history.value.filter(g => g.id !== gameId)
      calendar.value = calendar.value.filter(g => g.id !== gameId)
      all.value = all.value.filter(g => g.id !== gameId)

      console.log('Game deleted successfully')
    } catch (e) {
      error.value = e.message
      console.error('Failed to delete game:', e)
      throw e
    }
  }

  async function updateGameMatch(gameId, igdbId) {
    try {
      const updatedGame = await api.updateGameMatch(gameId, igdbId)

      // Update game in all lists where it appears
      const updateInList = (list) => {
        const index = list.findIndex(g => g.id === gameId)
        if (index !== -1) {
          list[index] = updatedGame
        }
      }

      updateInList(backlog.value)
      updateInList(playing.value)
      updateInList(history.value)
      updateInList(calendar.value)
      updateInList(all.value)

      console.log('Game match updated successfully')
      return updatedGame
    } catch (e) {
      console.error('Failed to update game match:', e)
      throw e
    }
  }

  return {
    backlog,
    playing,
    history,
    calendar,
    all,
    loading,
    error,
    fetchGames,
    createGame,
    updateStatus,
    searchIGDB,
    deleteGame,
    updateGameMatch
  }
})
