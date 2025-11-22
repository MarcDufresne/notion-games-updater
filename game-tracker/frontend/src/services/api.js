import { authService } from './auth'

const API_BASE_URL = '/api/v1'

async function fetchWithAuth(url, options = {}) {
  try {
    const token = await authService.getIdToken()
    
    const headers = {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
      ...options.headers
    }

    const response = await fetch(url, {
      ...options,
      headers
    })

    if (!response.ok) {
      const error = await response.text()
      throw new Error(error || `HTTP error! status: ${response.status}`)
    }

    return await response.json()
  } catch (error) {
    console.error('API request failed:', error)
    throw error
  }
}

export const api = {
  // Get games with optional view filter
  async getGames(view = null) {
    const url = view ? `${API_BASE_URL}/games?view=${view}` : `${API_BASE_URL}/games`
    return await fetchWithAuth(url)
  },

  // Create or update a game
  async saveGame(game) {
    return await fetchWithAuth(`${API_BASE_URL}/games`, {
      method: 'POST',
      body: JSON.stringify(game)
    })
  },

  // Update game status
  async updateGameStatus(gameId, status, datePlayed = null) {
    const body = { status }
    if (datePlayed) {
      body.date_played = datePlayed
    }
    
    return await fetchWithAuth(`${API_BASE_URL}/games/${gameId}/status`, {
      method: 'POST',
      body: JSON.stringify(body)
    })
  },

  // Health check
  async healthCheck() {
    const response = await fetch(`${API_BASE_URL}/health`)
    return await response.json()
  }
}
