<template>
  <div>
    <v-row class="mb-4">
      <v-col cols="12" sm="8">
        <h1 class="text-h3">New Practice Session</h1>
        <p class="text-body-1">Start tracking your practice time</p>
      </v-col>
      <v-col cols="12" sm="4" class="d-flex justify-end align-center">
        <v-btn
          v-if="!sessionStarted"
          color="primary"
          prepend-icon="mdi-play"
          :loading="startingSession"
          @click="startSession"
        >
          Start Session
        </v-btn>
        <v-btn
          v-else
          color="success"
          prepend-icon="mdi-check"
          :loading="finishingSession"
          @click="finishSession"
        >
          Finish Session
        </v-btn>
      </v-col>
    </v-row>

    <v-card class="mb-4">
      <v-card-text>
        <v-row>
          <v-col cols="12" md="6">
            <div class="text-h6 mb-2">Session Details</div>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  v-model="notes"
                  label="Session Notes"
                  variant="outlined"
                  rows="2"
                  placeholder="Add notes about this practice session (optional)"
                  :disabled="finishingSession"
                ></v-text-field>
              </v-col>
            </v-row>

            <div class="text-subtitle-1 font-weight-medium mb-1">
              Start Time:
              <span v-if="sessionStarted" class="font-weight-regular">
                {{ formatDateTime(startTime) }}
              </span>
              <span v-else class="font-weight-regular text-grey">
                Not started yet
              </span>
            </div>

            <div class="text-subtitle-1 font-weight-medium mb-1">
              Current Duration:
              <span class="font-weight-regular" v-if="sessionStarted">
                {{ sessionDuration }}
              </span>
              <span v-else class="font-weight-regular text-grey">
                0:00:00
              </span>
            </div>
          </v-col>

          <v-col cols="12" md="6">
            <div class="text-h6 mb-2">Exercises in this Session</div>
            <v-skeleton-loader
              v-if="loadingSessionExercises"
              type="list-item-two-line"
              :loading="loadingSessionExercises"
            ></v-skeleton-loader>

            <div v-else-if="sessionExercises.length === 0" class="text-center pa-4">
              <v-icon icon="mdi-music-note-off" size="64" color="grey-lighten-1" class="mb-2"></v-icon>
              <p class="text-body-2 text-grey">No exercises added yet</p>
              <v-btn
                color="primary"
                prepend-icon="mdi-plus"
                variant="text"
                @click="openAddExerciseDialog"
                :disabled="!sessionStarted"
              >
                Add Exercises
              </v-btn>
            </div>

            <div v-else>
              <v-expansion-panels variant="accordion">
                <v-expansion-panel
                  v-for="(exercise, index) in sessionExercises"
                  :key="exercise.id"
                >
                  <v-expansion-panel-title>
                    <div class="d-flex align-center">
                      <span>{{ exercise.name }}</span>
                      <v-chip 
                        size="small" 
                        class="ml-2" 
                        :color="getExerciseStatusColor(exercise)"
                      >
                        {{ getExerciseStatus(exercise) }}
                      </v-chip>
                      <span v-if="exercise.duration" class="ml-2 text-caption">
                        {{ exercise.duration }}
                      </span>
                    </div>

                    <template v-slot:actions>
                      <v-btn
                        icon
                        variant="text"
                        size="small"
                        color="error"
                        @click.stop="removeExercise(exercise)"
                        :disabled="exercise.isActive"
                      >
                        <v-icon>mdi-delete</v-icon>
                      </v-btn>
                    </template>
                  </v-expansion-panel-title>
                  <v-expansion-panel-text>
                    <v-row>
                      <v-col cols="12" md="6">
                        <v-combobox
                          v-model="exercise.bpms"
                          label="BPMs"
                          multiple
                          chips
                          hide-selected
                          clearable
                          variant="outlined"
                          hide-details
                          class="mb-4"
                          :disabled="!sessionStarted || !exercise.isActive"
                          placeholder="Enter BPM values (e.g. 100, 120)"
                          @update:model-value="validateBpms(exercise)"
                        />
                      </v-col>
                      <v-col cols="12" md="6">
                        <v-select
                          v-model="exercise.timeSignature"
                          :items="timeSignatureOptions"
                          label="Time Signature"
                          variant="outlined"
                          hide-details
                          class="mb-4"
                          :disabled="!sessionStarted || !exercise.isActive"
                        ></v-select>
                      </v-col>
                      <v-col cols="12">
                        <v-textarea
                          v-model="exercise.notes"
                          label="Notes for this exercise"
                          variant="outlined"
                          rows="2"
                          placeholder="Add specific notes about this exercise (optional)"
                          hide-details
                          :disabled="!sessionStarted || !exercise.isActive"
                        ></v-textarea>
                      </v-col>
                      <v-col cols="12" class="d-flex justify-end">
                        <v-btn
                          v-if="!exercise.isActive && !exercise.completed"
                          color="primary"
                          prepend-icon="mdi-play"
                          @click="startExercisePractice(exercise)"
                          :disabled="!sessionStarted || (hasActiveExercise && !exercise.isActive)"
                        >
                          Start Practice
                        </v-btn>
                        <v-btn
                          v-else-if="exercise.isActive"
                          color="error"
                          prepend-icon="mdi-stop"
                          @click="stopExercisePractice(exercise)"
                        >
                          Stop Practice
                        </v-btn>
                        <v-btn
                          v-else
                          color="success"
                          prepend-icon="mdi-check"
                          disabled
                        >
                          Completed
                        </v-btn>
                      </v-col>
                    </v-row>
                  </v-expansion-panel-text>
                </v-expansion-panel>
              </v-expansion-panels>

              <div class="d-flex justify-end mt-2">
                <v-btn
                  color="primary"
                  prepend-icon="mdi-plus"
                  variant="text"
                  @click="openAddExerciseDialog"
                  :disabled="!sessionStarted"
                >
                  Add More
                </v-btn>
              </div>
            </div>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- Exercise search and selection -->
    <v-card class="mb-4">
      <v-card-title>
        <div class="text-h6">Choose Exercises</div>
      </v-card-title>

      <v-card-text>
        <!-- Filters -->
        <v-row>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model="exerciseSearch"
              label="Search exercises"
              variant="outlined"
              density="compact"
              prepend-inner-icon="mdi-magnify"
              hide-details
              clearable
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="4">
            <v-select
              v-model="categoryFilter"
              :items="categoriesForSelect"
              label="Filter by category"
              variant="outlined"
              density="compact"
              hide-details
              clearable
            ></v-select>
          </v-col>
          <v-col cols="12" sm="4">
            <v-select
              v-model="tagFilter"
              :items="tagsForSelect"
              label="Filter by tag"
              variant="outlined"
              density="compact"
              hide-details
              clearable
            ></v-select>
          </v-col>
        </v-row>

        <!-- Exercise List -->
        <exercise-list
          :exercises="filteredExercises"
          :loading="exercisesStore.loading"
          display-type="grid"
          :allow-select="true"
          @select-exercise="toggleExerciseInSession"
          :selected-exercise-ids="sessionExerciseIds"
          class="mt-4"
        />
      </v-card-text>
    </v-card>

    <!-- Add Exercise Dialog -->
    <v-dialog v-model="addExerciseDialog" max-width="800">
      <v-card>
        <v-card-title>Add Exercises to Session</v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="12" sm="4">
              <v-text-field
                v-model="dialogExerciseSearch"
                label="Search exercises"
                variant="outlined"
                density="compact"
                prepend-inner-icon="mdi-magnify"
                hide-details
                clearable
                class="mb-4"
              ></v-text-field>
            </v-col>
            <v-col cols="12" sm="4">
              <v-select
                v-model="dialogCategoryFilter"
                :items="categoriesForSelect"
                label="Filter by category"
                variant="outlined"
                density="compact"
                hide-details
                clearable
              ></v-select>
            </v-col>
            <v-col cols="12" sm="4">
              <v-select
                v-model="dialogTagFilter"
                :items="tagsForSelect"
                label="Filter by tag"
                variant="outlined"
                density="compact"
                hide-details
                clearable
              ></v-select>
            </v-col>
          </v-row>

          <exercise-list
            :exercises="dialogFilteredExercises"
            :loading="exercisesStore.loading"
            display-type="grid"
            :allow-select="true"
            @select-exercise="toggleExerciseInSession"
            :selected-exercise-ids="sessionExerciseIds"
            class="mt-4"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="primary" variant="text" @click="addExerciseDialog = false">
            Done
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Cancel Session Dialog -->
    <v-dialog v-model="cancelDialog" max-width="500">
      <v-card>
        <v-card-title class="text-h5">Cancel Session?</v-card-title>
        <v-card-text>
          Are you sure you want to cancel this practice session? All progress will be lost.
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="grey-darken-1" variant="text" @click="cancelDialog = false">No, Continue</v-btn>
          <v-btn 
            color="error" 
            variant="flat" 
            @click="confirmCancelSession"
          >
            Yes, Cancel Session
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Finish Session Form Dialog -->
    <v-dialog v-model="finishSessionDialog" max-width="600" persistent>
      <v-card>
        <v-card-title class="text-h5">Complete Practice Session</v-card-title>
        <v-card-text>
          <v-form ref="finishForm" v-model="finishFormValid">
            <v-textarea
              v-model="notes"
              label="Session Notes"
              variant="outlined"
              rows="4"
              placeholder="Add any final notes about this practice session (optional)"
              class="mb-4"
            ></v-textarea>
            
            <p class="text-body-2 mt-4">
              This will complete your practice session with {{ sessionExercises.length }} exercises 
              and a total duration of {{ sessionDuration }}.
            </p>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="grey-darken-1" variant="text" @click="finishSessionDialog = false">
            Back to Session
          </v-btn>
          <v-btn
            color="success"
            variant="flat"
            @click="completeSession"
            :loading="finishingSession"
          >
            Complete Session
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useExercisesStore } from '@/stores/exercises.js'
import { useCategoriesStore } from '@/stores/categories.js'
import { useTagsStore } from '@/stores/tags.js'
import { useSessionsStore } from '@/stores/sessions.js'
import { useHistoryStore } from '@/stores/history.js'
import { useAppStore } from '@/stores/app.js'
import CategoryChip from '@/components/categories/CategoryChip.vue'
import TagChip from '@/components/tags/TagChip.vue'
import ExerciseList from '@/components/exercises/ExerciseList.vue'

