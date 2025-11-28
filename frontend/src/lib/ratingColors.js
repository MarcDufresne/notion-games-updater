/**
 * Get Metacritic-style color classes for a rating
 * @param {number} rating - Rating value (0-100)
 * @returns {object} Object with bg, text, and border color classes
 */
export function getRatingColor(rating) {
  if (rating >= 75) {
    return {
      bg: 'bg-green-600',
      text: 'text-white',
      border: 'border-green-500'
    }
  } else if (rating >= 50) {
    return {
      bg: 'bg-yellow-600',
      text: 'text-white',
      border: 'border-yellow-500'
    }
  } else {
    return {
      bg: 'bg-red-600',
      text: 'text-white',
      border: 'border-red-500'
    }
  }
}
