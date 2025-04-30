<template>
  <div>
    <!-- Loading state -->
    <div v-if="loading" class="d-flex justify-center my-8">
      <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
    </div>

    <!-- Error state -->
    <v-alert v-else-if="error" type="error" class="my-4">
      {{ error }}
    </v-alert>

    <!-- Session details -->
    <template v-else>
      <div v-if="!sessionsStore.currentSession">
        <v-alert type="info" class="my-4">
          Session not found. It may have been deleted.
          <div class="mt-4">
            <v-btn color="primary" to="/sessions">
              Back to Sessions
            </v-btn>
          </div>
        </v-alert>
      </div>
      <template v-else>
        <!-- Header section with actions -->
        <v-row class="mb-4">
          <v-col cols="12" sm="8">
            <v-btn variant="text" to="/sessions" prepend-icon="mdi-arrow-left" class="mb-2">
              Back to Sessions
            </v-btn>
            <h1 class="text-h3">Practice Session</h1>
            <p class="text-body-1">
              {{ formatDate(sessionsStore.currentSession.startTime) }}
            </p>
          </v-col>
          <v-col cols="12" sm="4" class="d-flex justify-end align-center">
            <v-btn color="error" variant="outlined" prepend-icon="mdi-delete" class="mr-2" @click="confirmDelete">
              Delete
            </v-btn>
            <v-btn color="primary" variant="outlined" prepend-icon="mdi-pencil" class="mr-2" @click="openEditDialog">
              Edit
            </v-btn>
          </v-col>
        </v-row>

        <!-- Session details card -->
        <v-card class="mb-6">
          <v-card-text>
            <v-row>
              <v-col cols="12" md="6">
                <div class="text-body-1 font-weight-medium mb-1">Start Time</div>
                <div class="text-body-2 mb-4">{{ formatDateTime(sessionsStore.currentSession.startTime) }}</div>

                <div class="text-body-1 font-weight-medium mb-1">End Time</div>
                <div class="text-body-2 mb-4">{{ formatDateTime(sessionsStore.currentSession.endTime) }}</div>

                <div class="text-body-1 font-weight-medium mb-1">Duration</div>
                <div class="text-body-2 mb-4">{{ formatDuration(sessionsStore.currentSession.startTime,
                  sessionsStore.currentSession.endTime) }}</div>

                <div class="text-body-1 font-weight-medium mb-1">Notes</div>
                <div class="text-body-2 mb-4">{{ sessionsStore.currentSession.notes || 'No notes provided' }}</div>
              </v-col>

              <v-col cols="12" md="6">
                <div class="text-body-1 font-weight-medium mb-2">Exercises</div>
                <v-list
                  v-if="sessionsStore.currentSession.exercises && sessionsStore.currentSession.exercises.length > 0">
                  <v-list-item v-for="exercise in sessionsStore.currentSession.exercises" :key="exercise.id">
                    <v-list-item-title>{{ exercise.exercise ? exercise.exercise.name : `Exercise
                      #${exercise.exerciseId}` }}</v-list-item-title>
                    <v-list-item-subtitle>
                      <span v-if="exercise.bpms">{{ exercise.bpms.join(', ') }} BPM</span>
                      <span v-if="exercise.timeSignature" class="ml-2">{{ exercise.timeSignature }}</span>
                      <span v-if="exercise.notes" class="d-block">{{ exercise.notes }}</span>
                    </v-list-item-subtitle>
                  </v-list-item>
                </v-list>
                <div v-else class="text-body-2 mb-4 text-grey">
                  No exercises were recorded in this session
                </div>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>

        <!-- Exercise Details Cards -->
        <h2 class="text-h5 mb-4"
          v-if="sessionsStore.currentSession.exercises && sessionsStore.currentSession.exercises.length > 0">
          Exercise Details
        </h2>
        <v-row v-if="sessionsStore.currentSession.exercises && sessionsStore.currentSession.exercises.length > 0">
          <v-col v-for="(exercise, index) in sessionsStore.currentSession.exercises" :key="index" cols="12" md="6"
            lg="4">
            <v-card variant="outlined" class="mb-4">
              <v-card-title>
                {{ exercise.exercise ? exercise.exercise.name : `Exercise #${exercise.exerciseId}` }}
              </v-card-title>
              <v-card-text>
                <div class="d-flex mb-2">
                  <div class="text-body-2 font-weight-medium mr-4">
                    <v-icon icon="mdi-metronome" size="small" class="mr-1"></v-icon>
                    {{ exercise.bpms.join(', ') || 'N/A' }} BPM
                  </div>
                  <div class="text-body-2 font-weight-medium">
                    <v-icon icon="mdi-music" size="small" class="mr-1"></v-icon>
                    {{ exercise.timeSignature || '4/4' }}
                  </div>
                </div>

                <div class="text-body-2 font-weight-medium mb-1">Duration</div>
                <div class="text-body-2 mb-2">{{ exercisesStore.getExerciseDuration(exercise) }}</div>

                <div class="text-body-2 font-weight-medium mb-1">Notes</div>
                <div class="text-body-2">{{ exercise.notes || 'No notes provided' }}</div>

                <v-divider class="my-3"></v-divider>

                <v-btn variant="text" :to="`/exercises/${exercise.exerciseId}`" color="primary" size="small">
                  View Exercise Details
                </v-btn>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </template>
    </template>

    <!-- Delete Confirmation Dialog -->
    <v-dialog v-model="deleteDialog" max-width="500">
      <v-card>
        <v-card-title class="text-h5">Delete Session</v-card-title>
        <v-card-text>
          Are you sure you want to delete this practice session? This action cannot be undone.
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="grey-darken-1" variant="text" @click="deleteDialog = false">Cancel</v-btn>
          <v-btn color="error" variant="flat" @click="deleteSession" :loading="deleteLoading">
            Delete
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="editDialog" max-width="600" persistent>
      <v-card>
        <v-card-title class="text-h5">Edit Session</v-card-title>
        <v-card-text>
          <v-form ref="editForm" v-model="editFormValid">
            <v-row>
              <v-col cols="12">
                <v-textarea v-model="editSessionData.notes" label="Session Notes" variant="outlined" rows="4"
                  placeholder="Add notes about this practice session (optional)" class="mb-4"></v-textarea>
              </v-col>

              <!-- Start Date/Time -->
              <v-col cols="12" md="6">
                <v-dialog v-model="startDateMenu" :close-on-content-click="false" width="auto">
                  <template v-slot:activator="{ props }">
                    <v-text-field v-model="formattedStartDate" label="Start Date" prepend-inner-icon="mdi-calendar"
                      variant="outlined" readonly v-bind="props"></v-text-field>
                  </template>
                  <v-date-picker v-model="startDate">
                    <v-spacer></v-spacer>
                    <v-btn color="primary" variant="text" @click="startDateMenu = false">Close</v-btn>
                  </v-date-picker>
                </v-dialog>
              </v-col>

              <v-col cols="12" md="6">
                <v-text-field v-model="startTime" label="Start Time (HH:MM)" prepend-inner-icon="mdi-clock-outline"
                  variant="outlined" placeholder="e.g. 14:30" hint="24-hour format" persistent-hint></v-text-field>
              </v-col>

              <!-- End Date/Time -->
              <v-col cols="12" md="6">
                <v-dialog v-model="endDateMenu" :close-on-content-click="false" width="auto">
                  <template v-slot:activator="{ props }">
                    <v-text-field v-model="formattedEndDate" label="End Date" prepend-inner-icon="mdi-calendar"
                      variant="outlined" readonly v-bind="props"></v-text-field>
                  </template>
                  <v-date-picker v-model="endDate">
                    <v-spacer></v-spacer>
                    <v-btn color="primary" variant="text" @click="endDateMenu = false">Close</v-btn>
                  </v-date-picker>
                </v-dialog>
              </v-col>

              <v-col cols="12" md="6">
                <v-text-field v-model="endTime" label="End Time (HH:MM)" prepend-inner-icon="mdi-clock-outline"
                  variant="outlined" placeholder="e.g. 15:45" hint="24-hour format" persistent-hint></v-text-field>
              </v-col>

              <!-- Duration (calculated, read-only) -->
              <v-col cols="12">
                <v-text-field :model-value="calculatedDuration" label="Duration" variant="outlined" readonly
                  disabled></v-text-field>
              </v-col>
            </v-row>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="grey-darken-1" variant="text" @click="cancelEdit">
            Cancel
          </v-btn>
          <v-btn color="primary" variant="flat" @click="saveSessionEdit" :loading="updatingSession">
            Save
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useSessionsStore } from '@/stores/sessions.js'
import { useAppStore } from '@/stores/app.js'
import { useCategoriesStore } from '@/stores/categories.js'
import { useExercisesStore } from '@/stores/exercises.js'