const router = useRouter()
const route = useRoute()
const exercisesStore = useExercisesStore()
const categoriesStore = useCategoriesStore()
const tagsStore = useTagsStore()
const sessionsStore = useSessionsStore()
const historyStore = useHistoryStore()
const appStore = useAppStore()

// Data for session
const sessionStarted = ref(false)
const startTime = ref(null)
const notes = ref('')
const sessionExercises = ref([])
const startingSession = ref(false)
const finishingSession = ref(false)
const loadingSessionExercises = ref(false)
const sessionDuration = ref('0:00:00')
const durationInterval = ref(null)
const sessionId = ref(null)

// Dialog controls
const addExerciseDialog = ref(false)
const dialogExerciseSearch = ref('')
const dialogCategoryFilter = ref(null)
const dialogTagFilter = ref(null)
const exerciseSearch = ref('')
const categoryFilter = ref(null)
const tagFilter = ref(null)
const cancelDialog = ref(false)
const finishSessionDialog = ref(false)
const finishForm = ref(null)
const finishFormValid = ref(true)

// Time signature options
const timeSignatureOptions = [
  '4/4', '3/4', '2/4', '6/8', '7/8', '5/4', '9/8', '12/8'
]

// Exercise tracking
const hasActiveExercise = computed(() => {
  return sessionExercises.value.some(ex => ex.isActive)
})

