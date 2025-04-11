<template>
  <div class="session-details">
    <!-- Back navigation -->
    <div class="mb-4">
      <router-link to="/practice" class="btn btn-outline-secondary">
        <i class="bi bi-arrow-left me-1"></i> Back to Sessions
      </router-link>
    </div>
    
    <!-- Loading indicator -->
    <div v-if="sessionsStore.loading" class="text-center my-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
      <p class="mt-2">Loading session details...</p>
    </div>
    
    <!-- Error message -->
    <div v-else-if="sessionsStore.error" class="alert alert-danger" role="alert">
      <i class="bi bi-exclamation-triangle-fill me-2"></i>
      {{ sessionsStore.error }}
    </div>
    
    <!-- Not found message -->
    <div v-else-if="!session" class="alert alert-warning" role="alert">
      <i class="bi bi-question-circle-fill me-2"></i>
      Session not found. It may have been deleted or you may have followed an invalid link.
    </div>
    
    <!-- Session details -->
    <div v-else class="session-content">
      <div class="row">
        <!-- Left column: Session info -->
        <div class="col-md-8">
          <div class="card mb-4">
            <div class="card-body">
              <div class="d-flex justify-content-between align-items-start mb-3">
                <h1 class="card-title mb-0">Practice Session</h1>
                <div class="button-group">
                  <button class="btn btn-outline-primary" @click="editSession">
                    <i class="bi bi-pencil me-1"></i> Edit
                  </button>
                </div>
              </div>
              
              <div class="session-meta mb-4">
                <div class="meta-item">
                  <i class="bi bi-calendar me-1"></i> 
                  <strong>Date:</strong> {{ formatDate(session.date) }}
                </div>
                <div class="meta-item">
                  <i class="bi bi-clock me-1"></i> 
                  <strong>Time:</strong> {{ formatTime(session.date) }}
                </div>
                <div class="meta-item">
                  <i class="bi bi-stopwatch me-1"></i> 
                  <strong>Duration:</strong> {{ formatDuration(session.durationMinutes) }}
                </div>
              </div>
              
              <!-- Notes -->
              <div class="notes mb-4">
                <h5>Notes</h5>
                <p v-if="session.notes" class="mb-0">{{ session.notes }}</p>
                <p v-else class="text-muted mb-0">No notes recorded for this session.</p>
              </div>
              
              <!-- Exercises -->
              <div class="exercises mb-4">
                <div class="d-flex justify-content-between align-items-center mb-3">
                  <h5 class="mb-0">Exercises Practiced</h5>
                  <button class="btn btn-sm btn-primary" @click="openAddExerciseModal">
                    <i class="bi bi-plus-circle me-1"></i> Add Exercise
                  </button>
                </div>
                
                <div v-if="!session.exercises || session.exercises.length === 0" class="text-muted">
                  No exercises added to this session yet.
                </div>
                
                <div v-else class="exercise-list">
                  <div 
                    v-for="(sessionExercise, index) in session.exercises" 
                    :key="sessionExercise.id" 
                    class="exercise-item"
                    :class="{ 'border-bottom': index < session.exercises.length - 1 }"
                  >
                    <div class="d-flex justify-content-between align-items-start">
                      <div class="exercise-info">
                        <h5 class="exercise-name">
                          <router-link :to="`/exercises/${sessionExercise.exercise_id}`">
                            {{ sessionExercise.exercise?.name || 'Unknown Exercise' }}
                          </router-link>
                        </h5>
                        <div class="exercise-meta">
                          <span class="me-3">
                            <i class="bi bi-stopwatch me-1"></i> {{ formatDuration(sessionExercise.durationMinutes) }}
                          </span>
                          <span class="me-3">
                            <i class="bi bi-music-note me-1"></i> {{ sessionExercise.bpm }} BPM
                          </span>
                          <span>
                            <i class="bi bi-clock me-1"></i> {{ sessionExercise.time_signature || '4/4' }}
                          </span>
                        </div>
                        <p v-if="sessionExercise.notes" class="exercise-notes mt-2">
                          {{ sessionExercise.notes }}
                        </p>
                      </div>
                      <div class="exercise-actions">
                        <button 
                          class="btn btn-sm btn-outline-primary me-1" 
                          @click="editExercise(sessionExercise)"
                          title="Edit Exercise"
                        >
                          <i class="bi bi-pencil"></i>
                        </button>
                        <button 
                          class="btn btn-sm btn-outline-danger" 
                          @click="confirmRemoveExercise(sessionExercise)"
                          title="Remove Exercise"
                        >
                          <i class="bi bi-trash"></i>
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Right column: Session summary -->
        <div class="col-md-4">
          <div class="card mb-4">
            <div class="card-header">
              <h5 class="mb-0">Session Summary</h5>
            </div>
            <div class="card-body">
              <div class="stats-grid">
                <div class="stat-item">
                  <div class="stat-value">{{ session.durationMinutes }}</div>
                  <div class="stat-label">Minutes</div>
                </div>
                <div class="stat-item">
                  <div class="stat-value">{{ session.exercises?.length || 0 }}</div>
                  <div class="stat-label">Exercises</div>
                </div>
                <div class="stat-item">
                  <div class="stat-value">{{ averageBpm }}</div>
                  <div class="stat-label">Avg. BPM</div>
                </div>
                <div class="stat-item">
                  <div class="stat-value">{{ exerciseDurationPercentage }}%</div>
                  <div class="stat-label">Time Tracked</div>
                </div>
              </div>
              
              <!-- Exercise time distribution -->
              <div v-if="session.exercises && session.exercises.length > 0" class="time-distribution mt-4">
                <h6>Time Distribution</h6>
                <div class="chart-container" style="position: relative; height: 200px;">
                  <canvas ref="timeDistributionChart"></canvas>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Edit Session Modal -->
    <div 
      class="modal fade" 
      id="editSessionModal" 
      tabindex="-1" 
      aria-labelledby="editSessionModalLabel" 
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="editSessionModalLabel">Edit Session</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="updateSession">
              <div class="mb-3">
                <label for="sessionDate" class="form-label">Date</label>
                <input 
                  type="datetime-local" 
                  class="form-control" 
                  id="sessionDate" 
                  v-model="sessionForm.date" 
                  required
                >
              </div>
              <div class="mb-3">
                <label for="sessionDuration" class="form-label">Duration (minutes)</label>
                <input 
                  type="number" 
                  class="form-control" 
                  id="sessionDuration" 
                  v-model="sessionForm.durationMinutes" 
                  min="1" 
                  required
                >
              </div>
              <div class="mb-3">
                <label for="sessionNotes" class="form-label">Notes</label>
                <textarea 
                  class="form-control" 
                  id="sessionNotes" 
                  v-model="sessionForm.notes" 
                  rows="3"
                ></textarea>
              </div>
              <div class="text-end">
                <button type="button" class="btn btn-secondary me-2" data-bs-dismiss="modal">Cancel</button>
                <button type="submit" class="btn btn-primary" :disabled="sessionsStore.loading">
                  <span v-if="sessionsStore.loading" class="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
                  Update Session
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Add Exercise Modal -->
    <div 
      class="modal fade" 
      id="addExerciseModal" 
      tabindex="-1" 
      aria-labelledby="addExerciseModalLabel" 
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="addExerciseModalLabel">Add Exercise</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="addExercise">
              <div class="mb-3">
                <label for="exerciseSelect" class="form-label">Select Exercise</label>
                <select 
                  class="form-select" 
                  id="exerciseSelect" 
                  v-model="exerciseForm.exercise_id" 
                  required
                >
                  <option value="" disabled>Choose an exercise</option>
                  <option 
                    v-for="exercise in availableExercises" 
                    :key="exercise.id" 
                    :value="exercise.id"
                  >
                    {{ exercise.name }}
                  </option>
                </select>
              </div>
              <div class="mb-3">
                <label for="exerciseDuration" class="form-label">Duration (minutes)</label>
                <input 
                  type="number" 
                  class="form-control" 
                  id="exerciseDuration" 
                  v-model="exerciseForm.durationMinutes" 
                  min="1" 
                  required
                >
              </div>
              <div class="mb-3">
                <label for="exerciseBpm" class="form-label">BPM</label>
                <input 
                  type="number" 
                  class="form-control" 
                  id="exerciseBpm" 
                  v-model="exerciseForm.bpm" 
                  min="1"
                >
              </div>
              <div class="mb-3">
                <label for="exerciseTimeSignature" class="form-label">Time Signature</label>
                <select 
                  class="form-select" 
                  id="exerciseTimeSignature" 
                  v-model="exerciseForm.time_signature"
                >
                  <option value="4/4">4/4 (Common Time)</option>
                  <option value="3/4">3/4 (Waltz Time)</option>
                  <option value="2/4">2/4 (March Time)</option>
                  <option value="6/8">6/8 (Compound Duple)</option>
                  <option value="9/8">9/8 (Compound Triple)</option>
                  <option value="12/8">12/8 (Compound Quadruple)</option>
                  <option value="5/4">5/4 (Quintuple)</option>
                  <option value="7/8">7/8 (Septuple)</option>
                </select>
              </div>
              <div class="mb-3">
                <label for="exerciseNotes" class="form-label">Notes</label>
                <textarea 
                  class="form-control" 
                  id="exerciseNotes" 
                  v-model="exerciseForm.notes" 
                  rows="3"
                ></textarea>
              </div>
              <div class="text-end">
                <button type="button" class="btn btn-secondary me-2" data-bs-dismiss="modal">Cancel</button>
                <button type="submit" class="btn btn-primary" :disabled="sessionsStore.loading">
                  <span v-if="sessionsStore.loading" class="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
                  Add Exercise
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Edit Exercise Modal -->
    <div 
      class="modal fade" 
      id="editExerciseModal" 
      tabindex="-1" 
      aria-labelledby="editExerciseModalLabel" 
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="editExerciseModalLabel">Edit Exercise</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="updateExercise">
              <div class="mb-3">
                <label for="editExerciseDuration" class="form-label">Duration (minutes)</label>
                <input 
                  type="number" 
                  class="form-control" 
                  id="editExerciseDuration" 
                  v-model="exerciseForm.durationMinutes" 
                  min="1" 
                  required
                >
              </div>
              <div class="mb-3">
                <label for="editExerciseBpm" class="form-label">BPM</label>
                <input 
                  type="number" 
                  class="form-control" 
                  id="editExerciseBpm" 
                  v-model="exerciseForm.bpm" 
                  min="1"
                >
              </div>
              <div class="mb-3">
                <label for="editExerciseTimeSignature" class="form-label">Time Signature</label>
                <select 
                  class="form-select" 
                  id="editExerciseTimeSignature" 
                  v-model="exerciseForm.time_signature"
                >
                  <option value="4/4">4/4 (Common Time)</option>
                  <option value="3/4">3/4 (Waltz Time)</option>
                  <option value="2/4">2/4 (March Time)</option>
                  <option value="6/8">6/8 (Compound Duple)</option>
                  <option value="9/8">9/8 (Compound Triple)</option>
                  <option value="12/8">12/8 (Compound Quadruple)</option>
                  <option value="5/4">5/4 (Quintuple)</option>
                  <option value="7/8">7/8 (Septuple)</option>
                </select>
              </div>
              <div class="mb-3">
                <label for="editExerciseNotes" class="form-label">Notes</label>
                <textarea 
                  class="form-control" 
                  id="editExerciseNotes" 
                  v-model="exerciseForm.notes" 
                  rows="3"
                ></textarea>
              </div>
              <div class="text-end">
                <button type="button" class="btn btn-secondary me-2" data-bs-dismiss="modal">Cancel</button>
                <button type="submit" class="btn btn-primary" :disabled="sessionsStore.loading">
                  <span v-if="sessionsStore.loading" class="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
                  Update Exercise
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Remove Exercise Confirmation Modal -->
    <div 
      class="modal fade" 
      id="removeExerciseModal" 
      tabindex="-1" 
      aria-labelledby="removeExerciseModalLabel" 
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="removeExerciseModalLabel">Confirm Remove</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <p>Are you sure you want to remove this exercise from the session?</p>
            <p>
              <strong>Exercise:</strong> 
              {{ exerciseToRemove?.exercise?.name || 'Unknown Exercise' }}
            </p>
            <p>
              <strong>Duration:</strong> 
              {{ exerciseToRemove ? formatDuration(exerciseToRemove.durationMinutes) : '' }}
            </p>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
            <button type="button" class="btn btn-danger" @click="removeExercise" :disabled="sessionsStore.loading">
              <span v-if="sessionsStore.loading" class="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
              Remove
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue';
import { useSessionsStore } from '@/store/modules/sessions';
import { useExercisesStore } from '@/store/modules/exercises';
import { formatDate, formatTime } from '@/utils/dateUtils';
import { formatDuration } from '@/utils/formatters';
import Chart from 'chart.js/auto';
import { Modal } from 'bootstrap'

