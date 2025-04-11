<template>
  <div class="history-view">
    <div class="container-fluid">
      <div class="page-header mb-4">
        <h1>Progress Tracking</h1>
        <p class="text-muted">Track your practice progress and improvements over time.</p>
      </div>
      
      <!-- Exercise selection and filters -->
      <div class="card mb-4">
        <div class="card-body">
          <h5 class="card-title mb-3">Exercise History</h5>
          
          <div class="row g-3">
            <div class="col-md-4">
              <div class="form-group">
                <label for="exerciseSelect" class="form-label">Select Exercise</label>
                <select 
                  class="form-select" 
                  id="exerciseSelect" 
                  v-model="selectedExerciseId"
                  @change="loadExerciseHistory"
                >
                  <option value="">All Exercises</option>
                  <option 
                    v-for="exercise in exercises" 
                    :key="exercise.id" 
                    :value="exercise.id"
                  >
                    {{ exercise.name }}
                  </option>
                </select>
              </div>
            </div>
            
            <div class="col-md-4">
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
            
            <div class="col-md-4">
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
          </div>
          
          <div class="d-flex justify-content-end mt-3">
            <button 
              class="btn btn-outline-secondary me-2" 
              @click="clearFilters"
            >
              Clear Filters
            </button>
            <button 
              class="btn btn-primary" 
              @click="loadExerciseHistory"
            >
              <i class="bi bi-filter me-1"></i> Apply Filters
            </button>
          </div>
        </div>
      </div>
      
      <!-- Progress charts -->
      <div v-if="selectedExerciseId" class="row mb-4">
        <div class="col-12">
          <div class="card">
            <div class="card-body">
              <h5 class="card-title">Progress Chart</h5>
              <div v-if="historyStore.loading" class="text-center my-4">
                <div class="spinner-border text-primary" role="status">
                  <span class="visually-hidden">Loading...</span>
                </div>
              </div>
              <div v-else-if="!historyStore.exerciseHistory.length" class="text-center my-4">
                <p class="text-muted">No history data available for the selected exercise.</p>
              </div>
              <div v-else class="chart-container" style="position: relative; height: 300px;">
                <canvas ref="progressChart"></canvas>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- History list -->
      <div class="card">
        <div class="card-header d-flex justify-content-between align-items-center">
          <h5 class="mb-0">Practice History</h5>
          <div class="btn-group">
            <button 
              class="btn btn-sm btn-outline-secondary" 
              title="This Week"
              @click="setThisWeek"
            >
              This Week
            </button>
            <button 
              class="btn btn-sm btn-outline-secondary" 
              title="This Month"
              @click="setThisMonth"
            >
              This Month
            </button>
            <button 
              class="btn btn-sm btn-outline-secondary" 
              title="Last 3 Months"
              @click="setLast3Months"
            >
              Last 3 Months
            </button>
          </div>
        </div>
        <div class="card-body">
          <div v-if="historyStore.loading" class="text-center my-4">
            <div class="spinner-border text-primary" role="status">
              <span class="visually-hidden">Loading...</span>
            </div>
            <p class="mt-2">Loading history...</p>
          </div>
          
          <div v-else-if="!historyStore.exerciseHistory.length" class="text-center my-4">
            <i class="bi bi-clock-history display-4 text-muted"></i>
            <p class="mt-2">No practice history data found with the current filters.</p>
            <button class="btn btn-outline-primary" @click="clearFilters">
              Clear Filters
            </button>
          </div>
          
          <div v-else>
            <div class="table-responsive">
              <table class="table table-hover">
                <thead>
                  <tr>
                    <th>Date</th>
                    <th>Exercise</th>
                    <th>BPM</th>
                    <th>Time Signature</th>
                    <th>Rating</th>
                    <th>Notes</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="entry in historyStore.exerciseHistory" :key="entry.id">
                    <td>{{ formatDate(entry.date) }}</td>
                    <td>
                      <router-link :to="`/exercises/${entry.exercise_id}`">
                        {{ entry.exercise ? entry.exercise.name : 'Unknown Exercise' }}
                      </router-link>
                    </td>
                    <td>{{ entry.bpm }}</td>
                    <td>{{ entry.time_signature }}</td>
                    <td>
                      <div class="stars">
                        <i 
                          v-for="i in 5" 
                          :key="i" 
                          :class="['bi', i <= entry.rating ? 'bi-star-fill' : 'bi-star']"
                        ></i>
                      </div>
                    </td>
                    <td>{{ truncateText(entry.notes, 50) }}</td>
                    <td>
                      <button 
                        class="btn btn-sm btn-outline-danger" 
                        @click="confirmDeleteHistory(entry)"
                        title="Delete History Entry"
                      >
                        <i class="bi bi-trash"></i>
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            
            <!-- Pagination -->
            <Pagination
              v-if="historyStore.totalCount > 0"
              :total-items="historyStore.totalCount"
              :current-page="currentPage"
              :page-size="pageSize"
              @page-change="handlePageChange"
              @page-size-change="handlePageSizeChange"
              class="mt-4"
            />
          </div>
        </div>
      </div>
      
      <!-- Delete Confirmation Modal -->
      <div 
        class="modal fade" 
        id="deleteHistoryModal" 
        tabindex="-1" 
        aria-labelledby="deleteHistoryModalLabel" 
        aria-hidden="true"
      >
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title" id="deleteHistoryModalLabel">Confirm Delete</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
              <p>Are you sure you want to delete this practice history entry?</p>
              <p>
                <strong>Exercise:</strong> 
                {{ historyToDelete?.exercise ? historyToDelete.exercise.name : 'Unknown Exercise' }}
              </p>
              <p><strong>Date:</strong> {{ historyToDelete ? formatDate(historyToDelete.date) : '' }}</p>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
              <button type="button" class="btn btn-danger" @click="deleteHistory" :disabled="historyStore.loading">
                <span v-if="historyStore.loading" class="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
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
import { ref, onMounted, computed, watch } from 'vue';
import { useRoute } from 'vue-router';
import Pagination from '@/components/common/Pagination.vue';
import { useHistoryStore } from '@/store/modules/history';
import { useExercisesStore } from '@/store/modules/exercises';
import { formatDate } from '@/utils/dateUtils';
import { truncateText } from '@/utils/formatters';
import { getCurrentWeekRange, getCurrentMonthRange, getLastNDaysRange } from '@/utils/dateUtils';
import Chart from 'chart.js/auto';