// Computed properties
const sessionExerciseIds = computed(() => {
  return sessionExercises.value.map(exercise => exercise.id)
})

const filteredExercises = computed(() => {
  let result = exercisesStore.exercises
  
  // Search filter
  if (exerciseSearch.value) {
    const searchLower = exerciseSearch.value.toLowerCase()
    result = result.filter(exercise => 
      exercise.name.toLowerCase().includes(searchLower) || 
      (exercise.description && exercise.description.toLowerCase().includes(searchLower))
    )
  }
  
  // Category filter
  if (categoryFilter.value) {
    result = result.filter(exercise => 
      exercise.categoryIds && exercise.categoryIds.includes(parseInt(categoryFilter.value))
    )
  }
  
  // Tag filter
  if (tagFilter.value) {
    result = result.filter(exercise => 
      exercise.tagIds && exercise.tagIds.includes(parseInt(tagFilter.value))
    )
  }
  
  return result
})

const dialogFilteredExercises = computed(() => {
  let result = exercisesStore.exercises
  
  // Search filter
  if (dialogExerciseSearch.value) {
    const search = dialogExerciseSearch.value.toLowerCase()
    result = result.filter(exercise => 
      exercise.name.toLowerCase().includes(search) || 
      (exercise.description && exercise.description.toLowerCase().includes(search))
    )
  }
  
  // Category filter
  if (dialogCategoryFilter.value) {
    result = result.filter(exercise => 
      exercise.categoryIds && exercise.categoryIds.includes(parseInt(dialogCategoryFilter.value))
    )
  }
  
  // Tag filter
  if (dialogTagFilter.value) {
    result = result.filter(exercise => 
      exercise.tagIds && exercise.tagIds.includes(parseInt(dialogTagFilter.value))
    )
  }
  
  return result
})

