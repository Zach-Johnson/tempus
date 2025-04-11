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
            <v-btn
              variant="text"
              to="/sessions"
              prepend-icon="mdi-arrow-left"
              class="mb-2"
            >
              Back to Sessions
            </v-btn>
            <h1 class="text-h3">Practice Session</h1>
            <p class="text-body-1">
              {{ formatDate(sessionsStore.currentSession.startTime) }}
            </p>
          </v-col>
          <v-col cols="12" sm="4" class="d-flex justify-end align-center">
            <v-btn
              color="error"
              variant="outlined"
              prepend-icon="mdi-delete"
              class="mr-2"
              @click="confirmDelete"
            >
              Delete
            </v-btn>
            <v-btn
              color="primary"
              variant="flat"
              prepend-icon="mdi-pencil"
              @click="editSession"
            >
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
                <div class="text-body-2 mb-4">{{ formatDuration(sessionsStore.currentSession.startTime, sessionsStore.currentSession.endTime) }}</div>
                
                <div class="text-body-1 font-weight-medium mb-1">Notes</div>
                <div class="text-body-2 mb-4">{{ sessionsStore.currentSession.notes || 'No notes provided' }}</div>
              </v-col>
              
              <v-col cols="12" md="6">
                <div class="text-body-1 font-weight-medium mb-2">Exercises</div>
                <v-list v-if="sessionsStore.currentSession.exercises && sessionsStore.currentSession.exercises.length > 0">
                  <v-list-item
                    v-for="exercise in sessionsStore.currentSession.exercises"
                    :key="exercise.id"
                    :to="`/exercises/${exercise.exerciseId}`"
                  >
                    <v-list-item-title>{{ exercise.exercise ? exercise.exercise.name : `Exercise #${exercise.exerciseId}` }}</v-list-item-title>
                    <v-list-item-subtitle>
                      <span v-if="exercise.bpm">{{ exercise.bpm }} BPM</span>
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
        <h2 class="text-h5 mb-4" v-if="sessionsStore.currentSession.exercises && sessionsStore.currentSession.exercises.length > 0">
          Exercise Details
        </h2>
        <v-row v-if="sessionsStore.currentSession.exercises && sessionsStore.currentSession.exercises.length > 0">
          <v-col 
            v-for="exercise in sessionsStore.currentSession.exercises" 
            :key="exercise.id"
            cols="12" 
            md="6" 
            lg="4"
          >
            <v-card variant="outlined" class="mb-4">
              <v-card-title>
                {{ exercise.exercise ? exercise.exercise.name : `Exercise #${exercise.exerciseId}` }}
              </v-card-title>
              <v-card-text>
                <div class="d-flex mb-2">
                  <div class="text-body-2 font-weight-medium mr-4">
                    <v-icon icon="mdi-metronome" size="small" class="mr-1"></v-icon>
                    {{ exercise.bpm || 'N/A' }} BPM
                  </div>
                  <div class="text-body-2 font-weight-medium">
                    <v-icon icon="mdi-music" size="small" class="mr-1"></v-icon>
                    {{ exercise.timeSignature || '4/4' }}
                  </div>
                </div>
                
                <div class="text-body-2 font-weight-medium mb-1">Duration</div>
                <div class="text-body-2 mb-2">{{ formatDuration(exercise.startTime, exercise.endTime) }}</div>
                
                <div class="text-body-2 font-weight-medium mb-1">Notes</div>
                <div class="text-body-2">{{ exercise.notes || 'No notes provided' }}</div>
                
                <v-divider class="my-3"></v-divider>
                
                <v-btn
                  variant="text"
                  :to="`/exercises/${exercise.exerciseId}`"
                  color="primary"
                  size="small"
                >
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
          <v-btn 
            color="error" 
            variant="flat" 
            @click="deleteSession" 
            :loading="deleteLoading"
          >
            Delete
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
import CategoryChip from '@/components/categories/CategoryChip.vue'

const router = useRouter()
const route = useRoute()
const sessionsStore = useSessionsStore()
const appStore = useAppStore()
const categoriesStore = useCategoriesStore()

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
