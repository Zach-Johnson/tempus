<template>
  <div class="dashboard-container">
    <h1 class="mb-4">Dashboard</h1>
    
    <div class="row">
      <!-- Practice summary -->
      <div class="col-md-6 col-lg-3 mb-4">
        <div class="card h-100">
          <div class="card-body">
            <h5 class="card-title">Practice Summary</h5>
            <div class="summary-stats">
              <div class="stat-item">
                <div class="stat-value">{{ recentSessionsCount }}</div>
                <div class="stat-label">Recent Sessions</div>
              </div>
              <div class="stat-item">
                <div class="stat-value">{{ totalMinutesPracticed }}</div>
                <div class="stat-label">Minutes Practiced</div>
              </div>
              <div class="stat-item">
                <div class="stat-value">{{ exercisesCount }}</div>
                <div class="stat-label">Exercises</div>
              </div>
            </div>
          </div>
          <div class="card-footer">
            <router-link to="/practice" class="btn btn-sm btn-primary">
              View All Sessions
            </router-link>
          </div>
        </div>
      </div>
      
      <!-- Recent progress -->
      <div class="col-md-6 col-lg-3 mb-4">
        <div class="card h-100">
          <div class="card-body">
            <h5 class="card-title">Recent Progress</h5>
            <div v-if="recentProgress.length > 0" class="progress-list">
              <div v-for="(item, index) in recentProgress" :key="index" class="progress-item">
                <div class="d-flex justify-content-between align-items-center mb-1">
                  <div class="progress-label">{{ truncateText(item.exercise.name, 20) }}</div>
                  <div class="progress-value">{{ item.bpm }} BPM</div>
                </div>
                <div class="progress mb-3">
                  <div 
                    class="progress-bar" 
                    role="progressbar" 
                    :style="{ width: `${calculateProgressPercentage(item.bpm)}%` }"
                    :aria-valuenow="item.bpm"
                    aria-valuemin="0" 
                    aria-valuemax="200"
                  ></div>
                </div>
              </div>
            </div>
            <div v-else class="text-muted">
              No recent progress data available.
            </div>
          </div>
          <div class="card-footer">
            <router-link to="/history" class="btn btn-sm btn-primary">
              View All Progress
            </router-link>
          </div>
        </div>
      </div>
      
      <!-- Recent sessions -->
      <div class="col-md-6 col-lg-3 mb-4">
        <div class="card h-100">
          <div class="card-body">
            <h5 class="card-title">Recent Sessions</h5>
            <div v-if="recentSessions.length > 0" class="session-list">
              <div v-for="session in recentSessions" :key="session.id" class="session-item">
                <router-link :to="`/practice/${session.id}`" class="session-link">
                  <div class="session-date">{{ formatDate(session.date) }}</div>
                  <div class="session-duration">{{ formatDuration(session.duration_minutes) }}</div>
                  <div v-if="session.exercises.length > 0" class="session-exercises">
                    {{ session.exercises.length }} exercises
                  </div>
                </router-link>
              </div>
            </div>
            <div v-else class="text-muted">
              No recent sessions available.
            </div>
          </div>
          <div class="card-footer">
            <button class="btn btn-sm btn-success" @click="startNewSession">
              Start New Session
            </button>
          </div>
        </div>
      </div>
      
      <!-- Quick exercises -->
      <div class="col-md-6 col-lg-3 mb-4">
        <div class="card h-100">
          <div class="card-body">
            <h5 class="card-title">Quick Exercises</h5>
            <div v-if="recentExercises.length > 0" class="exercise-list">
              <div v-for="exercise in recentExercises" :key="exercise.id" class="exercise-item">
                <router-link :to="`/exercises/${exercise.id}`" class="exercise-link">
                  <div class="exercise-name">{{ exercise.name }}</div>
                  <div class="exercise-categories">
                    <span 
                      v-for="category in exercise.categories.slice(0, 2)" 
                      :key="category.id" 
                      class="badge bg-secondary me-1"
                    >
                      {{ category.name }}
                    </span>
                  </div>
                </router-link>
              </div>
            </div>
            <div v-else class="text-muted">
              No exercises available.
            </div>
          </div>
          <div class="card-footer">
            <router-link to="/exercises" class="btn btn-sm btn-primary">
              View All Exercises
            </router-link>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Practice statistics chart -->
    <div class="row mt-4">
      <div class="col-12">
        <div class="card">
          <div class="card-body">
            <h5 class="card-title">Practice Statistics</h5>
            <div class="chart-container" style="position: relative; height: 300px;">
              <canvas ref="practiceStatsChart"></canvas>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useSessionsStore } from '@/store/modules/sessions';
