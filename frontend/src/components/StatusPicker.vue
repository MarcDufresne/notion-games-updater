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

    <!-- Date Picker Modal - Teleported to body -->
    <Teleport to="body">
      <div
        v-if="showDatePicker"
        class="fixed inset-0 z-[10000] flex items-center justify-center bg-black/70 backdrop-blur-sm"
        @click.self="cancelDateSelection"
      >
        <div class="bg-gray-800 rounded-lg shadow-2xl p-6 border border-gray-700 max-w-sm w-full mx-4" @click.stop>
          <h3 class="text-lg font-semibold text-white mb-4">
            When did you {{ selectedStatus === 'Done' ? 'complete' : selectedStatus === 'Abandoned' ? 'abandon' : 'decide not to play' }} this game?
          </h3>

          <div class="mb-6">
            <label class="block text-sm font-medium text-gray-400 mb-2">Date Played</label>
            <input
              ref="datePicker"
              type="date"
              v-model="selectedDate"
              class="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              :max="new Date().toISOString().split('T')[0]"
            />
          </div>

          <div class="flex gap-3">
            <button
              @click="cancelDateSelection"
              class="flex-1 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
            >
              Cancel
            </button>
            <button
              @click="confirmDateSelection"
              class="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
            >
              Confirm
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { isDatePlayedSentinel } from '../lib/dateUtils'

const props = defineProps({
  modelValue: {
    type: String,
    required: true
  },
  datePlayed: {
    type: String,
    default: null
  }
})

const emit = defineEmits(['update:modelValue'])

const isOpen = ref(false)
const showDatePicker = ref(false)
const selectedStatus = ref(null)
const selectedDate = ref(null)
const statusPicker = ref(null)
const statusButton = ref(null)
const dropdown = ref(null)
const datePicker = ref(null)

const dropdownStyle = ref({})

const todoStatuses = ['Backlog', 'Break']
const inProgressStatuses = ['Playing']
const completeStatuses = ['Done', 'Abandoned', "Won't Play"]

const STATUS_COLORS = new Map([
  ['Backlog', 'bg-gray-600 text-gray-100 hover:bg-gray-500 focus:ring-gray-500'],
  ['Break', 'bg-yellow-600 text-yellow-100 hover:bg-yellow-500 focus:ring-yellow-500'],
  ['Playing', 'bg-blue-600 text-blue-100 hover:bg-blue-500 focus:ring-blue-500'],
  ['Done', 'bg-green-600 text-green-100 hover:bg-green-500 focus:ring-green-500'],
  ['Abandoned', 'bg-red-600 text-red-100 hover:bg-red-500 focus:ring-red-500'],
  ["Won't Play", 'bg-gray-700 text-gray-200 hover:bg-gray-600 focus:ring-gray-600 border border-gray-600']
])

const STATUS_TEXT_COLORS = new Map([
  ['Backlog', 'text-gray-300'],
  ['Break', 'text-yellow-300'],
  ['Playing', 'text-blue-300'],
  ['Done', 'text-green-300'],
  ['Abandoned', 'text-red-300'],
  ["Won't Play", 'text-gray-400']
])

function calculateDropdownPosition() {
  if (!statusButton.value) return

  const buttonRect = statusButton.value.getBoundingClientRect()
  const dropdownWidth = buttonRect.width

  // Estimate dropdown height (3 groups with headers + buttons)
  // To-do: 2 items, In Progress: 1 item, Complete: 3 items
  // Header height ~32px, button height ~40px
  const estimatedDropdownHeight = (3 * 32) + (6 * 40) + 20 // headers + buttons + padding

  const spaceBelow = window.innerHeight - buttonRect.bottom
  const spaceAbove = buttonRect.top

  // Check if there's enough space below, otherwise position above
  const shouldPositionAbove = spaceBelow < estimatedDropdownHeight && spaceAbove > spaceBelow

  if (shouldPositionAbove) {
    dropdownStyle.value = {
      bottom: `${window.innerHeight - buttonRect.top + 8}px`,
      left: `${buttonRect.left}px`,
      width: `${dropdownWidth}px`,
      minWidth: '200px'
    }
  } else {
    dropdownStyle.value = {
      top: `${buttonRect.bottom + 8}px`,
      left: `${buttonRect.left}px`,
      width: `${dropdownWidth}px`,
      minWidth: '200px'
    }
  }
}

async function toggleDropdown() {
  isOpen.value = !isOpen.value
  if (isOpen.value) {
    await nextTick()
    calculateDropdownPosition()
  }
}

// Recalculate position on scroll or resize (throttled for performance)
let scrollThrottle = null
let resizeThrottle = null

function handleScroll() {
  if (!isOpen.value) return

  if (scrollThrottle) return

  scrollThrottle = setTimeout(() => {
    calculateDropdownPosition()
    scrollThrottle = null
  }, 16) // ~60fps
}

function handleResize() {
  if (!isOpen.value) return

  if (resizeThrottle) return

  resizeThrottle = setTimeout(() => {
    calculateDropdownPosition()
    resizeThrottle = null
  }, 100)
}

function selectStatus(status) {
  const completedStatuses = ['Done', 'Abandoned', "Won't Play"]

  if (completedStatuses.includes(status)) {
    selectedStatus.value = status
    isOpen.value = false
    showDatePicker.value = true

    if (props.datePlayed && !isDatePlayedSentinel(props.datePlayed)) {
      const date = new Date(props.datePlayed)
      selectedDate.value = date.toISOString().split('T')[0]
    } else {
      const today = new Date().toISOString().split('T')[0]
      selectedDate.value = today
    }
  } else {
    emit('update:modelValue', status)
    isOpen.value = false
  }
}

function confirmDateSelection() {
  if (selectedStatus.value && selectedDate.value) {
    // Convert date string (YYYY-MM-DD) to ISO 8601 datetime string
    // Set time to noon UTC to avoid timezone issues
    const dateObj = new Date(selectedDate.value + 'T12:00:00Z')
    const isoDateString = dateObj.toISOString()
    emit('update:modelValue', selectedStatus.value, isoDateString)
  }
  showDatePicker.value = false
  selectedStatus.value = null
  selectedDate.value = null
}

function cancelDateSelection() {
  showDatePicker.value = false
  selectedStatus.value = null
  selectedDate.value = null
}

function getStatusColor(status) {
  return STATUS_COLORS.get(status) || 'bg-gray-700 text-gray-200 hover:bg-gray-600 focus:ring-gray-600'
}

function getStatusTextColor(status) {
  return STATUS_TEXT_COLORS.get(status) || 'text-gray-300'
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
  if (!isOpen.value) return // Don't process if dropdown isn't open

  // Check if click is inside the status picker button
  if (statusPicker.value && statusPicker.value.contains(event.target)) {
    return
  }

  // Check if click is inside the dropdown (which is teleported to body)
  if (dropdown.value && dropdown.value.contains(event.target)) {
    return
  }

  // Click is outside both - close the dropdown
  isOpen.value = false
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside, true) // Use capture phase
  window.addEventListener('scroll', handleScroll, { capture: true, passive: true })
  window.addEventListener('resize', handleResize, { passive: true })
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside, true) // Use capture phase
  window.removeEventListener('scroll', handleScroll, { capture: true, passive: true })
  window.removeEventListener('resize', handleResize, { passive: true })

  // Clear any pending throttles
  if (scrollThrottle) clearTimeout(scrollThrottle)
  if (resizeThrottle) clearTimeout(resizeThrottle)
})
</script>

<style scoped>
.bg-gray-750 {
  background-color: #2d3748;
}
</style>
