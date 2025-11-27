import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '../lib/api'

const sortByReleaseDate = (games, descending = false) => {
  return games.sort((a, b) => {
    const dateA = a.release_date ? new Date(a.release_date) : new Date(0)
    const dateB = b.release_date ? new Date(b.release_date) : new Date(0)
    return descending ? dateB - dateA : dateA - dateB
  })
}

const shouldAddToCalendar = (game) => {
  if (!game.release_date) return false
  const releaseDate = new Date(game.release_date)
  const oneMonthAgo = new Date()
  oneMonthAgo.setMonth(oneMonthAgo.getMonth() - 1)
  return releaseDate >= oneMonthAgo
}

const addGameToStatusList = (game, lists) => {
  const status = game.status

  if (status === 'Backlog' || status === 'Break') {
    lists.backlog.value.push(game)
    sortByReleaseDate(lists.backlog.value)

    if (shouldAddToCalendar(game)) {
      lists.calendar.value.push(game)
      sortByReleaseDate(lists.calendar.value)
    }
  } else if (status === 'Playing') {
    lists.playing.value.unshift(game)
  } else {
    lists.history.value.push(game)
    // Sort history by date_played descending (most recent first)
    lists.history.value.sort((a, b) => {
      const dateA = a.date_played ? new Date(a.date_played).getTime() : 0
      const dateB = b.date_played ? new Date(b.date_played).getTime() : 0
      return dateB - dateA
    })
  }
}

const removeGameFromAllLists = (gameId, lists) => {
  lists.backlog.value = lists.backlog.value.filter(g => g.id !== gameId)
  lists.playing.value = lists.playing.value.filter(g => g.id !== gameId)
  lists.history.value = lists.history.value.filter(g => g.id !== gameId)
  lists.calendar.value = lists.calendar.value.filter(g => g.id !== gameId)
}

const updateGameInList = (gameId, updatedGame, list) => {
  const index = list.findIndex(g => g.id === gameId)
  if (index !== -1) {
    list[index] = updatedGame
  }
}

export const useGamesStore = defineStore('games', () => {
  const backlog = ref([])
  const playing = ref([])
  const history = ref([])
  const calendar = ref([])
  const all = ref([])
  const loading = ref(false)
  const error = ref(null)

  const lists = { backlog, playing, history, calendar, all }

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
    try {
      const game = await api.createGame(gameData)

      addGameToStatusList(game, lists)

      all.value.push(game)
      sortByReleaseDate(all.value, true)

      return game
    } catch (e) {
      console.error('Failed to create game:', e)
      throw e
    }
  }

  async function updateStatus(gameId, status, datePlayed = null) {
    error.value = null
    try {
      const updatedGame = await api.updateGameStatus(gameId, status, datePlayed)

      removeGameFromAllLists(gameId, lists)
      updateGameInList(gameId, updatedGame, all.value)
      addGameToStatusList(updatedGame, lists)

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

      removeGameFromAllLists(gameId, lists)
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

      // Remove the game from all lists and re-add it to properly sort it
      // based on its updated metadata (like release date)
      removeGameFromAllLists(gameId, lists)

      // Re-add to appropriate status lists with proper sorting
      addGameToStatusList(updatedGame, lists)

      // Update in the 'all' list - remove old entry and add with proper sorting
      all.value = all.value.filter(g => g.id !== gameId)
      all.value.push(updatedGame)
      sortByReleaseDate(all.value, true)

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
