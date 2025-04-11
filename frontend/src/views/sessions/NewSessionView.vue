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
          :loading="savingSession"
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
                  :disabled="savingSession"
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
              >
                Add Exercises
              </v-btn>
            </div>

            <v-list v-else density="compact">
              <v-list-item
                v-for="exercise in sessionExercises"
                :key="exercise.id"
                :title="exercise.name"
                :subtitle="`Added at ${formatTime(exercise.addedAt)}`"
              >
                <template v-slot:prepend>
                  <v-avatar size="36" color="primary" class="text-white">
                    {{ exercise.name.charAt(0).toUpperCase() }}
                  </v-avatar>
                </template>

                <template v-slot:append>
                  <v-btn
                    icon
                    variant="text"
                    size="small"
                    color="error"
                    @click="removeExercise(exercise)"
                  >
                    <v-icon>mdi-delete</v-icon>
                  </v-btn>
                </template>
              </v-list-item>
            </v-list>

            <div class="d-flex justify-end mt-2" v-if="sessionExercises.length > 0">
              <v-btn
                color="primary"
                prepend-icon="mdi-plus"
                variant="text"
                @click="openAddExerciseDialog"
              >
                Add More
              </v-btn>
            </div>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- Exercise search and selection -->
    <v-card v-if="exercisesStore.exercises.length > 0 && !sessionStarted">
      <v-card-title>
        <div class="text-h6">Choose Exercises</div>
        <v-text-field
          v-model="exerciseSearch"
          append-inner-icon="mdi-magnify"
          label="Search exercises"
          single-line
          hide-details
          variant="underlined"
          class="mx-4"
        ></v-text-field>
      </v-card-title>

      <v-card-text>
        <v-row>
          <v-col cols="12" sm="6" md="4" lg="3" v-for="exercise in filteredExercises" :key="exercise.id">
            <v-card 
              variant="outlined" 
              @click="addExerciseToSession(exercise)"
              class="exercise-card"
              :class="{ 'selected-exercise': isExerciseInSession(exercise) }"
            >
              <v-card-title>{{ exercise.name }}</v-card-title>
              <v-card-text>
                <p class="text-body-2 text-truncate">{{ exercise.description || 'No description' }}</p>
                
                <div v-if="exercise.categoryIds && exercise.categoryIds.length > 0" class="mt-2">
                  <span v-if="exercise.categoryIds.length > 0" class="text-caption text-grey">
                    {{ exercise.categoryIds.length }} {{ exercise.categoryIds.length === 1 ? 'category' : 'categories' }}
                  </span>
                </div>
              </v-card-text>
              
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn
                  variant="text"
                  :color="isExerciseInSession(exercise) ? 'error' : 'primary'"
                  size="small"
                >
                  {{ isExerciseInSession(exercise) ? 'Remove' : 'Add' }}
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- Add Exercise Dialog -->
    <v-dialog v-model="addExerciseDialog" max-width="800">
      <v-card>
        <v-card-title>Add Exercises to Session</v-card-title>
        <v-card-text>
          <v-text-field
            v-model="dialogExerciseSearch"
            label="Search exercises"
            variant="outlined"
            prepend-inner-icon="mdi-magnify"
            clearable
            class="mb-4"
          ></v-text-field>

          <v-row>
            <v-col cols="12" sm="6" md="4" v-for="exercise in dialogFilteredExercises" :key="exercise.id">
              <v-card
                variant="outlined"
                :class="{ 'selected-exercise': isExerciseInSession(exercise) }"
              >
                <v-card-title>{{ exercise.name }}</v-card-title>
                <v-card-text>
                  <p class="text-body-2 text-truncate">{{ exercise.description || 'No description' }}</p>
                </v-card-text>
                <v-card-actions>
                  <v-spacer></v-spacer>
                  <v-btn
                    variant="text"
                    :color="isExerciseInSession(exercise) ? 'error' : 'primary'"
                    @click="toggleExerciseInSession(exercise)"
                  >
                    {{ isExerciseInSession(exercise) ? 'Remove' : 'Add' }}
                  </v-btn>
                </v-card-actions>
              </v-card>
            </v-col>
          </v-row>
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
    <v-dialog v-model="finishSessionDialog" max-width="800" persistent>
      <v-card>
        <v-card-title class="text-h5">Complete Practice Session</v-card-title>
        <v-card-text>
          <v-form ref="finishForm" v-model="finishFormValid">
            <v-textarea
              v-model="notes"
              label="Session Notes"
              variant="outlined"
              rows="3"
              placeholder="Add notes about this practice session (optional)"
              class="mb-4"
            ></v-textarea>

            <div class="text-h6 mb-2">Exercise Details</div>
            <v-expansion-panels variant="accordion" v-if="sessionExercises.length > 0">
              <v-expansion-panel
                v-for="(exercise, index) in sessionExercises"
                :key="exercise.id"
              >
                <v-expansion-panel-title>
                  {{ exercise.name }}
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <v-row>
                    <v-col cols="12" md="6">
                      <v-text-field
                        v-model="exercise.bpm"
                        label="BPM"
                        type="number"
                        min="1"
                        variant="outlined"
                        hide-details
                        class="mb-4"
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-select
                        v-model="exercise.timeSignature"
                        :items="timeSignatureOptions"
                        label="Time Signature"
                        variant="outlined"
                        hide-details
                        class="mb-4"
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
                      ></v-textarea>
                    </v-col>
                  </v-row>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>
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
            @click="saveSession"
            :loading="savingSession"
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
import { useAppStore } from '@/stores/app.js'
import CategoryChip from '@/components/categories/CategoryChip.vue'
import TagChip from '@/components/tags/TagChip.vue'

