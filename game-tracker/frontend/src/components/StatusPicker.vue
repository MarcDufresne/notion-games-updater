<template>
  <div class="relative" ref="statusPicker">
    <!-- Status Badge (clickable) -->
    <button
      ref="statusButton"
      @click="toggleDropdown"
      :class="[
        'text-sm rounded-full px-4 py-2 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-800 w-full font-medium transition-all cursor-pointer text-center flex items-center justify-center gap-2',
        getStatusColor(modelValue)
      ]"
    >
      <!-- Status indicator circle -->
      <div v-html="getStatusIcon(modelValue)" class="w-3 h-3 flex-shrink-0"></div>
      <span>{{ modelValue }}</span>
      <!-- Dropdown arrow -->
      <svg :class="['w-4 h-4 transition-transform flex-shrink-0', isOpen ? 'rotate-180' : '']" fill="currentColor" viewBox="0 0 20 20">
        <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
      </svg>
    </button>

    <!-- Custom Dropdown - Teleported to body -->
    <Teleport to="body">
      <div
        v-if="isOpen"
        ref="dropdown"
        :style="dropdownStyle"
        class="fixed z-[9999] bg-gray-800 border border-gray-600 rounded-lg shadow-2xl overflow-hidden"
      >
      <!-- To-do Group -->
      <div class="border-b border-gray-700">
        <div class="px-3 py-2 text-xs font-semibold text-gray-400 uppercase tracking-wide bg-gray-750">
          To-do
        </div>
        <button
          v-for="status in todoStatuses"
          :key="status"
          @click="selectStatus(status)"
          :class="[
            'w-full px-4 py-2 text-left text-sm font-medium transition-colors flex items-center gap-2',
            modelValue === status ? 'bg-gray-700' : 'hover:bg-gray-700/50'
          ]"
        >
          <div v-html="getStatusIcon(status)" class="w-3 h-3 flex-shrink-0" :class="getStatusTextColor(status)"></div>
          <span :class="getStatusTextColor(status)">{{ status }}</span>
        </button>
      </div>

      <!-- In Progress Group -->
      <div class="border-b border-gray-700">
        <div class="px-3 py-2 text-xs font-semibold text-gray-400 uppercase tracking-wide bg-gray-750">
          In Progress
        </div>
        <button
          v-for="status in inProgressStatuses"
          :key="status"
          @click="selectStatus(status)"
          :class="[
            'w-full px-4 py-2 text-left text-sm font-medium transition-colors flex items-center gap-2',
            modelValue === status ? 'bg-gray-700' : 'hover:bg-gray-700/50'
          ]"
        >
          <div v-html="getStatusIcon(status)" class="w-3 h-3 flex-shrink-0" :class="getStatusTextColor(status)"></div>
          <span :class="getStatusTextColor(status)">{{ status }}</span>
        </button>
      </div>

      <!-- Complete Group -->
      <div>
        <div class="px-3 py-2 text-xs font-semibold text-gray-400 uppercase tracking-wide bg-gray-750">
          Complete
        </div>
        <button
          v-for="status in completeStatuses"
          :key="status"
          @click="selectStatus(status)"
          :class="[
            'w-full px-4 py-2 text-left text-sm font-medium transition-colors flex items-center gap-2',
            modelValue === status ? 'bg-gray-700' : 'hover:bg-gray-700/50'
          ]"
        >
          <div v-html="getStatusIcon(status)" class="w-3 h-3 flex-shrink-0" :class="getStatusTextColor(status)"></div>
          <span :class="getStatusTextColor(status)">{{ status }}</span>
        </button>
      </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['update:modelValue'])

const isOpen = ref(false)
const statusPicker = ref(null)
const statusButton = ref(null)
const dropdown = ref(null)

const dropdownStyle = ref({})

const todoStatuses = ['Backlog', 'Break']
const inProgressStatuses = ['Playing']
const completeStatuses = ['Done', 'Abandoned', "Won't Play"]

function calculateDropdownPosition() {
  if (!statusButton.value) return

  const buttonRect = statusButton.value.getBoundingClientRect()
  const dropdownWidth = buttonRect.width

  dropdownStyle.value = {
    top: `${buttonRect.bottom + 8}px`,
    left: `${buttonRect.left}px`,
    width: `${dropdownWidth}px`,
    minWidth: '200px'
  }
}

async function toggleDropdown() {
  isOpen.value = !isOpen.value
  if (isOpen.value) {
    await nextTick()
    calculateDropdownPosition()
  }
}

// Recalculate position on scroll or resize
function handleScroll() {
  if (isOpen.value) {
    calculateDropdownPosition()
  }
}

function handleResize() {
  if (isOpen.value) {
    calculateDropdownPosition()
  }
}

function selectStatus(status) {
  emit('update:modelValue', status)
  isOpen.value = false
}

function getStatusColor(status) {
  switch (status) {
    case 'Backlog':
      return 'bg-gray-600 text-gray-100 hover:bg-gray-500 focus:ring-gray-500'
    case 'Break':
      return 'bg-yellow-600 text-yellow-100 hover:bg-yellow-500 focus:ring-yellow-500'
    case 'Playing':
      return 'bg-blue-600 text-blue-100 hover:bg-blue-500 focus:ring-blue-500'
    case 'Done':
      return 'bg-green-600 text-green-100 hover:bg-green-500 focus:ring-green-500'
    case 'Abandoned':
      return 'bg-red-600 text-red-100 hover:bg-red-500 focus:ring-red-500'
    case "Won't Play":
      return 'bg-gray-800 text-gray-300 hover:bg-gray-700 focus:ring-gray-700'
    default:
      return 'bg-gray-700 text-gray-200 hover:bg-gray-600 focus:ring-gray-600'
  }
}

function getStatusTextColor(status) {
  switch (status) {
    case 'Backlog':
      return 'text-gray-300'
    case 'Break':
      return 'text-yellow-300'
    case 'Playing':
      return 'text-blue-300'
    case 'Done':
      return 'text-green-300'
    case 'Abandoned':
      return 'text-red-300'
    case "Won't Play":
      return 'text-gray-400'
    default:
      return 'text-gray-300'
  }
}

function getStatusIcon(status) {
  // To-do statuses (Backlog, Break) - Dotted circle
  if (status === 'Backlog' || status === 'Break') {
    return `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-dasharray="2,2">
      <circle cx="12" cy="12" r="10"/>
    </svg>`
  }

  // In progress status (Playing) - Half-filled circle
  if (status === 'Playing') {
    return `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <circle cx="12" cy="12" r="10"/>
      <path d="M12 2 A10 10 0 0 1 12 22 Z" fill="currentColor"/>
    </svg>`
  }

  // Complete statuses (Done, Abandoned, Won't Play) - Filled circle
  return `<svg viewBox="0 0 24 24" fill="currentColor">
    <circle cx="12" cy="12" r="10"/>
  </svg>`
}

// Close dropdown when clicking outside
function handleClickOutside(event) {
  if (statusPicker.value && !statusPicker.value.contains(event.target)) {
    // Also check if click is inside the dropdown (which is teleported to body)
    if (dropdown.value && !dropdown.value.contains(event.target)) {
      isOpen.value = false
    }
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  window.addEventListener('scroll', handleScroll, true) // Use capture phase
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  window.removeEventListener('scroll', handleScroll, true)
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.bg-gray-750 {
  background-color: #2d3748;
}
</style>