export default {
  name: 'SessionDetails',
  
  props: {
    sessionId: {
      type: Number,
      required: true
    }
  },
  
  setup(props) {
    const sessionsStore = useSessionsStore();
    const exercisesStore = useExercisesStore();
    
    const timeDistributionChart = ref(null);
    let chart = null;
    
    // Form state for session editing
    const sessionForm = ref({
      id: null,
      date: '',
      durationMinutes: 0,
      notes: ''
    });
    
    // Form state for exercise
    const exerciseForm = ref({
      id: null,
      session_id: null,
      exercise_id: '',
      durationMinutes: 15,
      bpm: 100,
      time_signature: '4/4',
      notes: ''
    });
    
    // State for exercise removal
    const exerciseToRemove = ref(null);
    
    // Modal references
    let editSessionModal = null;
    let addExerciseModal = null;
    let editExerciseModal = null;
    let removeExerciseModal = null;
    
    // Load data on component mount
    onMounted(async () => {
      // Initialize modals
      const editSessionModalEl = document.getElementById('editSessionModal')
      if (editSessionModalEl) {
          editSessionModal = new Modal(editSessionModalEl)
      }
      const addExerciseModalEl = document.getElementById('addExerciseModal')
      if (addExerciseModalEl) {
          addExerciseModal = new Modal(addExerciseModalEl)
      }
      const editExerciseModalEl = document.getElementById('editExerciseModal')
      if (editExerciseModalEl) {
          editExerciseModal = new Modal(editExerciseModalEl)
      }
      const removeExerciseModalEl = document.getElementById('removeExerciseModal')
      if (removeExerciseModalEl) {
          removeExerciseModal = new Modal(removeExerciseModalEl)
      }
      
      // Load exercises if not already loaded
      if (exercisesStore.exercises.length === 0) {
        await exercisesStore.fetchExercises({ page_size: 100 });
      }
      
      // Fetch session details
      await sessionsStore.fetchSessionById(props.sessionId);
      
      // Initialize chart
      if (session.value && session.value.exercises && session.value.exercises.length > 0) {
        initTimeDistributionChart();
      }
    });
    
    // Watch for session ID changes
    watch(() => props.sessionId, async (newId) => {
      if (newId) {
        await sessionsStore.fetchSessionById(newId);
        
        // Initialize chart if needed
        if (session.value && session.value.exercises && session.value.exercises.length > 0) {
          if (chart) {
            chart.destroy();
          }
          initTimeDistributionChart();
        }
      }
    });
    
    // Watch for changes to current session to update chart
    watch(() => sessionsStore.currentSession, () => {
      if (session.value && session.value.exercises && session.value.exercises.length > 0) {
        if (chart) {
          chart.destroy();
        }
        initTimeDistributionChart();
      }
    }, { deep: true });
    
    // Access the current session
    const session = computed(() => {
      return sessionsStore.currentSession;
    });
    
    // Compute available exercises (those not already in the session)
    const availableExercises = computed(() => {
      if (!session.value || !session.value.exercises) {
        return exercisesStore.exercises;
      }
      
      // Get IDs of exercises already in the session
      const existingExerciseIds = session.value.exercises.map(se => se.exercise_id);
      
      // Filter out exercises already in the session
      return exercisesStore.exercises.filter(exercise => 
        !existingExerciseIds.includes(exercise.id)
      );
    });
    
    // Calculate average BPM across exercises
    const averageBpm = computed(() => {
      if (!session.value || !session.value.exercises || session.value.exercises.length === 0) {
        return 0;
      }
      
      const exercisesWithBpm = session.value.exercises.filter(se => se.bpm);
      if (exercisesWithBpm.length === 0) {
        return 0;
      }
      
      const totalBpm = exercisesWithBpm.reduce((sum, se) => sum + se.bpm, 0);
      return Math.round(totalBpm / exercisesWithBpm.length);
    });
    
    // Calculate percentage of session time tracked with exercises
    const exerciseDurationPercentage = computed(() => {
      if (!session.value || !session.value.exercises || session.value.exercises.length === 0 || !session.value.durationMinutes) {
        return 0;
      }
      
      const totalExerciseDuration = session.value.exercises.reduce((sum, se) => 
        sum + (se.durationMinutes || 0), 0
      );
      
      const percentage = (totalExerciseDuration / session.value.durationMinutes) * 100;
      return Math.min(100, Math.round(percentage)); // Cap at 100%
    });
    
    // Initialize time distribution chart
    const initTimeDistributionChart = () => {
      if (!timeDistributionChart.value || !session.value || !session.value.exercises) return;
      
      // Process data for chart
      const exercises = session.value.exercises.filter(se => se.durationMinutes);
      
      // Group small exercises under "Other" to avoid cluttering the chart
      const threshold = session.value.durationMinutes * 0.05; // 5% of total time
      let otherTime = 0;
      
      const chartData = exercises.reduce((result, se) => {
        if (se.durationMinutes >= threshold) {
          result.labels.push(se.exercise?.name || 'Unknown');
          result.data.push(se.durationMinutes);
        } else {
          otherTime += se.durationMinutes;
        }
        return result;
      }, { labels: [], data: [] });
      
      // Add "Other" category if needed
      if (otherTime > 0) {
        chartData.labels.push('Other');
        chartData.data.push(otherTime);
      }
      
      // Generate colors
      const backgroundColors = generateChartColors(chartData.labels.length);
      
      // Create chart
      chart = new Chart(timeDistributionChart.value, {
        type: 'pie',
        data: {
          labels: chartData.labels,
          datasets: [
            {
              data: chartData.data,
              backgroundColor: backgroundColors
            }
          ]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            legend: {
              position: 'right',
              labels: {
                boxWidth: 15,
                font: {
                  size: 12
                }
              }
            },
            tooltip: {
              callbacks: {
                label: (context) => {
                  const value = context.raw;
                  const total = context.dataset.data.reduce((sum, val) => sum + val, 0);
                  const percentage = Math.round((value / total) * 100);
                  return `${context.label}: ${value} mins (${percentage}%)`;
                }
              }
            }
          }
        }
      });
    };
    
    // Generate colors for chart
    const generateChartColors = (count) => {
      const baseColors = [
        'rgba(54, 162, 235, 0.8)',   // Blue
        'rgba(255, 99, 132, 0.8)',   // Red
        'rgba(75, 192, 192, 0.8)',   // Green
        'rgba(255, 159, 64, 0.8)',   // Orange
        'rgba(153, 102, 255, 0.8)',  // Purple
        'rgba(255, 205, 86, 0.8)',   // Yellow
        'rgba(201, 203, 207, 0.8)'   // Gray
      ];
      
      // If we need more colors than in the base array, generate them
      if (count <= baseColors.length) {
        return baseColors.slice(0, count);
      }
      
      // Generate additional colors
      const colors = [...baseColors];
      for (let i = baseColors.length; i < count; i++) {
        const r = Math.floor(Math.random() * 255);
        const g = Math.floor(Math.random() * 255);
        const b = Math.floor(Math.random() * 255);
        colors.push(`rgba(${r}, ${g}, ${b}, 0.8)`);
      }
      
      return colors;
    };
    
    // Edit session
    const editSession = () => {
      if (!session.value) return;
      
      // Format date for input
      const date = new Date(session.value.date);
      const formattedDate = date.toISOString().slice(0, -8); // Format as YYYY-MM-DDTHH:MM
      
      sessionForm.value = {
        id: session.value.id,
        date: formattedDate,
        durationMinutes: session.value.durationMinutes,
        notes: session.value.notes || ''
      };
      
      editSessionModal.show();
    };
    
    // Update session
    const updateSession = async () => {
      try {
        await sessionsStore.updateSession(sessionForm.value);
        
        // Close modal
        editSessionModal.hide();
      } catch (error) {
        console.error('Error updating session:', error);
      }
    };
    
    // Open add exercise modal
    const openAddExerciseModal = () => {
      exerciseForm.value = {
        id: null,
        session_id: props.sessionId,
        exercise_id: '',
        durationMinutes: 15,
        bpm: 100,
        time_signature: '4/4',
        notes: ''
      };
      
      addExerciseModal.show();
    };
    
    // Add exercise to session
    const addExercise = async () => {
      try {
        await sessionsStore.addExerciseToSession(props.sessionId, exerciseForm.value);
        
        // Close modal and reset form
        addExerciseModal.hide();
        exerciseForm.value = {
          id: null,
          session_id: props.sessionId,
          exercise_id: '',
          durationMinutes: 15,
          bpm: 100,
          time_signature: '4/4',
          notes: ''
        };
        
        // Refresh session data
        await sessionsStore.fetchSessionById(props.sessionId);
      } catch (error) {
        console.error('Error adding exercise to session:', error);
      }
    };
    
    // Edit exercise
    const editExercise = (sessionExercise) => {
      exerciseForm.value = {
        id: sessionExercise.id,
        session_id: sessionExercise.session_id,
        exercise_id: sessionExercise.exercise_id,
        durationMinutes: sessionExercise.durationMinutes || 15,
        bpm: sessionExercise.bpm || 100,
        time_signature: sessionExercise.time_signature || '4/4',
        notes: sessionExercise.notes || ''
      };
      
      editExerciseModal.show();
    };
    
    // Update exercise
    const updateExercise = async () => {
      try {
        await sessionsStore.updateSessionExercise(exerciseForm.value);
        
        // Close modal
        editExerciseModal.hide();
        
        // Refresh session data
        await sessionsStore.fetchSessionById(props.sessionId);
      } catch (error) {
        console.error('Error updating exercise:', error);
      }
    };
    
    // Confirm remove exercise
    const confirmRemoveExercise = (sessionExercise) => {
      exerciseToRemove.value = sessionExercise;
      removeExerciseModal.show();
    };
    
    // Remove exercise from session
    const removeExercise = async () => {
      if (!exerciseToRemove.value) return;
      
      try {
        await sessionsStore.removeExerciseFromSession(
          exerciseToRemove.value.sessionId,
          exerciseToRemove.value.exerciseId
        );
        
        // Close modal
        removeExerciseModal.hide();
        exerciseToRemove.value = null;
        
        // Refresh session data
        await sessionsStore.fetchSessionById(props.sessionId);
      } catch (error) {
        console.error('Error removing exercise from session:', error);
      }
    };
    
    return {
      sessionsStore,
      exercisesStore,
      session,
      availableExercises,
      averageBpm,
      exerciseDurationPercentage,
      timeDistributionChart,
      sessionForm,
      exerciseForm,
      exerciseToRemove,
      editSession,
      updateSession,
      openAddExerciseModal,
      addExercise,
      editExercise,
      updateExercise,
      confirmRemoveExercise,
      removeExercise,
      formatDate,
      formatTime,
      formatDuration
    };
  }
};
</script>

<style scoped>
.session-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  font-size: 1rem;
}

.meta-item {
  padding: 0.5rem 1rem;
  background-color: #f8f9fa;
  border-radius: 0.25rem;
}

.exercise-list {
  margin-top: 1rem;
}

.exercise-item {
  padding: 1rem 0;
  border-bottom: 1px solid #e9ecef;
}

.exercise-name {
  margin-bottom: 0.5rem;
}

.exercise-meta {
  font-size: 0.875rem;
  color: #6c757d;
}

.exercise-notes {
  font-size: 0.875rem;
  color: #495057;
  margin-bottom: 0;
}

.exercise-actions {
  display: flex;
}

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
</style>