const router = useRouter()
const route = useRoute()
const sessionsStore = useSessionsStore()
const appStore = useAppStore()
const categoriesStore = useCategoriesStore()
const exercisesStore = useExercisesStore()

// Data
const loading = ref(true)
const error = ref(null)
const deleteDialog = ref(false)
const deleteLoading = ref(false)

// Computed
const sessionId = computed(() => parseInt(route.params.id))

// Methods
function formatDate(dateString) {
  return appStore.formatDate(dateString)
}

function formatDateTime(dateString) {
  return appStore.formatDateTime(dateString)
}

function formatDuration(startTime, endTime) {
  return appStore.formatDuration(startTime, endTime)
}

function getCategoryById(categoryId) {
  return categoriesStore.categoryById(categoryId)
}

function editSession() {
  // Navigate to edit session page (to be implemented)
  appStore.showInfoMessage('Session editing will be implemented in a future update')
}

function confirmDelete() {
  deleteDialog.value = true
}

async function deleteSession() {
  deleteLoading.value = true
  try {
    await sessionsStore.deleteSession(sessionId.value)
    appStore.showSuccessMessage('Practice session deleted successfully')
    deleteDialog.value = false
    router.push({ name: 'sessions' })
  } catch (error) {
    appStore.showErrorMessage(`Error deleting session: ${error.message}`)
  } finally {
    deleteLoading.value = false
  }
}

