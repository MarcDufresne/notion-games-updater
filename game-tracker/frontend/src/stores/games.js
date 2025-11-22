import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useGameStore = defineStore('games', () => {
  const backlog = ref([])
  const playing = ref([])
  const history = ref([])
  const loading = ref(false)
  const error = ref(null)

  // Mock token for now (should come from Firebase Auth)
  const token = "MOCK_TOKEN"

  async function fetchGames(view) {
    loading.value = true
    error.value = null
    try {
      const response = await fetch(`/api/v1/games?view=${view}`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })
      if (!response.ok) throw new Error('Failed to fetch games')
      const data = await response.json()

      if (view === 'backlog') backlog.value = data || []
      else if (view === 'playing') playing.value = data || []
      else if (view === 'history') history.value = data || []

    } catch (err) {
      error.value = err.message
      // Fallback for dev without backend
      console.warn("Using fallback data for " + view)
      if (view === 'playing') {
        playing.value = [
            { id: '1', title: 'Zelda: TOTK', status: 'Playing', cover_url: 'https://via.placeholder.com/150', platforms: ['Switch'] },
            { id: '2', title: 'Baldur\'s Gate 3', status: 'Playing', cover_url: 'https://via.placeholder.com/150', platforms: ['PC'] }
        ]
      } else if (view === 'backlog') {
          backlog.value = [
              { id: '3', title: 'Hollow Knight', status: 'Backlog', cover_url: 'https://via.placeholder.com/150', platforms: ['PC'] },
              { id: '4', title: 'Elden Ring', status: 'Break', cover_url: 'https://via.placeholder.com/150', platforms: ['PS5'] }
          ]
      } else if (view === 'history') {
          history.value = [
              { id: '5', title: 'Celeste', status: 'Done', date_played: '2023-01-15T00:00:00Z', cover_url: 'https://via.placeholder.com/150', platforms: ['Switch'] }
          ]
      }

    } finally {
      loading.value = false
    }
  }

  async function updateStatus(id, newStatus) {
    // Implementation for updating status
    console.log(`Updating game ${id} to ${newStatus}`)
  }

  return { backlog, playing, history, loading, error, fetchGames, updateStatus }
})