export default {
  name: 'HistoryView',
  
  components: {
    Pagination
  },
  
  setup() {
    const route = useRoute();
    const historyStore = useHistoryStore();
    const exercisesStore = useExercisesStore();
    
    // Pagination state
    const currentPage = ref(1);
    const pageSize = ref(10);
    
    // Filter state
    const selectedExerciseId = ref('');
    const startDate = ref('');
    const endDate = ref('');
    
    // Chart ref
    const progressChart = ref(null);
    let chart = null;
    
    // Delete confirmation
    const historyToDelete = ref(null);
    let deleteModal = null;
    
    // Check URL params for exercise_id
    onMounted(async () => {
      deleteModal = new window.bootstrap.Modal(document.getElementById('deleteHistoryModal'));
      
      // Load exercises for the dropdown
      if (exercisesStore.exercises.length === 0) {
        await exercisesStore.fetchExercises({ page_size: 100 });
      }
      
      // Check if there's an exercise_id in the route query
      if (route.query.exercise_id) {
        selectedExerciseId.value = parseInt(route.query.exercise_id, 10);
      }
      
      // Set default date range to last month
      const { startDate: start, endDate: end } = getLastNDaysRange(30);
      startDate.value = start.toISOString().split('T')[0];
      endDate.value = end.toISOString().split('T')[0];
      
      // Load initial history
      await loadExerciseHistory();
    });
    
    // Watch for changes to chart data
    watch(() => historyStore.exerciseHistory, () => {
      if (selectedExerciseId.value && historyStore.exerciseHistory.length > 1) {
        initProgressChart();
      }
    });
    
    // Get exercises from store
    const exercises = computed(() => {
      return exercisesStore.exercises;
    });
    
    // Load exercise history
    const loadExerciseHistory = async () => {
      historyStore.setExerciseId(selectedExerciseId.value || null);
      
      if (startDate.value) {
        historyStore.setDateRange(
          new Date(startDate.value), 
          endDate.value ? new Date(endDate.value) : null
        );
      } else {
        historyStore.clearDateRange();
      }
      
      await historyStore.fetchExerciseHistory({
        page: currentPage.value,
        page_size: pageSize.value
      });
      
      // Initialize chart if we have data and a selected exercise
      if (selectedExerciseId.value && historyStore.exerciseHistory.length > 1) {
        initProgressChart();
      }
    };
    
    // Filter date presets
    const setThisWeek = () => {
      const { startDate: start, endDate: end } = getCurrentWeekRange();
      startDate.value = start.toISOString().split('T')[0];
      endDate.value = end.toISOString().split('T')[0];
      loadExerciseHistory();
    };
    
    const setThisMonth = () => {
      const { startDate: start, endDate: end } = getCurrentMonthRange();
      startDate.value = start.toISOString().split('T')[0];
      endDate.value = end.toISOString().split('T')[0];
      loadExerciseHistory();
    };
    
    const setLast3Months = () => {
      const { startDate: start, endDate: end } = getLastNDaysRange(90);
      startDate.value = start.toISOString().split('T')[0];
      endDate.value = end.toISOString().split('T')[0];
      loadExerciseHistory();
    };
    
    // Clear all filters
    const clearFilters = () => {
      selectedExerciseId.value = '';
      startDate.value = '';
      endDate.value = '';
      historyStore.setExerciseId(null);
      historyStore.clearDateRange();
      currentPage.value = 1;
      loadExerciseHistory();
    };
    
    // Pagination handlers
    const handlePageChange = (page) => {
      currentPage.value = page;
      loadExerciseHistory();
    };
    
    const handlePageSizeChange = (size) => {
      pageSize.value = size;
      currentPage.value = 1; // Reset to first page
      loadExerciseHistory();
    };
    
    // Initialize progress chart
    const initProgressChart = () => {
      if (!progressChart.value) return;
      
      // Destroy existing chart if it exists
      if (chart) {
        chart.destroy();
      }
      
      // Process data for chart - sort by date
      const entries = [...historyStore.exerciseHistory]
        .sort((a, b) => new Date(a.date) - new Date(b.date));
      
      const dates = entries.map(entry => formatDate(entry.date, 'MMM d'));
      const bpms = entries.map(entry => entry.bpm);
      
      // Create chart
      chart = new Chart(progressChart.value, {
        type: 'line',
        data: {
          labels: dates,
          datasets: [
            {
              label: 'BPM',
              data: bpms,
              borderColor: 'rgba(54, 162, 235, 1)',
              backgroundColor: 'rgba(54, 162, 235, 0.2)',
              borderWidth: 2,
              tension: 0.3,
              fill: true
            }
          ]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          scales: {
            y: {
              beginAtZero: false,
              title: {
                display: true,
                text: 'BPM'
              }
            }
          }
        }
      });
    };
    
    // Delete history
    const confirmDeleteHistory = (history) => {
      historyToDelete.value = history;
      deleteModal.show();
    };
    
    const deleteHistory = async () => {
      if (!historyToDelete.value) return;
      
      try {
        await historyStore.deleteExerciseHistory(historyToDelete.value.id);
        
        // Close modal
        deleteModal.hide();
        historyToDelete.value = null;
        
        // Reload history
        await loadExerciseHistory();
      } catch (error) {
        console.error('Error deleting history:', error);
      }
    };
    
    return {
      historyStore,
      exercises,
      selectedExerciseId,
      startDate,
      endDate,
      currentPage,
      pageSize,
      progressChart,
      historyToDelete,
      loadExerciseHistory,
      setThisWeek,
      setThisMonth,
      setLast3Months,
      clearFilters,
      handlePageChange,
      handlePageSizeChange,
      confirmDeleteHistory,
      deleteHistory,
      formatDate,
      truncateText
    };
  }
};
</script>

<style scoped>
.stars {
  color: #ffc107;
  font-size: 0.875rem;
}

.chart-container {
  margin-top: 1rem;
}

.table td {
  vertical-align: middle;
}
</style>
