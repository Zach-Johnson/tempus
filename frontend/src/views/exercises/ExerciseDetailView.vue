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

    <!-- Exercise details -->
    <template v-else>
      <div v-if="!exercisesStore.currentExercise">
        <v-alert type="info" class="my-4">
          Exercise not found. It may have been deleted.
          <div class="mt-4">
            <v-btn color="primary" to="/exercises">
              Back to Exercises
            </v-btn>
          </div>
        </v-alert>
      </div>
      <template v-else>
        <!-- Header section with actions -->
        <v-row class="mb-4">
          <v-col cols="12" sm="8">
            <v-btn variant="text" to="/exercises" prepend-icon="mdi-arrow-left" class="mb-2">
              Back to Exercises
            </v-btn>
            <h1 class="text-h3">{{ exercisesStore.currentExercise.name }}</h1>
          </v-col>
          <v-col cols="12" sm="4" class="d-flex justify-end align-center">
            <v-btn color="primary" variant="outlined" class="mr-2" prepend-icon="mdi-pencil" @click="openEditDialog">
              Edit
            </v-btn>
            <v-btn color="primary" variant="flat" prepend-icon="mdi-play" @click="startPractice">
              Practice Now
            </v-btn>
          </v-col>
        </v-row>

        <!-- Exercise details card -->
        <v-card class="mb-6">
          <v-card-text>
            <v-row>
              <v-col cols="12" md="6">
                <div class="text-body-1 font-weight-medium mb-1">Description</div>
                <div class="text-body-2 mb-4">{{ exercisesStore.currentExercise.description || 'No description provided'
                }}</div>

                <div class="text-body-1 font-weight-medium mb-1">Categories</div>
                <div v-if="exerciseCategories.length === 0" class="text-body-2 mb-4 text-grey">
                  No categories assigned
                </div>
                <div v-else class="mb-4">
                  <category-list :categories="exerciseCategories" />
                </div>

                <div class="text-body-1 font-weight-medium mb-1">Tags</div>
                <div v-if="exerciseTags.length === 0" class="text-body-2 mb-4 text-grey">
                  No tags assigned
                </div>
                <div v-else class="mb-4">
                  <tag-list :tags="exerciseTags" />
                </div>

                <div class="text-body-1 font-weight-medium mb-1">Created</div>
                <div class="text-body-2 mb-4">{{ formatDate(exercisesStore.currentExercise.createdAt) }}</div>

                <div class="text-body-1 font-weight-medium mb-1">Last Updated</div>
                <div class="text-body-2">{{ formatDate(exercisesStore.currentExercise.updatedAt) }}</div>
              </v-col>

              <v-col cols="12" md="6">
                <div class="text-body-1 font-weight-medium mb-1">External Resources</div>

                <div v-if="!exercisesStore.currentExercise.links || exercisesStore.currentExercise.links.length === 0"
                  class="text-body-2 mb-4 text-grey">
                  No external resources available
                </div>

                <v-list v-else density="compact" class="mb-4">
                  <v-list-item v-for="(link, index) in exercisesStore.currentExercise.links" :key="index">
                    <template v-slot:prepend>
                      <v-icon icon="mdi-link" size="small"></v-icon>
                    </template>

                    <v-list-item-title>
                      <a :href="link.url" target="_blank" rel="noopener noreferrer">
                        {{ link.description || link.url }}
                      </a>
                    </v-list-item-title>
                  </v-list-item>
                </v-list>

                <!-- Add resource button -->
                <div class="d-flex justify-end mt-2">
                  <v-btn size="small" prepend-icon="mdi-plus" @click="openAddResourceDialog">
                    Add Resource
                  </v-btn>
                </div>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>

        <!-- Exercise images card -->
        <v-container>
          <v-row v-if="exercisesStore.currentExercise.images.length > 0" class="justify-center">
            <v-col cols="12" sm="8" md="6" v-for="(image, index) in exercisesStore.currentExercise.images" :key="index"
              class="d-flex justify-center">
              <v-card class="mb-4" width="100%" max-width="1000">
                <v-img :src="'data:image/jpeg;base64,' + image.imageData" class="bg-grey-lighten-4" contain></v-img>
              </v-card>
            </v-col>
          </v-row>
        </v-container>

        <!-- Practice Statistics Card -->
        <v-row>
          <v-col cols="12">
            <v-card>
              <v-card-title class="d-flex align-center">
                Practice Statistics
                <v-spacer></v-spacer>
                <v-btn variant="text" prepend-icon="mdi-refresh" @click="loadStats" :loading="statsLoading">
                  Refresh
                </v-btn>
              </v-card-title>

              <v-card-text v-if="statsLoading">
                <div class="d-flex justify-center py-4">
                  <v-progress-circular indeterminate color="primary"></v-progress-circular>
                </div>
              </v-card-text>

              <template v-else-if="exerciseStats">
                <v-card-text>
                  <v-row>
                    <v-col cols="12" sm="6" md="3">
                      <div class="text-h4 text-center primary--text">{{ exerciseStats.practiceCount }}</div>
                      <div class="text-body-2 text-center">Times Practiced</div>
                    </v-col>
                    <v-col cols="12" sm="6" md="3">
                      <div class="text-h4 text-center primary--text">{{
                        formatTime(exerciseStats.totalPracticeDurationSeconds / 60) }}</div>
                      <div class="text-body-2 text-center">Total Practice Time</div>
                    </v-col>
                    <v-col cols="12" sm="6" md="3">
                      <div class="text-h4 text-center primary--text">{{ exerciseStats.maxBpm || 'N/A' }}</div>
                      <div class="text-body-2 text-center">Max BPM</div>
                    </v-col>
                    <v-col cols="12" sm="6" md="3">
                      <div class="text-h4 text-center primary--text">{{ exerciseStats.avgRating ?
                        exerciseStats.avgRating.toFixed(1) : 'N/A' }}</div>
                      <div class="text-body-2 text-center">Average Rating</div>
                    </v-col>
                  </v-row>

                  <!-- BPM Progress Chart will be added here when we have data -->
                </v-card-text>
              </template>

              <v-card-text v-else class="text-center py-6">
                <v-icon icon="mdi-chart-timeline-variant" size="64" color="grey-lighten-1" class="mb-4"></v-icon>
                <h3 class="text-h6 mb-2">No practice data available</h3>
                <p class="text-body-2 text-grey mb-4">
                  Start practicing this exercise to collect statistics
                </p>
                <v-btn color="primary" prepend-icon="mdi-play" @click="startPractice">
                  Practice Now
                </v-btn>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>

        <!-- Recent Practice Sessions -->
        <v-row class="mt-4">
          <v-col cols="12">
            <v-card>
              <v-card-title>
                Recent Practice Sessions
                <v-spacer></v-spacer>
                <v-btn variant="text" to="/sessions" :disabled="!hasRecentSessions">
                  View All
                </v-btn>
              </v-card-title>

              <v-list v-if="recentSessions.length > 0">
                <v-list-item v-for="session in recentSessions" :key="session.id" :to="`/sessions/${session.id}`">
                  <v-list-item-title>{{ formatDate(session.startTime) }}</v-list-item-title>
                  <v-list-item-subtitle>
                    {{ formatDuration(session.startTime, session.endTime) }}
                    <span v-if="session.bpm" class="ml-2">{{ session.bpm }} BPM</span>
                  </v-list-item-subtitle>
                </v-list-item>
              </v-list>

              <v-card-text v-else class="text-center py-6">
                <p class="text-body-2 text-grey">No recent practice sessions</p>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </template>
    </template>

    <!-- Edit Dialog -->
    <exercise-form-dialog v-if="exercisesStore.currentExercise" v-model="dialogVisible"
      :exercise="exercisesStore.currentExercise" :is-edit="true" @save="saveExercise" />

    <!-- Add Resource Dialog -->
    <v-dialog v-model="resourceDialog" max-width="600">
      <v-card>
        <v-card-title>Add External Resource</v-card-title>
        <v-card-text>
          <v-form ref="resourceForm" @submit.prevent="saveResource" v-model="resourceFormValid">
            <v-text-field v-model="resourceData.url" label="URL" :rules="urlRules" required variant="outlined"
              placeholder="https://example.com" class="mb-4"></v-text-field>

            <v-text-field v-model="resourceData.description" label="Description" variant="outlined"
              placeholder="Description of the resource (optional)"></v-text-field>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="grey-darken-1" variant="text" @click="resourceDialog = false">Cancel</v-btn>
          <v-btn color="primary" variant="flat" @click="saveResource" :disabled="!resourceFormValid"
            :loading="addingResource">
            Add
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useExercisesStore } from '@/stores/exercises.js'
import { useCategoriesStore } from '@/stores/categories.js'
import { useTagsStore } from '@/stores/tags.js'
import { useAppStore } from '@/stores/app.js'
import { useSessionsStore } from '@/stores/sessions.js'
import ExerciseFormDialog from '@/components/exercises/ExerciseFormDialog.vue'
import CategoryList from '@/components/categories/CategoryList.vue'
import TagList from '@/components/tags/TagList.vue'

