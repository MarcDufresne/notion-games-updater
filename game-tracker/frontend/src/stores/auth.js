import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authService } from '@/services/auth'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(null)
  const loading = ref(true)
  const error = ref(null)

  // Computed
  const isAuthenticated = computed(() => user.value !== null)
  const userDisplayName = computed(() => user.value?.displayName || 'User')
  const userEmail = computed(() => user.value?.email || '')
  const userPhotoURL = computed(() => user.value?.photoURL || '')

  // Actions
  async function signInWithGoogle() {
    loading.value = true
    error.value = null
    
    try {
      const loggedInUser = await authService.signInWithGoogle()
      user.value = loggedInUser
    } catch (err) {
      error.value = err.message
      console.error('Failed to sign in:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function signOut() {
    loading.value = true
    error.value = null
    
    try {
      await authService.signOut()
      user.value = null
    } catch (err) {
      error.value = err.message
      console.error('Failed to sign out:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  function setUser(newUser) {
    user.value = newUser
    loading.value = false
  }

  function clearError() {
    error.value = null
  }

  // Initialize auth state listener
  function initAuthListener() {
    authService.onAuthStateChanged((newUser) => {
      user.value = newUser
      loading.value = false
    })
  }

  return {
    // State
    user,
    loading,
    error,
    // Computed
    isAuthenticated,
    userDisplayName,
    userEmail,
    userPhotoURL,
    // Actions
    signInWithGoogle,
    signOut,
    setUser,
    clearError,
    initAuthListener
  }
})