const categoriesForSelect = computed(() => {
  return categoriesStore.categories.map(category => ({
    title: category.name,
    value: category.id
  }))
})

const tagsForSelect = computed(() => {
  return tagsStore.tags.map(tag => ({
    title: tag.name,
    value: tag.id
  }))
})

// Methods
function formatDateTime(dateTime) {
  if (!dateTime) return ''
  return appStore.formatDateTime(dateTime)
}

function formatTime(time) {
  if (!time) return ''
  
  const date = new Date(time)
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

function getCategoryById(categoryId) {
  const category = categoriesStore.categoryById(categoryId)
  return category // This may return undefined if category is not found
}

function getTagById(tagId) {
  return tagsStore.tagById(tagId)
}

function isExerciseInSession(exercise) {
  return sessionExercises.value.some(e => e.id === exercise.id)
}

function getExerciseStatus(exercise) {
  if (exercise.isActive) return 'Active'
  if (exercise.completed) return 'Completed'
  return 'Pending'
}

function getExerciseStatusColor(exercise) {
  if (exercise.isActive) return 'success'
  if (exercise.completed) return 'primary'
  return 'grey'
}

function addExerciseToSession(exercise) {
  if (isExerciseInSession(exercise)) {
    removeExercise(exercise)
    return
  }
  
  // Add the exercise with default values and current time
  sessionExercises.value.push({
    ...exercise,
    addedAt: new Date(),
    bpms: [], // Empty initially so user can set it
    timeSignature: '4/4', // Default time signature
    notes: '', // Empty initially for user input
    sessionTags: [], // New field for session-specific tags
    startTime: null,
    endTime: null,
    isActive: false,
    completed: false,
    duration: null
  })
  
  appStore.showSuccessMessage(`Added ${exercise.name} to session`)
}

function removeSessionTag(exercise, index) {
  if (exercise.sessionTags) {
    exercise.sessionTags.splice(index, 1)
  }
}

function toggleExerciseInSession(exercise) {
  if (isExerciseInSession(exercise)) {
    removeExercise(exercise)
  } else {
    addExerciseToSession(exercise)
  }
}

async function removeExercise(exercise) {
  // The exercise passed in may come from one of the sub-components, which doesn't have the full local context
  // of the hisotryID
  const fullExercise = sessionExercises.value.find(e => e.id === exercise.id)

  // Don't allow removing active exercises
  if (exercise.isActive) {
    appStore.showWarningMessage(`Can't remove an active exercise. Stop it first.`)
    return
  }
  
  sessionExercises.value = sessionExercises.value.filter(e => e.id !== exercise.id)

  // Also remove from the backend
  if (fullExercise.completed && fullExercise.historyID) {
      await historyStore.deleteHistoryEntry(fullExercise.historyID)
  }
  
  appStore.showInfoMessage(`Removed ${exercise.name} from session`)
}

function openAddExerciseDialog() {
  dialogExerciseSearch.value = ''
  dialogCategoryFilter.value = null
  dialogTagFilter.value = null
  addExerciseDialog.value = true
}

async function startSession() {
  startingSession.value = true
  
  try {
    // Set start time to now
    startTime.value = new Date()
    
    // Create session in the database immediately
    const sessionData = {
      start_time: startTime.value.toISOString(),
      end_time: startTime.value.toISOString(), // Temporary end time, will be updated later
      notes: notes.value || ''
    }
    
    const newSession = await sessionsStore.createSession(sessionData)
    sessionId.value = newSession.id
    
    // Start the timer to show elapsed time
    startDurationTimer()
    
    sessionStarted.value = true
    appStore.showSuccessMessage('Practice session started!')
  } catch (error) {
    console.error('Error creating session:', error)
    appStore.showErrorMessage(`Failed to start session: ${error.message}`)
  } finally {
    startingSession.value = false
  }
}

function startDurationTimer() {
  durationInterval.value = setInterval(() => {
    if (startTime.value) {
      const now = new Date()
      const diffMs = now - startTime.value
      
      // Format the duration
      const seconds = Math.floor((diffMs / 1000) % 60)
      const minutes = Math.floor((diffMs / (1000 * 60)) % 60)
      const hours = Math.floor(diffMs / (1000 * 60 * 60))
      
      sessionDuration.value = `${hours}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
      
      // Update duration for active exercises
      sessionExercises.value.forEach(exercise => {
        if (exercise.isActive && exercise.startTime) {
          const exDiffMs = now - exercise.startTime
          const exSeconds = Math.floor((exDiffMs / 1000) % 60)
          const exMinutes = Math.floor((exDiffMs / (1000 * 60)) % 60)
          const exHours = Math.floor(exDiffMs / (1000 * 60 * 60))
          
          exercise.duration = `${exHours}:${exMinutes.toString().padStart(2, '0')}:${exSeconds.toString().padStart(2, '0')}`
        }
      })
    }
  }, 1000)
}

function finishSession() {
  // Check if any exercise is still active
  const activeExercise = sessionExercises.value.find(ex => ex.isActive)
  if (activeExercise) {
    appStore.showWarningMessage(`Please stop the active exercise "${activeExercise.name}" before finishing the session.`)
    return
  }
  
  finishSessionDialog.value = true
}

async function completeSession() {
  finishingSession.value = true
  
  try {
    const endTime = new Date()
    
    // Update the session with the end time and notes
    const sessionData = {
      end_time: endTime.toISOString(),
      notes: notes.value || ''
    }
    
    await sessionsStore.updateSession(sessionId.value, sessionData, 'endTime,notes')
    
    appStore.showSuccessMessage('Practice session saved successfully!')
    
    // Navigate to session detail
    router.push({
      name: 'session-detail',
      params: { id: sessionId.value }
    })
  } catch (error) {
    console.error('Failed to complete session:', error)
    appStore.showErrorMessage(`Failed to complete session: ${error.message}`)
    finishingSession.value = false
  }
}

function startExercisePractice(exercise) {
  if (hasActiveExercise.value) {
    // Already checked in the UI, but double-check here
    appStore.showWarningMessage('Please finish the current exercise first.')
    return
  }
  
  // Set exercise as active and record start time
  exercise.isActive = true
  exercise.startTime = new Date()
  exercise.duration = '0:00:00'
  
  appStore.showSuccessMessage(`Started practice: ${exercise.name}`)
}

async function stopExercisePractice(exercise) {
  // Set end time and mark as completed
  exercise.isActive = false
  exercise.endTime = new Date()
  exercise.completed = true
  
  try {
    // Prepare the session tags - convert to array of strings if objects
    const sessionTagsArray = exercise.sessionTags?.map(tag => {
      return typeof tag === 'string' ? tag : (tag.title || tag.text || tag)
    }) || []
    
    // Save the exercise history entry immediately
    const exerciseData = {
      session_id: sessionId.value,
      exercise_id: exercise.id,
      start_time: exercise.startTime.toISOString(),
      end_time: exercise.endTime.toISOString(),
      bpms: exercise.bpms ? exercise.bpms : [],
      time_signature: exercise.timeSignature || '4/4',
      notes: exercise.notes || '',
      additional_tags: sessionTagsArray
    }
    
    const newHistoryEntry = await historyStore.createHistoryEntry(exerciseData)
    exercise.historyID = newHistoryEntry.id
    appStore.showSuccessMessage(`Completed practice: ${exercise.name}`)
  } catch (error) {
    console.error(`Error saving exercise history for ${exercise.name}:`, error)
    appStore.showErrorMessage(`Failed to save exercise history: ${error.message}`)
  }
}

function validateBpms(exercise) {
  // Filter out non-numeric values and convert strings to numbers
  if (exercise.bpms && Array.isArray(exercise.bpms)) {
    exercise.bpms = exercise.bpms
      .map(bpm => {
        // Convert to string first to handle all input types
        const bpmStr = String(bpm).trim()
        // Parse it as a number
        const bpmNum = parseInt(bpmStr, 10)
        // Return the number if valid, otherwise null
        return !isNaN(bpmNum) && bpmNum > 0 ? bpmNum : null
      })
      .filter(bpm => bpm !== null) // Remove nulls
  }
}

function openCancelDialog() {
  if (sessionStarted.value) {
    cancelDialog.value = true
  } else {
    // If session hasn't started, just go back
    router.push({ name: 'sessions' })
  }
}

async function confirmCancelSession() {
  cancelDialog.value = false
  
  // If we have created a session in the database, delete it
  if (sessionId.value) {
    try {
      await sessionsStore.deleteSession(sessionId.value)
    } catch (error) {
      console.error('Error deleting session:', error)
    }
  }
  
  // Clean up
  if (durationInterval.value) {
    clearInterval(durationInterval.value)
  }
  
  // Navigate back
  router.push({ name: 'sessions' })
}

// Check for URL parameter with exercise ID to pre-add
onMounted(async () => {
  // Load exercises if not already loaded
  if (exercisesStore.exercises.length === 0) {
    await exercisesStore.fetchExercises()
  }
  
  if (categoriesStore.categories.length === 0) {
    await categoriesStore.fetchCategories()
  }
  
  if (tagsStore.tags.length === 0) {
    await tagsStore.fetchTags()
  }
  
  // Check for exercise ID in URL
  const exerciseId = Number(route.query.exercise)
  if (exerciseId) {
    const exercise = exercisesStore.exerciseById(exerciseId)
    if (exercise) {
      addExerciseToSession(exercise)
    } else {
      // If not in store, try to fetch it
      try {
        await exercisesStore.fetchExercise(exerciseId)
        const exercise = exercisesStore.exerciseById(exerciseId)
        if (exercise) {
          addExerciseToSession(exercise)
        }
      } catch (error) {
        console.error('Failed to fetch exercise:', error)
      }
    }
  }
})

// Clean up the timer when component is destroyed
onBeforeUnmount(() => {
  if (durationInterval.value) {
    clearInterval(durationInterval.value)
  }
})

// Handle browser back button or leaving the page
window.addEventListener('beforeunload', (event) => {
  if (sessionStarted.value) {
    // Standard way to show dialog before leaving page
    event.preventDefault()
    event.returnValue = ''
  }
})
</script>

<style scoped>
.exercise-card {
  cursor: pointer;
  transition: all 0.2s ease;
}

.exercise-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.selected-exercise {
  border: 2px solid var(--v-primary-base) !important;
  background-color: rgba(var(--v-primary-base), 0.05);
}
</style>