const router = useRouter()
const route = useRoute()
const exercisesStore = useExercisesStore()
const categoriesStore = useCategoriesStore()
const tagsStore = useTagsStore()
const appStore = useAppStore()
const sessionsStore = useSessionsStore()

// Data
const loading = ref(true)
const error = ref(null)
const dialogVisible = ref(false)
const statsLoading = ref(false)
const exerciseStats = ref(null)
const recentSessions = ref([])

// Add Resource Dialog
const resourceDialog = ref(false)
const resourceForm = ref(null)
const resourceFormValid = ref(false)
const resourceData = ref({
  url: '',
  description: ''
})
const addingResource = ref(false)

// Validation rules
const urlRules = [
  v => !!v || 'URL is required',
  v => /^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w.-]*)*\/?$/.test(v) || 'Enter a valid URL'
]

// Computed
const exerciseId = computed(() => parseInt(route.params.id))

const exerciseCategories = computed(() => {
  if (!exercisesStore.currentExercise || !exercisesStore.currentExercise.categoryIds) {
    return []
  }

  return exercisesStore.currentExercise.categoryIds
    .map(id => categoriesStore.categoryById(id))
    .filter(Boolean) // Filter out undefined values
})

const exerciseTags = computed(() => {
  if (!exercisesStore.currentExercise || !exercisesStore.currentExercise.tagIds) {
    return []
  }

  return exercisesStore.currentExercise.tagIds
    .map(id => tagsStore.tagById(id))
    .filter(Boolean) // Filter out undefined values
})

