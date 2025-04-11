<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
}

.stat-item {
  text-align: center;
  padding: 0.75rem;
  background-color: #f8f9fa;
  border-radius: 0.25rem;
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

.session-card {
  margin-bottom: 1rem;
  border: 1px solid #e9ecef;
  border-radius: 0.25rem;
  transition: transform 0.2s, box-shadow 0.2s;
}

.session-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.session-meta {
  font-size: 0.875rem;
}

.exercise-tags {
  margin-top: 0.5rem;
}
</style>
<template>
  <div class="practice-sessions-view">
    <div class="container-fluid">
      <div class="page-header mb-4">
        <h1>Practice Sessions</h1>
        <p class="text-muted">Track and manage your practice sessions.</p>
      </div>
      
      <!-- Show the main session list when no ID in route -->
      <div v-if="!$route.params.id">
        <!-- Date range filter -->
        <div class="card mb-4">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-center mb-3">
              <h5 class="card-title mb-0">Date Range</h5>
              <div class="btn-group">
                <button 
                  class="btn btn-sm btn-outline-secondary" 
                  @click="setThisWeek"
                >
                  This Week
                </button>
                <button 
                  class="btn btn-sm btn-outline-secondary" 
                  @click="setThisMonth"
                >
                  This Month
                </button>
                <button 
                  class="btn btn-sm btn-outline-secondary" 
                  @click="setLast30Days"
                >
                  Last 30 Days
                </button>
              </div>
            </div>
            
            <div class="row g-3">
              <div class="col-md-5">
                <div class="form-group">
                  <label for="startDate" class="form-label">Start Date</label>
                  <input 
                    type="date" 
                    class="form-control" 
                    id="startDate" 
                    v-model="startDate"
                  >
                </div>
              </div>
              <div class="col-md-5">
                <div class="form-group">
                  <label for="endDate" class="form-label">End Date</label>
                  <input 
                    type="date" 
                    class="form-control" 
                    id="endDate" 
                    v-model="endDate"
                  >
                </div>
              </div>
              <div class="col-md-2 d-flex align-items-end">
                <button 
                  class="btn btn-primary w-100" 
                  @click="loadSessions"
                >
                  Apply Filter
                </button>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Session stats -->
        <div class="row mb-4">
          <div class="col-md-4">
            <div class="card h-100">
              <div class="card-body">
                <h5 class="card-title">Practice Summary</h5>
                <div class="stats-grid">
                  <div class="stat-item">
                    <div class="stat-value">{{ sessionsStore.sessions.length }}</div>
                    <div class="stat-label">Sessions</div>
                  </div>
                  <div class="stat-item">
                    <div class="stat-value">{{ totalPracticeMinutes }}</div>
                    <div class="stat-label">Minutes</div>
                  </div>
                  <div class="stat-item">
                    <div class="stat-value">{{ averageSessionDuration }}</div>
                    <div class="stat-label">Avg. Duration</div>
                  </div>
                  <div class="stat-item">
                    <div class="stat-value">{{ totalExercisesPracticed }}</div>
                    <div class="stat-label">Exercises</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="col-md-8">
            <div class="card h-100">
              <div class="card-body">
                <h5 class="card-title">Practice Time Distribution</h5>
                <div v-if="sessionsStore.loading" class="text-center my-4">
                  <div class="spinner-border text-primary" role="status">
                    <span class="visually-hidden">Loading...</span>
                  </div>
                </div>
                <div v-else-if="sessionsStore.sessions.length === 0" class="text-center my-4">
                  <p class="text-muted">No session data available for the selected date range.</p>
                </div>
                <div v-else class="chart-container" style="position: relative; height: 200px;">
                  <canvas ref="practiceChart"></canvas>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Session list -->
        <div class="card">
          <div class="card-header d-flex justify-content-between align-items-center">
            <h5 class="mb-0">Practice Sessions</h5>
            <button class="btn btn-success" @click="openNewSessionModal">
              <i class="bi bi-plus-circle me-1"></i> New Session
            </button>
          </div>
          <div class="card-body">
            <div v-if="sessionsStore.loading" class="text-center my-4">
              <div class="spinner-border text-primary" role="status">
                <span class="visually-hidden">Loading...</span>
              </div>
              <p class="mt-2">Loading sessions...</p>
            </div>
            
            <div v-else-if="sessionsStore.sessions.length === 0" class="text-center my-4">
              <i class="bi bi-calendar-x display-4 text-muted"></i>
              <p class="mt-2">No practice sessions found in the selected date range.</p>
              <button class="btn btn-success mt-2" @click="openNewSessionModal">
                Start Your First Session
              </button>
            </div>
            
            <div v-else class="session-list">
              <div 
                v-for="session in sessionsStore.sessions" 
                :key="session.id" 
                class="session-card"
              >
                <div class="row g-0">
                  <div class="col-md-9">
                    <div class="card-body">
                      <div class="d-flex justify-content-between align-items-start">
                        <h5 class="card-title">
                          <router-link :to="`/practice/${session.id}`">
                            Practice Session - {{ formatDate(session.date) }}
                          </router-link>
                        </h5>
                        <span class="badge bg-primary">{{ formatDuration(session.duration_minutes) }}</span>
                      </div>
                      
                      <div class="session-meta mb-2">
                        <span class="text-muted">
                          <i class="bi bi-clock me-1"></i> {{ formatTime(session.date) }}
                        </span>
                        <span class="text-muted ms-3">
                          <i class="bi bi-music-note me-1"></i> 
                          {{ session.exercises?.length || 0 }} exercises
                        </span>
                      </div>
                      
                      <p class="card-text" v-if="session.notes">{{ truncateText(session.notes, 150) }}</p>
                      
                      <div v-if="session.exercises && session.exercises.length > 0" class="exercise-tags">
                        <span 
                          v-for="exercise in session.exercises.slice(0, 3)" 
                          :key="exercise.id" 
                          class="badge bg-light text-dark me-1 mb-1"
                        >
                          {{ exercise.exercise?.name || 'Unknown Exercise' }}
                        </span>
                        <span 
                          v-if="session.exercises.length > 3" 
                          class="badge bg-light text-dark me-1 mb-1"
                        >
                          +{{ session.exercises.length - 3 }} more
                        </span>
                      </div>
                    </div>
                  </div>
                  <div class="col-md-3">
                    <div class="card-body text-end">
                      <router-link 
                        :to="`/practice/${session.id}`"
                        class="btn btn-outline-primary btn-sm mb-2"
                      >
                        <i class="bi bi-eye me-1"></i> View Details
                      </router-link>
                      <button 
                        class="btn btn-outline-danger btn-sm" 
                        @click="confirmDeleteSession(session)"
                      >
                        <i class="bi bi-trash me-1"></i> Delete
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- Pagination -->
            <Pagination
              v-if="sessionsStore.totalCount > 0"
              :total-items="sessionsStore.totalCount"
              :current-page="currentPage"
              :page-size="pageSize"
              @page-change="handlePageChange"
              @page-size-change="handlePageSizeChange"
              class="mt-4"
            />
          </div>
        </div>
      </div>
      
      <!-- Router outlet for session details -->
      <router-view v-else />
      
      <!-- New Session Modal -->
      <div 
        class="modal fade" 
        id="newSessionModal" 
        tabindex="-1" 
        aria-labelledby="newSessionModalLabel" 
        aria-hidden="true"
      >
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title" id="newSessionModalLabel">Start New Practice Session</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
              <form @submit.prevent="createSession">
                <div class="mb-3">
                  <label for="sessionDate" class="form-label">Date</label>
                  <input 
                    type="datetime-local" 
                    class="form-control" 
                    id="sessionDate" 
                    v-model="newSession.date" 
                    required
                  >
                </div>
                <div class="mb-3">
                  <label for="sessionDuration" class="form-label">Duration (minutes)</label>
                  <input 
                    type="number" 
                    class="form-control" 
                    id="sessionDuration" 
                    v-model="newSession.duration_minutes" 
                    min="1" 
                    required
                  >
                </div>
                <div class="mb-3">
                  <label for="sessionNotes" class="form-label">Notes</label>
                  <textarea 
                    class="form-control" 
                    id="sessionNotes" 
                    v-model="newSession.notes" 
                    rows="3"
                  ></textarea>
                </div>
                <div class="text-end">
                  <button type="button" class="btn btn-secondary me-2" data-bs-dismiss="modal">Cancel</button>
                  <button type="submit" class="btn btn-primary" :disabled="sessionsStore.loading">
                    <span v-if="sessionsStore.loading" class="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
                    Start Session
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Delete Confirmation Modal -->
      <div 
        class="modal fade" 
        id="deleteSessionModal" 
        tabindex="-1" 
        aria-labelledby="deleteSessionModalLabel" 
        aria-hidden="true"
      >
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title" id="deleteSessionModalLabel">Confirm Delete</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
              <p>Are you sure you want to delete this practice session?</p>
              <p><strong>Date:</strong> {{ sessionToDelete ? formatDate(sessionToDelete.date) : '' }}</p>
              <p><strong>Duration:</strong> {{ sessionToDelete ? formatDuration(sessionToDelete.duration_minutes) : '' }}</p>
              <p class="text-danger">This action cannot be undone.</p>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
              <button type="button" class="btn btn-danger" @click="deleteSession" :disabled="sessionsStore.loading">
                <span v-if="sessionsStore.loading" class="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
                Delete
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import Pagination from '@/components/common/Pagination.vue';
import { useSessionsStore } from '@/store/modules/sessions';
import { formatDate, formatTime, getCurrentWeekRange, getCurrentMonthRange, getLastNDaysRange } from '@/utils/dateUtils';
import { formatDuration, truncateText } from '@/utils/formatters';
import Chart from 'chart.js/auto';
import { Modal } from 'bootstrap'

export default {
  name: 'PracticeSessionsView',
  
  components: {
    Pagination
  },
  
  setup() {
    const router = useRouter();
    const sessionsStore = useSessionsStore();
    
    // Chart ref
    const practiceChart = ref(null);
    let chart = null;
    
    // Pagination state
    const currentPage = ref(1);
    const pageSize = ref(10);
    
    // Filter state
    const startDate = ref('');
    const endDate = ref('');
    
    // New session form
    const newSession = ref({
      date: new Date().toISOString().slice(0, -8), // Format as YYYY-MM-DDTHH:MM
      duration_minutes: 30,
      notes: ''
    });
    
    // Delete confirmation
    const sessionToDelete = ref(null);
    
    // Modal references
    let sessionModal = null;
    let deleteModal = null;
    
    // Initialize component
    onMounted(() => {
      // Initialize modals
      const sessionModalEl = document.getElementById('sessionsModal')
      if (sessionModalEl) {
          sessionsModal = new Modal(sessionModalEl)
      }
      const deleteModalEl = document.getElementById('deleteModal')
      if (deleteModalEl) {
          deleteModal = new Modal(deleteModalEl)
      }
      
      // Set default date range to current month
      setThisMonth();
      
      // Add event listener for navbar "New Session" button
      document.addEventListener('openNewSessionModal', openNewSessionModal);
    });
    
    // Cleanup event listener on unmount
    onUnmounted(() => {
      document.removeEventListener('openNewSessionModal', openNewSessionModal);
    });
    
    // Watch for updates to sessions to refresh chart
    watch(() => sessionsStore.sessions, () => {
      if (sessionsStore.sessions.length > 0) {
        initPracticeChart();
      }
    });
    
    // Stats computed properties
    const totalPracticeMinutes = computed(() => {
      return sessionsStore.sessions.reduce((total, session) => {
        return total + (session.duration_minutes || 0);
      }, 0);
    });
    
    const averageSessionDuration = computed(() => {
      if (sessionsStore.sessions.length === 0) return 0;
      
      const average = totalPracticeMinutes.value / sessionsStore.sessions.length;
      return Math.round(average);
    });
    
    const totalExercisesPracticed = computed(() => {
      return sessionsStore.sessions.reduce((total, session) => {
        return total + (session.exercises?.length || 0);
      }, 0);
    });
    
    // Date filter presets
    const setThisWeek = () => {
      const { startDate: start, endDate: end } = getCurrentWeekRange();
      startDate.value = start.toISOString().split('T')[0];
      endDate.value = end.toISOString().split('T')[0];
      loadSessions();
    };
    
    const setThisMonth = () => {
      const { startDate: start, endDate: end } = getCurrentMonthRange();
      startDate.value = start.toISOString().split('T')[0];
      endDate.value = end.toISOString().split('T')[0];
      loadSessions();
    };
    
    const setLast30Days = () => {
      const { startDate: start, endDate: end } = getLastNDaysRange(30);
      startDate.value = start.toISOString().split('T')[0];
      endDate.value = end.toISOString().split('T')[0];
      loadSessions();
    };
    
    // Load sessions with date filters
    const loadSessions = () => {
      if (startDate.value) {
        sessionsStore.setDateRange(
          new Date(startDate.value),
          endDate.value ? new Date(endDate.value) : null
        );
      } else {
        sessionsStore.clearDateRange();
      }
      
      sessionsStore.fetchSessions({
        page: currentPage.value,
        page_size: pageSize.value
      }).then(() => {
        if (sessionsStore.sessions.length > 0) {
          initPracticeChart();
        }
      });
    };
    
    // Initialize chart with practice data
    const initPracticeChart = () => {
      if (!practiceChart.value) return;
      
      // Destroy existing chart if it exists
      if (chart) {
        chart.destroy();
      }
      
      // Process data for chart
      const sessions = [...sessionsStore.sessions]
        .sort((a, b) => new Date(a.date) - new Date(b.date));
      
      const dates = sessions.map(session => formatDate(session.date, 'MMM d'));
      const durations = sessions.map(session => session.duration_minutes || 0);
      
      // Create chart
      chart = new Chart(practiceChart.value, {
        type: 'bar',
        data: {
          labels: dates,
          datasets: [
            {
              label: 'Practice Duration (minutes)',
              data: durations,
              backgroundColor: 'rgba(40, 167, 69, 0.5)',
              borderColor: 'rgba(40, 167, 69, 1)',
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
    
    // Pagination handlers
    const handlePageChange = (page) => {
      currentPage.value = page;
      loadSessions();
    };
    
    const handlePageSizeChange = (size) => {
      pageSize.value = size;
      currentPage.value = 1; // Reset to first page
      loadSessions();
    };
    
    // New session modal
    const openNewSessionModal = () => {
      // Reset the form data with current date
      newSession.value = {
        date: new Date().toISOString().slice(0, -8),
        duration_minutes: 30,
        notes: ''
      };
      
      // Show the modal
      sessionModal.show();
    };
    
    const createSession = async () => {
      try {
        // Create the session
        const session = await sessionsStore.createSession(newSession.value);
        
        // Hide the modal
        sessionModal.hide();
        
        if (session && session.id) {
          // Navigate to the new session
          router.push(`/practice/${session.id}`);
        }
      } catch (error) {
        console.error('Failed to create session:', error);
      }
    };
    
    // Delete session
    const confirmDeleteSession = (session) => {
      sessionToDelete.value = session;
      deleteModal.show();
    };
    
    const deleteSession = async () => {
      if (!sessionToDelete.value) return;
      
      try {
        await sessionsStore.deleteSession(sessionToDelete.value.id);
        
        // Close modal
        deleteModal.hide();
        sessionToDelete.value = null;
        
        // Reload sessions
        loadSessions();
      } catch (error) {
        console.error('Error deleting session:', error);
      }
    };
    
    return {
      sessionsStore,
      currentPage,
      pageSize,
      startDate,
      endDate,
      practiceChart,
      newSession,
      sessionToDelete,
      totalPracticeMinutes,
      averageSessionDuration,
      totalExercisesPracticed,
      setThisWeek,
      setThisMonth,
      setLast30Days,
      loadSessions,
      handlePageChange,
      handlePageSizeChange,
      openNewSessionModal,
      createSession,
      confirmDeleteSession,
      deleteSession,
      formatDate,
      formatTime,
      formatDuration,
      truncateText
    };
  }
};
</script>