// Load data
async function loadData() {
  loading.value = true
  error.value = null

  try {
    await sessionsStore.fetchSession(sessionId.value)
  } catch (err) {
    console.error(`Error loading session: ${err.message}`, err)
    error.value = `Error loading session: ${err.message}`
  } finally {
    loading.value = false
  }
}

// Editing
const editDialog = ref(false)
const editForm = ref(null)
const editFormValid = ref(true)
const editSessionData = ref({
  notes: '',
  startTime: null,
  endTime: null
})
const updatingSession = ref(false)

// For date and time inputs
const startDateMenu = ref(false)
const endDateMenu = ref(false)
const startDate = ref('')
const endDate = ref('')
const startTime = ref('')
const endTime = ref('')

// Formatted date displays
const formattedStartDate = computed(() => startDate.value || '')
const formattedEndDate = computed(() => endDate.value || '')

// Calculate session duration for display
const calculatedDuration = computed(() => {
  if (!startDate.value || !startTime.value || !endDate.value || !endTime.value) {
    return ''
  }

  try {
    const startDateTime = combineDateAndTime(startDate.value, startTime.value)
    const endDateTime = combineDateAndTime(endDate.value, endTime.value)

    if (!startDateTime || !endDateTime) return ''

    return formatDuration(startDateTime, endDateTime)
  } catch (error) {
    return 'Invalid date/time'
  }
})

// Methods for the edit session functionality
function openEditDialog() {
  // Initialize form with current session data
  if (sessionsStore.currentSession) {
    editSessionData.value = {
      notes: sessionsStore.currentSession.notes || '',
      startTime: sessionsStore.currentSession.startTime,
      endTime: sessionsStore.currentSession.endTime
    }

    // Set date and time inputs to current values
    if (sessionsStore.currentSession.startTime) {
      const date = new Date(sessionsStore.currentSession.startTime)
      startDate.value = date.toISOString().split('T')[0]
      startTime.value = `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
    }

    if (sessionsStore.currentSession.endTime) {
      const date = new Date(sessionsStore.currentSession.endTime)
      endDate.value = date.toISOString().split('T')[0]
      endTime.value = `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
    }
  }
  editDialog.value = true
}

function cancelEdit() {
  editDialog.value = false
}

function combineDateAndTime(dateString, timeString) {
  if (!dateString || !timeString) return null

  // Validate time format (HH:MM)
  const timeRegex = /^([0-1]?[0-9]|2[0-3]):([0-5][0-9])$/
  if (!timeRegex.test(timeString)) {
    return null
  }

  // Parse the date and time strings
  const date = new Date(dateString)
  const [hours, minutes] = timeString.split(':').map(Number)

  // Set the time on the date object
  date.setHours(hours, minutes, 0, 0)

  return date.toISOString()
}

async function saveSessionEdit() {
  if (!editFormValid.value) return

  updatingSession.value = true
  try {
    // Validate time format
    const timeRegex = /^([0-1]?[0-9]|2[0-3]):([0-5][0-9])$/
    if (!timeRegex.test(startTime.value) || !timeRegex.test(endTime.value)) {
      appStore.showErrorMessage('Please enter valid time in HH:MM format')
      updatingSession.value = false
      return
    }

    // Combine date and time values into ISO strings
    const combinedStartTime = combineDateAndTime(startDate.value, startTime.value)
    const combinedEndTime = combineDateAndTime(endDate.value, endTime.value)

    if (!combinedStartTime || !combinedEndTime) {
      appStore.showErrorMessage('Please enter valid date and time values')
      updatingSession.value = false
      return
    }

    // Validate that start time is before end time
    const startTimeDate = new Date(combinedStartTime)
    const endTimeDate = new Date(combinedEndTime)

    if (startTimeDate >= endTimeDate) {
      appStore.showErrorMessage('Start time must be before end time')
      updatingSession.value = false
      return
    }

    // Prepare the update payload
    const sessionData = {
      notes: editSessionData.value.notes,
      start_time: combinedStartTime,
      end_time: combinedEndTime
    }

    // Update the session
    await sessionsStore.updateSession(
      sessionId.value,
      sessionData,
      'notes,startTime,endTime'
    )

    appStore.showSuccessMessage('Session updated successfully')
    editDialog.value = false

    // Refresh the session data
    await loadData()
  } catch (error) {
    appStore.showErrorMessage(`Error updating session: ${error.message}`)
  } finally {
    updatingSession.value = false
  }
}

// Watch for route param changes to reload data
watch(() => route.params.id, () => {
  loadData()
})

// Lifecycle
onMounted(() => {
  loadData()
})
</script>

<style scoped>
.font-weight-medium {
  font-weight: 500;
}
</style>
