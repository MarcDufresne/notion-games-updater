import { auth } from './firebase'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

async function getAuthHeaders() {
  const user = auth.currentUser
  if (!user) {
    throw new Error('Not authenticated')
  }
  const token = await user.getIdToken()
  return {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
}

export const api = {
  async getGames(view = '') {
    const headers = await getAuthHeaders()
    const url = view ? `${API_URL}/api/v1/games?view=${view}` : `${API_URL}/api/v1/games`
    const response = await fetch(url, { headers })
    if (!response.ok) throw new Error('Failed to fetch games')
    return response.json()
  },

  async createGame(gameData) {
    const headers = await getAuthHeaders()
    const response = await fetch(`${API_URL}/api/v1/games`, {
      method: 'POST',
      headers,
      body: JSON.stringify(gameData)
    })
    if (!response.ok) {
      // Try to extract error message from response
      const errorText = await response.text()
      throw new Error(errorText || 'Failed to create game')
    }
    return response.json()
  },

  async updateGameStatus(gameId, status, datePlayed = null) {
    const headers = await getAuthHeaders()
    const body = { status }
    if (datePlayed) {
      body.date_played = datePlayed
    }
    const response = await fetch(`${API_URL}/api/v1/games/${gameId}/status`, {
      method: 'POST',
      headers,
      body: JSON.stringify(body)
    })
    if (!response.ok) throw new Error('Failed to update game status')
    return response.json()
  },

  async searchGames(query) {
    const headers = await getAuthHeaders()
    const response = await fetch(`${API_URL}/api/v1/search?q=${encodeURIComponent(query)}`, {
      headers
    })
    if (!response.ok) throw new Error('Failed to search games')
    return response.json()
  },

  async deleteGame(gameId) {
    const headers = await getAuthHeaders()
    const response = await fetch(`${API_URL}/api/v1/games/${gameId}`, {
      method: 'DELETE',
      headers
    })
    if (!response.ok) throw new Error('Failed to delete game')
  },

  async updateGameMatch(gameId, igdbId) {
    const headers = await getAuthHeaders()
    const response = await fetch(`${API_URL}/api/v1/games/${gameId}/match`, {
      method: 'POST',
      headers,
      body: JSON.stringify({ igdb_id: igdbId })
    })
    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(errorText || 'Failed to update game match')
    }
    return response.json()
  }
}