const hasRecentSessions = computed(() => recentSessions.value.length > 0)

// Methods
function formatDate(dateString) {
  return appStore.formatDate(dateString)
}

function formatTime(minutes) {
  return appStore.formatMinutes(minutes)
}

function formatDuration(startTime, endTime) {
  return appStore.formatDuration(startTime, endTime)
}

function openEditDialog() {
  dialogVisible.value = true
}

function openAddResourceDialog() {
  resourceData.value = {
    url: '',
    description: ''
  }
  resourceDialog.value = true
}

async function saveExercise(exerciseData) {
  try {
    await exercisesStore.updateExercise(exerciseId.value, exerciseData, 'name,description,categoryIds,tagIds')
    appStore.showSuccessMessage(`Exercise "${exerciseData.name}" updated successfully`)
    dialogVisible.value = false
  } catch (error) {
    appStore.showErrorMessage(`Error updating exercise: ${error.message}`)
  }
}

async function saveResource() {
  if (!resourceFormValid.value) return

  addingResource.value = true
  try {
    await exercisesStore.addExerciseLink(exerciseId.value, {
      url: resourceData.value.url,
      description: resourceData.value.description
    })

    appStore.showSuccessMessage('Resource added successfully')
    resourceDialog.value = false

    // Refresh the exercise data
    await exercisesStore.fetchExercise(exerciseId.value)
  } catch (error) {
    appStore.showErrorMessage(`Error adding resource: ${error.message}`)
  } finally {
    addingResource.value = false
  }
}

async function startPractice() {
  // Navigate to new session with this exercise pre-selected
  router.push({
    name: 'new-session',
    query: { exercise: exerciseId.value }
  })
}

async function loadStats() {
  statsLoading.value = true
  try {
    const stats = await exercisesStore.fetchExerciseStats(exerciseId.value)
    exerciseStats.value = stats
  } catch (error) {
    console.error('Error loading exercise stats:', error)
    appStore.showErrorMessage('Failed to load exercise statistics')
  } finally {
    statsLoading.value = false
  }
}

async function loadRecentSessions() {
  try {
    // Fetch sessions that include this exercise
    await sessionsStore.fetchSessions({ exercise_id: exerciseId.value, page_size: 5 })
    recentSessions.value = sessionsStore.sessions
      .filter(session =>
        session.exercises &&
        session.exercises.some(ex => ex.exerciseId === exerciseId.value)
      )
      .slice(0, 5) // Limit to 5 items
  } catch (error) {
    console.error('Error loading recent sessions:', error)
  }
}

// Load data
async function loadData() {
  loading.value = true
  error.value = null

  try {
    // Load exercise details
    await exercisesStore.fetchExercise(exerciseId.value)

    // Load related data
    if (categoriesStore.categories.length === 0) {
      await categoriesStore.fetchCategories()
    }

    if (tagsStore.tags.length === 0) {
      await tagsStore.fetchTags()
    }

    // Load exercise stats
    await loadStats()

    // Load recent sessions
    await loadRecentSessions()
  } catch (err) {
    console.error(`Error loading exercise: ${err.message}`, err)
    error.value = `Error loading exercise: ${err.message}`
  } finally {
    loading.value = false
  }
}

// Watch for route param changes to reload data
watch(() => route.params.id, () => {
  loadData()
})

// Lifecycle
onMounted(async () => {
  await loadData()
})
</script>

<style scoped>
.primary--text {
  color: rgb(var(--v-theme-primary)) !important;
}
</style>
