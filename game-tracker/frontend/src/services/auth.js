import { auth, googleProvider } from './firebase'
import { signInWithPopup, signOut, onAuthStateChanged } from 'firebase/auth'

export const authService = {
  // Sign in with Google
  async signInWithGoogle() {
    try {
      const result = await signInWithPopup(auth, googleProvider)
      return result.user
    } catch (error) {
      console.error('Error signing in with Google:', error)
      throw error
    }
  },

  // Sign out
  async signOut() {
    try {
      await signOut(auth)
    } catch (error) {
      console.error('Error signing out:', error)
      throw error
    }
  },

  // Get current user's ID token
  async getIdToken() {
    const user = auth.currentUser
    if (!user) {
      throw new Error('No user logged in')
    }
    return await user.getIdToken()
  },

  // Listen to auth state changes
  onAuthStateChanged(callback) {
    return onAuthStateChanged(auth, callback)
  },

  // Get current user
  getCurrentUser() {
    return auth.currentUser
  }
}
