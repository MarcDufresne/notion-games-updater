const DATE_PLAYED_SENTINEL = new Date('1970-01-01T00:00:00Z')

export function isReleaseDateSentinel(dateString) {
  if (!dateString) return false
  const date = new Date(dateString)
  const year = date.getFullYear()
  return year === 2099 || year === 2100
}

export function isDatePlayedSentinel(dateString) {
  if (!dateString) return false
  const date = new Date(dateString)
  return date.getTime() === 0 || (date.getFullYear() === 1970 && date.getMonth() === 0 && date.getDate() === 1)
}

export function formatReleaseDate(dateString, options = { year: 'numeric', month: 'long', day: 'numeric' }) {
  if (!dateString || isReleaseDateSentinel(dateString)) {
    return 'TBD'
  }
  const date = new Date(dateString)
  if (date.getFullYear() >= 2099) {
    return 'TBD'
  }

  const year = date.getUTCFullYear()
  const month = date.getUTCMonth()
  const day = date.getUTCDate()

  const utcDate = new Date(year, month, day)
  return utcDate.toLocaleDateString('en-US', options)
}

export function formatDatePlayed(dateString, options = { year: 'numeric', month: 'long', day: 'numeric' }) {
  if (!dateString || isDatePlayedSentinel(dateString)) {
    return null
  }
  const date = new Date(dateString)

  const year = date.getUTCFullYear()
  const month = date.getUTCMonth()
  const day = date.getUTCDate()

  const utcDate = new Date(year, month, day)
  return utcDate.toLocaleDateString('en-US', options)
}
