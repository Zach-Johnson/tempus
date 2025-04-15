<template>
  <div>
    <v-row>
      <v-col cols="12">
        <h1 class="text-h3 mb-5">Dashboard</h1>
      </v-col>
    </v-row>

    <!-- Summary Cards -->
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
                <v-list-item-title>{{ formatDate(session.startTime) }}</v-list-item-title>
                <v-list-item-subtitle>
                  {{ formatDuration(session.startTime, session.endTime) }}
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
                  <div class="text-h4 text-center primary--text">{{ statsStore.totalSessions || 0 }}</div>
                  <div class="text-body-2 text-center">Total Sessions</div>
                </v-col>
                <v-col cols="6">
                  <div class="text-h4 text-center primary--text">{{ formatTime(statsStore.totalPracticeTime) }}</div>
                  <div class="text-body-2 text-center">Total Practice Time</div>
                </v-col>
              </v-row>
              <v-row class="mt-4">
                <v-col cols="6">
                  <div class="text-h4 text-center primary--text">{{ exercisesStore.exercises.length || 0 }}</div>
                  <div class="text-body-2 text-center">Exercises</div>
                </v-col>
                <v-col cols="6">
                  <div class="text-h4 text-center primary--text">{{ categoriesStore.categories.length || 0 }}</div>
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
            <v-btn variant="text" to="/exercises">View All</v-btn>
          </v-card-title>
          <v-card-text v-if="loading">
            <v-progress-circular indeterminate color="primary"></v-progress-circular>
          </v-card-text>
          <v-list v-else density="compact">
            <template v-if="statsStore.topExercises.length > 0">
              <v-list-item 
                v-for="exercise in statsStore.topExercises" 
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
          <v-card-title class="d-flex align-center text-h6">
            Practice Time (Last 30 Days)
            <v-spacer></v-spacer>
            <v-btn
              variant="text"
              prepend-icon="mdi-refresh"
              @click="loadStats"
              :loading="statsStore.loading"
              size="small"
            >
              Refresh
            </v-btn>
          </v-card-title>
          <v-card-text style="height: 300px">
            <category-practice-time-chart
              :loading="statsStore.loading"
            />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Detailed Statistics -->
    <!-- <v-row class="mt-4"> -->
    <!--   <v-col cols="12" md="6"> -->
    <!--     <v-card> -->
    <!--       <v-card-title class="text-h6"> -->
    <!--         Practice by Category -->
    <!--       </v-card-title> -->
    <!--       <v-card-text> -->
    <!--         <category-summary-stats -->
    <!--           :loading="statsStore.loading" -->
    <!--         /> -->
    <!--       </v-card-text> -->
    <!--     </v-card> -->
    <!--   </v-col> -->
    <!---->
    <!--   <v-col cols="12" md="6"> -->
    <!--     <v-card> -->
    <!--       <v-card-title class="text-h6"> -->
    <!--         Top Exercises -->
    <!--       </v-card-title> -->
    <!--       <v-card-text> -->
    <!--         <top-exercises-stats -->
    <!--           :loading="statsStore.loading" -->
    <!--           :limit="5" -->
    <!--         /> -->
    <!--       </v-card-text> -->
    <!--     </v-card> -->
    <!--   </v-col> -->
    <!-- </v-row> -->

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
          to="/tags"
          prepend-icon="mdi-tag-multiple"
        >
          Tags
        </v-btn>
      </v-col>
    </v-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useAppStore } from '@/stores/app.js';
import { useSessionsStore } from '@/stores/sessions.js';
import { useExercisesStore } from '@/stores/exercises.js';
import { useCategoriesStore } from '@/stores/categories.js';
import { useStatsStore } from '@/stores/stats.js';
import CategoryPracticeTimeChart from '@/components/charts/CategoryPracticeTimeChart.vue';
import CategorySummaryStats from '@/components/stats/CategorySummaryStats.vue';
import TopExercisesStats from '@/components/stats/TopExercisesStats.vue';

const appStore = useAppStore();
const sessionsStore = useSessionsStore();
const exercisesStore = useExercisesStore();
const categoriesStore = useCategoriesStore();
const statsStore = useStatsStore();

// Data states
const loading = ref(true);
const recentSessions = ref([]);

onMounted(async () => {
  try {
    // Load all required data
    await Promise.all([
      loadCategories(),
      loadExercises(),
      fetchRecentSessions(),
      loadStats()
    ]);
  } catch (error) {
    console.error('Error loading dashboard data:', error);
    appStore.showErrorMessage('Error loading dashboard data');
  } finally {
    loading.value = false;
  }
});

// Load categories
async function loadCategories() {
  if (categoriesStore.categories.length === 0) {
    await categoriesStore.fetchCategories();
  }
}

// Load exercises
async function loadExercises() {
  if (exercisesStore.exercises.length === 0) {
    await exercisesStore.fetchExercises();
  }
}

// Fetch recent sessions
async function fetchRecentSessions() {
  try {
    const response = await sessionsStore.fetchSessions({ page_size: 5 });
    recentSessions.value = response.sessions || [];
  } catch (error) {
    console.error('Error fetching recent sessions:', error);
  }
}

// Load stats
async function loadStats() {
  try {
    // Get last 30 days range
    const endDate = new Date();
    const startDate = new Date();
    startDate.setDate(startDate.getDate() - 30);
    
    // Format dates for the API
    const params = {
      start_date: startDate.toISOString(),
      end_date: endDate.toISOString()
    };
    
    console.log('Fetching practice stats with params:', params);
    const result = await statsStore.fetchPracticeStats(params);
    
    // Debug logs
    console.log('Practice stats API response:', result);
    console.log('Practice frequency in store:', statsStore.practiceFrequency);
    console.log('Category distribution in store:', statsStore.categoryDistribution);
    
    // Check for empty data
    if (!statsStore.practiceFrequency || statsStore.practiceFrequency.length === 0) {
      console.warn('No practice frequency data available');
    }
    
    if (!statsStore.categoryDistribution || statsStore.categoryDistribution.length === 0) {
      console.warn('No category distribution data available');
    }
    
    return result;
  } catch (error) {
    console.error('Error loading statistics:', error);
    appStore.showErrorMessage('Error loading statistics');
  }
}

// Format helpers
function formatDate(dateString) {
  return appStore.formatDate(dateString);
}

function formatTime(minutes) {
  return appStore.formatMinutes(minutes);
}

function formatDuration(startTime, endTime) {
  return appStore.formatDuration(startTime, endTime);
}
</script>

<style scoped>
.primary--text {
  color: rgb(var(--v-theme-primary)) !important;
}
</style>