import { useExercisesStore } from '@/store/modules/exercises';
import { useHistoryStore } from '@/store/modules/history';
import { getCurrentWeekRange, formatDate } from '@/utils/dateUtils';
import { formatDuration, truncateText } from '@/utils/formatters';
import Chart from 'chart.js/auto';

export default {
  name: 'Dashboard',
  
  setup() {
    const router = useRouter();
    const sessionsStore = useSessionsStore();
    const exercisesStore = useExercisesStore();
    const historyStore = useHistoryStore();
    
    const practiceStatsChart = ref(null);
    let chart = null;
    
    // Load data on component mount
    onMounted(async () => {
      // Set date range to current week
      const { startDate, endDate } = getCurrentWeekRange();
      sessionsStore.setDateRange(startDate, endDate);
      
      // Load sessions and exercises
      await Promise.all([
        sessionsStore.fetchSessions(),
        exercisesStore.fetchExercises(),
        historyStore.fetchExerciseHistory()
      ]);
      
      // Initialize chart
      initChart();
    });
    
    // Recent sessions
    const recentSessions = computed(() => {
      return sessionsStore.sessions.slice(0, 5);
    });
    
    // Recent exercises
    const recentExercises = computed(() => {
      return exercisesStore.exercises.slice(0, 5);
    });
    
    // Recent progress
    const recentProgress = computed(() => {
      return historyStore.exerciseHistory.slice(0, 3);
    });
    
    // Practice summary stats
    const recentSessionsCount = computed(() => {
      return sessionsStore.sessions.length;
    });
    
    const totalMinutesPracticed = computed(() => {
      return sessionsStore.sessions.reduce((total, session) => {
        return total + (session.duration_minutes || 0);
      }, 0);
    });
    
    const exercisesCount = computed(() => {
      return exercisesStore.totalCount;
    });
    
    // Initialize chart with practice data
    const initChart = () => {
      if (!practiceStatsChart.value) return;
      
      // Prepare data for chart
      const sessionDates = sessionsStore.sessions.map(session => 
        formatDate(session.date, 'MMM d')
      );
      
      const sessionDurations = sessionsStore.sessions.map(session => 
        session.duration_minutes || 0
      );
      
      // Create chart
      chart = new Chart(practiceStatsChart.value, {
        type: 'bar',
        data: {
          labels: sessionDates,
          datasets: [
            {
              label: 'Practice Minutes',
              data: sessionDurations,
              backgroundColor: 'rgba(54, 162, 235, 0.5)',
              borderColor: 'rgba(54, 162, 235, 1)',
              borderWidth: 1
            }
          ]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          scales: {
            y: {
              beginAtZero: true,
              title: {
                display: true,
                text: 'Minutes'
              }
            }
          }
        }
      });
    };
    
    // Calculate progress percentage for progress bars (assuming max BPM is 200)
    const calculateProgressPercentage = (bpm) => {
      return Math.min(100, (bpm / 200) * 100);
    };
    
    // Start new session handler
    const startNewSession = () => {
      // Use the navbar's modal functionality via a global event
      const event = new CustomEvent('openNewSessionModal');
      document.dispatchEvent(event);
    };
    
    return {
      recentSessions,
      recentExercises,
      recentProgress,
      recentSessionsCount,
      totalMinutesPracticed,
      exercisesCount,
      practiceStatsChart,
      calculateProgressPercentage,
      startNewSession,
      formatDate,
      formatDuration,
      truncateText
    };
  }
};
</script>

<style scoped>
.summary-stats {
  display: flex;
  justify-content: space-between;
  margin-top: 1rem;
}

.stat-item {
  text-align: center;
  padding: 0.5rem;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: #007bff;
}

.stat-label {
  font-size: 0.875rem;
  color: #6c757d;
}

.progress-list, .session-list, .exercise-list {
  margin-top: 1rem;
}

.progress-item {
  margin-bottom: 0.5rem;
}

.progress-label {
  font-size: 0.875rem;
  font-weight: 500;
}

.progress-value {
  font-size: 0.875rem;
  font-weight: 500;
  color: #28a745;
}

.session-item, .exercise-item {
  padding: 0.5rem 0;
  border-bottom: 1px solid #e9ecef;
}

.session-item:last-child, .exercise-item:last-child {
  border-bottom: none;
}

.session-link, .exercise-link {
  display: block;
  text-decoration: none;
  color: inherit;
}

.session-link:hover, .exercise-link:hover {
  background-color: rgba(0, 0, 0, 0.03);
}

.session-date, .exercise-name {
  font-weight: 500;
}

.session-duration, .session-exercises, .exercise-categories {
  font-size: 0.875rem;
  color: #6c757d;
  margin-top: 0.25rem;
}

.chart-container {
  margin-top: 1rem;
}
</style>