const router = useRouter()
const route = useRoute()
const exercisesStore = useExercisesStore()
const categoriesStore = useCategoriesStore()
const tagsStore = useTagsStore()
const sessionsStore = useSessionsStore()
const appStore = useAppStore()

// Data for session
const sessionStarted = ref(false)
const startTime = ref(null)
const notes = ref('')
const sessionExercises = ref([])
const startingSession = ref(false)
const savingSession = ref(false)
const loadingSessionExercises = ref(false)
const sessionDuration = ref('0:00:00')
const durationInterval = ref(null)

// Dialog controls
const addExerciseDialog = ref(false)
const dialogExerciseSearch = ref('')
const exerciseSearch = ref('')
const cancelDialog = ref(false)
const finishSessionDialog = ref(false)
const finishForm = ref(null)
const finishFormValid = ref(true)

// Time signature options
const timeSignatureOptions = [
  '4/4', '3/4', '2/4', '6/8', '7/8', '5/4', '9/8', '12/8'
]

// Computed properties
const filteredExercises = computed(() => {
  if (!exerciseSearch.value) {
    return exercisesStore.exercises
  }
  
  const search = exerciseSearch.value.toLowerCase()
  return exercisesStore.exercises.filter(exercise => 
    exercise.name.toLowerCase().includes(search) || 
    (exercise.description && exercise.description.toLowerCase().includes(search))
  )
})

const dialogFilteredExercises = computed(() => {
  if (!dialogExerciseSearch.value) {
    return exercisesStore.exercises
  }
  
  const search = dialogExerciseSearch.value.toLowerCase()
  return exercisesStore.exercises.filter(exercise => 
    exercise.name.toLowerCase().includes(search) || 
    (exercise.description && exercise.description.toLowerCase().includes(search))
  )
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

function addExerciseToSession(exercise) {
  if (isExerciseInSession(exercise)) {
    removeExercise(exercise)
    return
  }
  
  // Add the exercise with default values and current time
  sessionExercises.value.push({
    ...exercise,
    addedAt: new Date(),
    bpm: '',
    timeSignature: '4/4',
    notes: '',
    startTime: new Date(),
    endTime: null
  })
  
  appStore.showSuccessMessage(`Added ${exercise.name} to session`)
}

function toggleExerciseInSession(exercise) {
  if (isExerciseInSession(exercise)) {
    removeExercise(exercise)
  } else {
    addExerciseToSession(exercise)
  }
}

function removeExercise(exercise) {
  sessionExercises.value = sessionExercises.value.filter(e => e.id !== exercise.id)
  appStore.showInfoMessage(`Removed ${exercise.name} from session`)
}

function openAddExerciseDialog() {
  dialogExerciseSearch.value = ''
  addExerciseDialog.value = true
}

function startSession() {
  startingSession.value = true
  
  setTimeout(() => {
    // Set start time to now
    startTime.value = new Date()
    sessionStarted.value = true
    startingSession.value = false
    
    // Start the timer to show elapsed time
    startDurationTimer()
    
    appStore.showSuccessMessage('Practice session started!')
  }, 500) // Simulated brief delay
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
    }
  }, 1000)
}

function finishSession() {
  finishSessionDialog.value = true
}

async function saveSession() {
  if (sessionExercises.value.length === 0) {
    appStore.showWarningMessage('Please add at least one exercise to this session')
    return
  }
  
  savingSession.value = true
  
  try {
    const endTime = new Date()
    
    // Prepare session data
    const sessionData = {
      start_time: startTime.value.toISOString(),
      end_time: endTime.toISOString(),
      notes: notes.value || '',
      exercises: sessionExercises.value.map(exercise => ({
        exercise_id: exercise.id,
        start_time: exercise.startTime.toISOString(),
        end_time: endTime.toISOString(),
        bpm: exercise.bpm ? parseInt(exercise.bpm) : null,
        time_signature: exercise.timeSignature || '4/4',
        notes: exercise.notes || ''
      }))
    }
    
    // Send to server
    const newSession = await sessionsStore.createSession(sessionData)
    
    appStore.showSuccessMessage('Practice session saved successfully!')
    
    // Navigate to session detail
    router.push({
      name: 'session-detail',
      params: { id: newSession.id }
    })
  } catch (error) {
    console.error('Failed to save session:', error)
    appStore.showErrorMessage(`Failed to save session: ${error.message}`)
    savingSession.value = false
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

function confirmCancelSession() {
  cancelDialog.value = false
  
  // Clean up
  if (durationInterval.value) {
    clearInterval(durationInterval.value)
  }
  
  // Navigate back
  router.push({ name: 'home' })
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
