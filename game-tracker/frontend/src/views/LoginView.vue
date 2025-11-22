<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          Game Tracker
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          Track your gaming library
        </p>
      </div>
      
      <div class="mt-8 space-y-6">
        <button
          @click="handleGoogleSignIn"
          :disabled="loading"
          class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
        >
          <span v-if="loading">Signing in...</span>
          <span v-else>Sign in with Google</span>
        </button>
        
        <p v-if="error" class="text-center text-sm text-red-600">
          {{ error }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const error = ref(null)

async function handleGoogleSignIn() {
  loading.value = true
  error.value = null
  
  try {
    await authStore.signInWithGoogle()
    router.push('/backlog')
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}
</script>
