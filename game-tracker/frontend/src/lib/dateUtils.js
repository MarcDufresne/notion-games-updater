// Sentinel date values used in the database
const RELEASE_DATE_SENTINEL = new Date('2100-01-01T00:00:00Z')
const DATE_PLAYED_SENTINEL = new Date('1970-01-01T00:00:00Z')

/**
 * Check if a date string is the release date sentinel value (2100-01-01)
 * Note: Due to timezone conversions, the date might appear as 2099-12-31 in some timezones
 * @param {string|Date} dateString - The date to check
 * @returns {boolean} - True if the date is the sentinel value
 */
export function isReleaseDateSentinel(dateString) {
  if (!dateString) return false
  const date = new Date(dateString)
  const year = date.getFullYear()
  // Check for 2099 or 2100 to handle timezone conversions
  return year === 2099 || year === 2100
}

/**
 * Check if a date string is the date played sentinel value (1970-01-01)
 * @param {string|Date} dateString - The date to check
 * @returns {boolean} - True if the date is the sentinel value
 */
export function isDatePlayedSentinel(dateString) {
  if (!dateString) return false
  const date = new Date(dateString)
  return date.getTime() === 0 || (date.getFullYear() === 1970 && date.getMonth() === 0 && date.getDate() === 1)
}

/**
 * Format a release date, showing "TBD" for sentinel values
 * Uses UTC to avoid timezone conversion issues (games release on a date, not at a time)
 * @param {string|Date} dateString - The date to format
 * @param {Object} options - Intl.DateTimeFormat options
 * @returns {string} - Formatted date or "TBD"
 */
export function formatReleaseDate(dateString, options = { year: 'numeric', month: 'long', day: 'numeric' }) {
  if (!dateString || isReleaseDateSentinel(dateString)) {
    return 'TBD'
  }
  const date = new Date(dateString)
  // Additional check: if year is 2099 or 2100, treat as TBD
  if (date.getFullYear() >= 2099) {
    return 'TBD'
  }

  // Format using UTC to avoid timezone conversion issues
  // Extract year, month, day in UTC
  const year = date.getUTCFullYear()
  const month = date.getUTCMonth()
  const day = date.getUTCDate()

  // Create a date object in local timezone but with UTC date values
  // This ensures the displayed date matches the actual release date regardless of timezone
  const utcDate = new Date(year, month, day)

  return utcDate.toLocaleDateString('en-US', options)
}

/**
 * Format a date played, showing nothing for sentinel values
 * Uses UTC to avoid timezone conversion issues (played on a date, not at a time)
 * @param {string|Date} dateString - The date to format
 * @param {Object} options - Intl.DateTimeFormat options
 * @returns {string|null} - Formatted date or null
 */
export function formatDatePlayed(dateString, options = { year: 'numeric', month: 'long', day: 'numeric' }) {
  if (!dateString || isDatePlayedSentinel(dateString)) {
    return null
  }
  const date = new Date(dateString)

  // Format using UTC to avoid timezone conversion issues
  // Extract year, month, day in UTC
  const year = date.getUTCFullYear()
  const month = date.getUTCMonth()
  const day = date.getUTCDate()

  // Create a date object in local timezone but with UTC date values
  const utcDate = new Date(year, month, day)

  return utcDate.toLocaleDateString('en-US', options)
}
