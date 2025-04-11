<template>
  <div>
    <v-row>
      <v-col cols="12">
        <h1 class="text-h3 mb-5">Dashboard</h1>
      </v-col>
    </v-row>

    <!-- Recent Activity Summary -->
    <v-row>
      <v-col cols="12" md="4">
        <v-card class="mb-4">
          <v-card-title class="text-h6">
            Recent Sessions
            <v-spacer></v-spacer>
            <v-btn variant="text" to="/sessions">View All</v-btn>
          </v-card-title>
          <v-card-text v-if="loading">
            <v-progress-circular indeterminate color="primary"></v-progress-circular>
          </v-card-text>
          <v-list v-else>
            <template v-if="recentSessions.length > 0">
              <v-list-item 
                v-for="session in recentSessions" 
                :key="session.id" 
                :to="`/sessions/${session.id}`">
                <v-list-item-title>{{ formatDate(session.start_time) }}</v-list-item-title>
                <v-list-item-subtitle>
                  {{ formatDuration(session.start_time, session.end_time) }}
                </v-list-item-subtitle>
              </v-list-item>
            </template>
            <v-list-item v-else>
              <v-list-item-title>No recent sessions</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card>
      </v-col>

      <v-col cols="12" md="4">
        <v-card class="mb-4">
          <v-card-title class="text-h6">
            Quick Stats
          </v-card-title>
          <v-card-text>
            <div v-if="loading" class="d-flex justify-center">
              <v-progress-circular indeterminate color="primary"></v-progress-circular>
            </div>
            <div v-else>
              <v-row>
                <v-col cols="6">
                  <div class="text-h4 text-center primary--text">{{ stats.totalSessions || 0 }}</div>
                  <div class="text-body-2 text-center">Total Sessions</div>
                </v-col>
                <v-col cols="6">
                  <div class="text-h4 text-center primary--text">{{ formatTime(stats.totalDuration) }}</div>
                  <div class="text-body-2 text-center">Total Practice Time</div>
                </v-col>
              </v-row>
              <v-row class="mt-4">
                <v-col cols="6">
                  <div class="text-h4 text-center primary--text">{{ stats.exerciseCount || 0 }}</div>
                  <div class="text-body-2 text-center">Exercises</div>
                </v-col>
                <v-col cols="6">
                  <div class="text-h4 text-center primary--text">{{ stats.categoryCount || 0 }}</div>
                  <div class="text-body-2 text-center">Categories</div>
                </v-col>
              </v-row>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" md="4">
        <v-card class="mb-4">
          <v-card-title class="text-h6">
            Most Practiced
            <v-spacer></v-spacer>
            <v-btn variant="text" to="/stats">View Stats</v-btn>
          </v-card-title>
          <v-card-text v-if="loading">
            <v-progress-circular indeterminate color="primary"></v-progress-circular>
          </v-card-text>
          <v-list v-else density="compact">
            <template v-if="topExercises.length > 0">
              <v-list-item 
                v-for="exercise in topExercises" 
                :key="exercise.id"
                :to="`/exercises/${exercise.id}`">
                <v-list-item-title>{{ exercise.name }}</v-list-item-title>
                <v-list-item-subtitle>
                  {{ formatTime(exercise.duration) }}
                </v-list-item-subtitle>
              </v-list-item>
            </template>
            <v-list-item v-else>
              <v-list-item-title>No data available</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card>
      </v-col>
    </v-row>

    <!-- Practice Time Chart -->
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title class="text-h6">
            Practice Time (Last 30 Days)
          </v-card-title>
          <v-card-text style="height: 300px">
            <div v-if="loading" class="d-flex justify-center align-center" style="height: 100%">
              <v-progress-circular indeterminate color="primary"></v-progress-circular>
            </div>
            <div v-else-if="practiceData.length > 0" style="height: 100%">
              <practice-time-chart :chart-data="practiceData" />
            </div>
            <div v-else class="d-flex justify-center align-center" style="height: 100%">
              <p class="text-body-1">No practice data available for the last 30 days</p>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Quick Actions -->
    <v-row class="mt-5">
      <v-col cols="12">
        <h2 class="text-h5 mb-4">Quick Actions</h2>
      </v-col>
      <v-col cols="6" md="3">
        <v-btn
          block
          color="primary"
          size="large"
          to="/sessions/new"
          prepend-icon="mdi-calendar-plus"
        >
          New Session
        </v-btn>
      </v-col>
      <v-col cols="6" md="3">
        <v-btn
          block
          size="large"
          to="/exercises"
          prepend-icon="mdi-music-note"
        >
          Exercises
        </v-btn>
      </v-col>
      <v-col cols="6" md="3">
        <v-btn
          block
          size="large"
          to="/categories"
          prepend-icon="mdi-folder"
        >
          Categories
        </v-btn>
      </v-col>
      <v-col cols="6" md="3">
        <v-btn
          block
          size="large"
          to="/stats"
          prepend-icon="mdi-chart-bar"
        >
          Statistics
        </v-btn>
      </v-col>
    </v-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { format, parseISO, differenceInMinutes } from 'date-fns'
