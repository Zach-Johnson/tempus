<template>
  <div>
    <v-row class="mb-4">
      <v-col cols="12" sm="8">
        <h1 class="text-h3">Practice Sessions</h1>
        <p class="text-body-1">View your practice history</p>
      </v-col>
      <v-col cols="12" sm="4" class="d-flex justify-end align-center">
        <v-btn
          color="primary"
          prepend-icon="mdi-plus"
          :to="{ name: 'new-session' }"
        >
          New Session
        </v-btn>
      </v-col>
    </v-row>

    <!-- Search and filter -->
    <v-card class="mb-6">
      <v-card-text>
        <v-row>
          <v-col cols="12" sm="6" md="4">
            <v-text-field
              v-model="search"
              label="Search sessions"
              variant="outlined"
              density="compact"
              prepend-inner-icon="mdi-magnify"
              hide-details
              clearable
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-select
              v-model="exerciseFilter"
              :items="exercisesForSelect"
              label="Filter by exercise"
              variant="outlined"
              density="compact"
              hide-details
              clearable
            ></v-select>
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-menu
              v-model="dateMenu"
              :close-on-content-click="false"
              location="bottom"
            >
              <template v-slot:activator="{ props }">
                <v-text-field
                  v-model="dateRangeText"
                  label="Date range"
                  variant="outlined"
                  density="compact"
                  prepend-inner-icon="mdi-calendar"
                  readonly
                  v-bind="props"
                  clearable
                  @click:clear="clearDateRange"
                ></v-text-field>
              </template>
              <v-date-picker
                v-model="dateRange"
                range
                @update:model-value="dateMenu = false"
              ></v-date-picker>
            </v-menu>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- Sessions list -->
    <v-card>
      <v-card-text v-if="sessionsStore.loading">
        <div class="d-flex justify-center align-center py-4">
          <v-progress-circular indeterminate color="primary"></v-progress-circular>
        </div>
      </v-card-text>
      
      <template v-else>
        <v-data-table
          v-if="filteredSessions.length > 0"
          :headers="headers"
          :items="filteredSessions"
          :items-per-page="10"
          class="elevation-1"
        >
          <template v-slot:item.date="{ item }">
            {{ formatDate(item.startTime) }}
          </template>
          
          <template v-slot:item.duration="{ item }">
            {{ formatDuration(item.startTime, item.endTime) }}
          </template>
          
          <template v-slot:item.exercises="{ item }">
            <v-chip
              v-if="item.exercises"
              color="primary"
              size="small"
              variant="outlined"
            >
              {{ item.exercises.length }} exercises
            </v-chip>
            <span v-else class="text-grey">No exercises</span>
          </template>
          
          <template v-slot:item.actions="{ item }">
            <v-btn
              icon
              variant="text"
              size="small"
              color="primary"
              :to="{ name: 'session-detail', params: { id: item.id }}"
            >
              <v-icon>mdi-eye</v-icon>
            </v-btn>
          </template>
        </v-data-table>
        
        <v-card-text v-else class="text-center py-8">
          <v-icon icon="mdi-calendar-blank" size="64" color="grey-lighten-1" class="mb-4"></v-icon>
          <h3 class="text-h6 mb-2">No sessions found</h3>
          <p class="text-body-2 text-grey">
            {{ 
              sessionsStore.sessions.length === 0 
                ? "You haven't recorded any practice sessions yet" 
                : "No sessions match your search criteria" 
            }}
          </p>
          <v-btn
            v-if="sessionsStore.sessions.length === 0"
            color="primary"
            class="mt-4"
            :to="{ name: 'new-session' }"
          >
            Start Practicing
          </v-btn>
        </v-card-text>
      </template>
    </v-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useSessionsStore } from '@/stores/sessions.js'
import { useExercisesStore } from '@/stores/exercises.js'
import { useAppStore } from '@/stores/app.js'
import { format, parseISO } from 'date-fns'

const router = useRouter()
const route = useRoute()
const sessionsStore = useSessionsStore()
const exercisesStore = useExercisesStore()
const appStore = useAppStore()

// Data
const search = ref('')
const exerciseFilter = ref(null)
const dateRange = ref([])
const dateMenu = ref(false)

// Table headers
const headers = [
  { title: 'Date', key: 'date' },
  { title: 'Duration', key: 'duration' },
  { title: 'Exercises', key: 'exercises' },
  { title: 'Notes', key: 'notes' },
  { title: 'Active', key: 'active', value: item => { return item.active ? 'yes' : 'no' }},
  { title: 'Actions', key: 'actions', sortable: false, align: 'end' }
]

// Computed
const dateRangeText = computed(() => {
  if (!dateRange.value || dateRange.value.length !== 2) return ''
  
  const [start, end] = dateRange.value
  if (!start || !end) return ''
  
  return `${start} to ${end}`
})

const exercisesForSelect = computed(() => {
  return exercisesStore.exercises.map(exercise => ({
    title: exercise.name,
    value: exercise.id
  }))
})

const filteredSessions = computed(() => {
  let result = sessionsStore.sessions
  
  // Search filter
  if (search.value) {
    const searchLower = search.value.toLowerCase()
    result = result.filter(session => 
      (session.notes && session.notes.toLowerCase().includes(searchLower))
    )
  }
  
  // Exercise filter
  if (exerciseFilter.value) {
    result = result.filter(session => 
      session.exercises && session.exercises.some(ex => 
        ex.exerciseId === parseInt(exerciseFilter.value) ||
        (ex.exercise && ex.exercise.id === parseInt(exerciseFilter.value))
      )
    )
  }
  
  // Date range filter
  if (dateRange.value && dateRange.value.length === 2 && dateRange.value[0] && dateRange.value[1]) {
    const startDate = new Date(dateRange.value[0])
    // Set to end of day for the end date
    const endDate = new Date(dateRange.value[1])
    endDate.setHours(23, 59, 59, 999)
    
    result = result.filter(session => {
      const sessionDate = new Date(session.startTime)
      return sessionDate >= startDate && sessionDate <= endDate
    })
  }
  
  // Sort by date, newest first
  return [...result].sort((a, b) => new Date(b.startTime) - new Date(a.startTime))
})

// Methods
function formatDate(dateString) {
  return appStore.formatDate(dateString)
}

function formatDuration(startTime, endTime) {
  return appStore.formatDuration(startTime, endTime)
}

function clearDateRange() {
  dateRange.value = []
}

// Lifecycle
onMounted(async () => {
  // Load sessions if not already loaded
  if (sessionsStore.sessions.length === 0) {
    await sessionsStore.fetchSessions()
  }
  
  // Load exercises if not already loaded
  if (exercisesStore.exercises.length === 0) {
    await exercisesStore.fetchExercises()
  }
})
</script>