import PracticeTimeChart from '@/components/charts/PracticeTimeChart.vue'

// Mock data states
const loading = ref(true)
const recentSessions = ref([])
const stats = ref({})
const topExercises = ref([])
const practiceData = ref([])

onMounted(async () => {
  try {
    // In a real app, these would be API calls
    await Promise.all([
      fetchRecentSessions(),
      fetchStats(),
      fetchTopExercises(),
      fetchPracticeData()
    ])
  } catch (error) {
    console.error('Error loading dashboard data:', error)
  } finally {
    loading.value = false
  }
})

// Mock data fetching functions
async function fetchRecentSessions() {
  // Simulate API call
  await new Promise(resolve => setTimeout(resolve, 500))
  
  // Mock data
  recentSessions.value = [
    { 
      id: 1, 
      start_time: '2023-11-01T18:30:00Z', 
      end_time: '2023-11-01T19:45:00Z'
    },
    { 
      id: 2, 
      start_time: '2023-10-29T17:00:00Z', 
      end_time: '2023-10-29T18:30:00Z'
    },
    { 
      id: 3, 
      start_time: '2023-10-27T19:15:00Z', 
      end_time: '2023-10-27T20:00:00Z'
    }
  ]
}

async function fetchStats() {
  // Simulate API call
  await new Promise(resolve => setTimeout(resolve, 700))
  
  // Mock data
  stats.value = {
    totalSessions: 42,
    totalDuration: 3760, // in minutes
    exerciseCount: 28,
    categoryCount: 5
  }
}

async function fetchTopExercises() {
  // Simulate API call
  await new Promise(resolve => setTimeout(resolve, 600))
  
  // Mock data
  topExercises.value = [
    { id: 5, name: 'Paradiddles', duration: 420 },
    { id: 8, name: 'Double Stroke Roll', duration: 380 },
    { id: 12, name: 'Flam Accent', duration: 310 },
    { id: 3, name: 'Six Stroke Roll', duration: 290 }
  ]
}

async function fetchPracticeData() {
  // Simulate API call
  await new Promise(resolve => setTimeout(resolve, 800))
  
  // Mock data for the last 30 days
  const data = []
  const now = new Date()
  
  for (let i = 29; i >= 0; i--) {
    const date = new Date()
    date.setDate(now.getDate() - i)
    
    // Random practice duration between 0-120 minutes (some days might have no practice)
    const randomDuration = Math.floor(Math.random() * 120) * (Math.random() > 0.3 ? 1 : 0)
    
    data.push({
      date: format(date, 'yyyy-MM-dd'),
      minutes: randomDuration
    })
  }
  
  practiceData.value = data
}

// Utility functions
function formatDate(dateString) {
  if (!dateString) return ''
  return format(parseISO(dateString), 'MMM d, yyyy')
}

function formatDuration(startTime, endTime) {
  if (!startTime || !endTime) return ''
  const minutes = differenceInMinutes(parseISO(endTime), parseISO(startTime))
  return formatTime(minutes)
}

function formatTime(minutes) {
  if (!minutes) return '0min'
  if (minutes < 60) return `${minutes}min`
  
  const hours = Math.floor(minutes / 60)
  const remainingMinutes = minutes % 60
  
  if (remainingMinutes === 0) return `${hours}h`
  return `${hours}h ${remainingMinutes}min`
}
</script>
