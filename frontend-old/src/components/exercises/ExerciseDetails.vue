<style scoped>
.meta-info {
  font-size: 0.875rem;
}

.description {
  white-space: pre-line;
}

.music-image-container {
  margin-top: 1rem;
  border: 1px solid #dee2e6;
  padding: 1rem;
  border-radius: 0.25rem;
  background-color: #f8f9fa;
  text-align: center;
}

.link-list {
  margin-top: 1rem;
}

.link-item {
  padding: 0.75rem 0;
  border-bottom: 1px solid #f0f0f0;
}

.link-item:last-child {
  border-bottom: none;
}

.link-url {
  text-decoration: none;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 90%;
  display: inline-block;
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

.last-practice-info {
  margin-top: 0.5rem;
}

.info-row {
  display: flex;
  margin-bottom: 0.5rem;
}

.info-label {
  font-weight: 500;
  width: 120px;
}

.info-row.notes {
  flex-direction: column;
}

.info-row.notes .info-label {
  margin-bottom: 0.25rem;
}

.stars {
  color: #ffc107;
}

.history-list {
  max-height: 300px;
  overflow-y: auto;
}

.history-item {
  padding: 0.75rem 1rem;
}

.history-date {
  font-weight: 500;
}

.history-bpm, .history-time-sig {
  font-size: 0.875rem;
  color: #6c757d;
}

.history-rating {
  font-size: 0.75rem;
}

.rating-input {
  display: flex;
  align-items: center;
}

.rating-text {
  margin-left: 1rem;
  font-size: 0.875rem;
  color: #6c757d;
}
</style>
<template>
  <div class="exercise-details">
    <!-- Back navigation -->
    <div class="mb-4">
      <router-link to="/exercises" class="btn btn-outline-secondary">
        <i class="bi bi-arrow-left me-1"></i> Back to Exercises
      </router-link>
    </div>
    
    <!-- Loading indicator -->
    <div v-if="exercisesStore.loading" class="text-center my-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
      <p class="mt-2">Loading exercise details...</p>
    </div>
    
    <!-- Error message -->
    <div v-else-if="exercisesStore.error" class="alert alert-danger" role="alert">
      <i class="bi bi-exclamation-triangle-fill me-2"></i>
      {{ exercisesStore.error }}
    </div>
    
    <!-- Not found message -->
    <div v-else-if="!exercise" class="alert alert-warning" role="alert">
      <i class="bi bi-question-circle-fill me-2"></i>
      Exercise not found. It may have been deleted or you may have followed an invalid link.
    </div>
    
    <!-- Exercise details -->
    <div v-else class="exercise-content">
      <div class="row">
        <!-- Left column: Exercise info -->
        <div class="col-md-8">
          <div class="card mb-4">
            <div class="card-body">
              <div class="d-flex justify-content-between align-items-start mb-3">
                <h1 class="card-title mb-0">{{ exercise.name }}</h1>
                <div class="button-group">
                  <button class="btn btn-primary me-2" @click="practiceExercise">
                    <i class="bi bi-play-circle me-1"></i> Practice Now
                  </button>
                  <button class="btn btn-outline-primary" @click="editExercise">
                    <i class="bi bi-pencil me-1"></i> Edit
                  </button>
                </div>
              </div>
              
              <div class="meta-info mb-3">
                <span class="text-muted me-3">
                  <i class="bi bi-calendar me-1"></i> Added: {{ formatDate(exercise.created_at) }}
                </span>
                <span class="text-muted">
                  <i class="bi bi-clock-history me-1"></i> Updated: {{ formatDate(exercise.updated_at) }}
                </span>
              </div>
              
              <!-- Categories and tags -->
              <div class="mb-4">
                <div class="mb-2">
                  <span 
                    v-for="category in exercise.categories" 
                    :key="category.id" 
                    class="badge bg-primary me-1 mb-1"
                  >
                    <i class="bi bi-folder me-1"></i> {{ category.name }}
                  </span>
                </div>
                <div>
                  <span 
                    v-for="tag in exercise.tags" 
                    :key="tag.id" 
                    class="badge bg-secondary me-1 mb-1"
                  >
                    <i class="bi bi-tag me-1"></i> {{ tag.name }}
                  </span>
                </div>
              </div>
              
              <!-- Description -->
              <div class="description mb-4">
                <h5>Description</h5>
                <p v-if="exercise.description" class="mb-0">{{ exercise.description }}</p>
                <p v-else class="text-muted mb-0">No description provided.</p>
              </div>
              
              <!-- Sheet music -->
              <div v-if="exercise.music_image_path" class="sheet-music mb-4">
                <h5>Sheet Music</h5>
                <div class="music-image-container">
                  <img :src="exercise.music_image_path" alt="Sheet music" class="img-fluid" />
                </div>
              </div>
              
              <!-- Links -->
              <div class="links mb-4">
                <h5 class="d-flex justify-content-between align-items-center">
                  <span>Reference Links</span>
                  <button class="btn btn-sm btn-outline-primary" @click="openAddLinkModal">
                    <i class="bi bi-plus-circle me-1"></i> Add Link
                  </button>
                </h5>
                
                <div v-if="exercise.links && exercise.links.length > 0" class="link-list">
                  <div v-for="link in exercise.links" :key="link.id" class="link-item">
                    <div class="d-flex justify-content-between align-items-center">
                      <a :href="link.url" target="_blank" rel="noopener noreferrer" class="link-url">
                        <i class="bi bi-link-45deg me-1"></i> {{ link.description || link.url }}
                      </a>
                      <button 
                        class="btn btn-sm btn-outline-danger" 
                        @click="confirmDeleteLink(link)"
                        title="Delete Link"
                      >
                        <i class="bi bi-trash"></i>
                      </button>
                    </div>
                  </div>
                </div>
                <div v-else class="text-muted">
                  No reference links added yet.
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Right column: Progress and actions -->
        <div class="col-md-4">
          <!-- Practice stats -->
          <div class="card mb-4">
            <div class="card-header">
              <h5 class="mb-0">Practice Stats</h5>
            </div>
            <div class="card-body">
              <div v-if="loadingHistory" class="text-center py-3">
                <div class="spinner-border spinner-border-sm text-primary" role="status">
                  <span class="visually-hidden">Loading...</span>
                </div>
                <p class="mt-2 mb-0">Loading history...</p>
              </div>
              <div v-else-if="exerciseHistory.length === 0" class="text-center py-3">
                <i class="bi bi-graph-up display-5 text-muted"></i>
                <p class="mt-2 mb-0">No practice data yet. Start practicing to track your progress!</p>
              </div>
              <div v-else class="stats-content">
                <div class="stats-grid">
                  <div class="stat-item">
                    <div class="stat-value">{{ exerciseHistory.length }}</div>
                    <div class="stat-label">Practices</div>
                  </div>
                  <div class="stat-item">
                    <div class="stat-value">{{ averageBpm }}</div>
                    <div class="stat-label">Avg. BPM</div>
                  </div>
                  <div class="stat-item">
                    <div class="stat-value">{{ maxBpm }}</div>
                    <div class="stat-label">Max BPM</div>
                  </div>
                  <div class="stat-item">
                    <div class="stat-value">{{ averageRating.toFixed(1) }}</div>
                    <div class="stat-label">Avg. Rating</div>
                  </div>
                </div>
                
                <!-- Progress chart -->
                <div v-if="exerciseHistory.length > 1" class="progress-chart mt-4">
                  <h6>BPM Progress</h6>
                  <div class="chart-container" style="position: relative; height: 200px;">
                    <canvas ref="progressChart"></canvas>
                  </div>
                </div>
                
                <!-- Last practiced info -->
                <div v-if="lastPracticed" class="last-practiced mt-4">
                  <h6>Last Practiced</h6>
                  <div class="last-practice-info">
                    <div class="info-row">
                      <span class="info-label">Date:</span>
                      <span class="info-value">{{ formatDate(lastPracticed.date) }}</span>
                    </div>
                    <div class="info-row">
                      <span class="info-label">BPM:</span>
                      <span class="info-value">{{ lastPracticed.bpm }}</span>
                    </div>
                    <div class="info-row">
                      <span class="info-label">Time Signature:</span>
                      <span class="info-value">{{ lastPracticed.time_signature }}</span>
                    </div>
                    <div class="info-row">
                      <span class="info-label">Rating:</span>
                      <span class="info-value">
                        <span class="stars">
                          <i 
                            v-for="i in 5" 
                            :key="i" 
                            :class="['bi', i <= lastPracticed.rating ? 'bi-star-fill' : 'bi-star']"
                          ></i>
                        </span>
                      </span>
                    </div>
                    <div v-if="lastPracticed.notes" class="info-row notes">
                      <div class="info-label">Notes:</div>
                      <div class="info-value">{{ lastPracticed.notes }}</div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="card-footer">
              <button class="btn btn-primary w-100" @click="practiceExercise">
                <i class="bi bi-play-circle me-1"></i> Practice Now
              </button>
            </div>
          </div>
          
          <!-- Recent history -->
          <div class="card">
            <div class="card-header d-flex justify-content-between align-items-center">
              <h5 class="mb-0">Recent History</h5>
              <router-link 
                :to="`/history?exercise_id=${exerciseId}`" 
                class="btn btn-sm btn-outline-primary"
              >
                View All
              </router-link>
            </div>
            <div class="card-body p-0">
              <div v-if="loadingHistory" class="text-center py-3">
                <div class="spinner-border spinner-border-sm text-primary" role="status">
                  <span class="visually-hidden">Loading...</span>
                </div>
                <p class="mt-2 mb-0">Loading history...</p>
              </div>
              <div v-else-if="exerciseHistory.length === 0" class="text-center py-3">
                <p class="mb-0">No practice history yet.</p>
              </div>
              <div v-else class="history-list">
                <div 
                  v-for="(entry, index) in exerciseHistory.slice(0, 5)" 
                  :key="entry.id" 
                  class="history-item"
                  :class="{ 'border-bottom': index < Math.min(exerciseHistory.length, 5) - 1 }"
                >
                  <div class="d-flex justify-content-between align-items-center">
                    <div class="history-date">{{ formatDate(entry.date) }}</div>
                    <div class="history-bpm">{{ entry.bpm }} BPM</div>
                  </div>
                  <div class="d-flex justify-content-between align-items-center mt-1">
                    <div class="history-time-sig">{{ entry.time_signature }}</div>
                    <div class="history-rating">
                      <small class="stars">
                        <i 
                          v-for="i in 5" 
                          :key="i" 
                          :class="['bi', i <= entry.rating ? 'bi-star-fill' : 'bi-star']"
                        ></i>
                      </small>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Practice Exercise Modal -->
    <div 
      class="modal fade" 
      id="practiceModal" 
      tabindex="-1" 
      aria-labelledby="practiceModalLabel" 
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="practiceModalLabel">
              Practice {{ exercise?.name }}
            </h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="recordPractice">
              <div class="mb-3">
                <label for="practiceBpm" class="form-label">BPM</label>
                <input 
                  type="number" 
                  class="form-control" 
                  id="practiceBpm" 
                  v-model="practiceForm.bpm" 
                  min="1"
                  required
                >
              </div>
              <div class="mb-3">
                <label for="practiceTimeSignature" class="form-label">Time Signature</label>
                <select 
                  class="form-select" 
                  id="practiceTimeSignature" 
                  v-model="practiceForm.time_signature"
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
                <label for="practiceRating" class="form-label">Rating (1-5)</label>
                <div class="rating-input">
                  <div class="stars">
                    <i 
                      v-for="rating in 5" 
                      :key="rating"
                      :class="['bi', rating <= practiceForm.rating ? 'bi-star-fill' : 'bi-star']"
                      @click="practiceForm.rating = rating"
                    ></i>
                  </div>
                  <span class="rating-text">{{ getRatingDescription(practiceForm.rating) }}</span>
                </div>
              </div>
              <div class="mb-3">
                <label for="practiceNotes" class="form-label">Notes</label>
                <textarea 
                  class="form-control" 
                  id="practiceNotes" 
                  v-model="practiceForm.notes" 
                  rows="3"
                ></textarea>
              </div>
              <div class="text-end">
                <button type="button" class="btn btn-secondary me-2" data-bs-dismiss="modal">
                  Cancel
                </button>
                <button type="submit" class="btn btn-primary" :disabled="historyStore.loading">
                  <span v-if="historyStore.loading" class="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
                  Save Practice
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Add Link Modal -->
    <div 
      class="modal fade" 
      id="addLinkModal" 
      tabindex="-1" 
      aria-labelledby="addLinkModalLabel" 
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="addLinkModalLabel">Add Reference Link</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="addLink">
              <div class="mb-3">
                <label for="linkUrl" class="form-label">URL</label>
                <input 
                  type="url" 
                  class="form-control" 
                  id="linkUrl" 
                  v-model="linkForm.url" 
                  placeholder="https://example.com"
                  required
                >
              </div>
              <div class="mb-3">
                <label for="linkDescription" class="form-label">Description</label>
                <input 
                  type="text" 
                  class="form-control" 
                  id="linkDescription" 
                  v-model="linkForm.description" 
                  placeholder="Describe this link (optional)"
                >
              </div>
              <div class="text-end">
                <button type="button" class="btn btn-secondary me-2" data-bs-dismiss="modal">
                  Cancel
                </button>
                <button type="submit" class="btn btn-primary" :disabled="exercisesStore.loading">
                  <span v-if="exercisesStore.loading" class="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
                  Add Link
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Delete Link Confirmation Modal -->
    <div 
      class="modal fade" 
      id="deleteLinkModal" 
      tabindex="-1" 
      aria-labelledby="deleteLinkModalLabel" 
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="deleteLinkModalLabel">Confirm Delete</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <p>Are you sure you want to delete this link?</p>
            <p><strong>{{ linkToDelete?.description || linkToDelete?.url }}</strong></p>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
            <button type="button" class="btn btn-danger" @click="deleteLink" :disabled="exercisesStore.loading">
              <span v-if="exercisesStore.loading" class="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import { useExercisesStore } from '@/store/modules/exercises';
import { useHistoryStore } from '@/store/modules/history';
import { formatDate } from '@/utils/dateUtils';
import { getRatingDescription } from '@/utils/formatters';
import Chart from 'chart.js/auto';
import { Modal } from 'bootstrap'

export default {
  name: 'ExerciseDetails',
  
  props: {
    exerciseId: {
      type: Number,
      required: true
    }
  },
  
  setup(props) {
    const router = useRouter();
    const exercisesStore = useExercisesStore();
    const historyStore = useHistoryStore();
    
    const progressChart = ref(null);
    let chart = null;
    
    // State for practice history
    const loadingHistory = ref(false);
    
    // Form state
    const practiceForm = ref({
      exercise_id: props.exerciseId,
      date: new Date().toISOString(),
      bpm: 100,
      time_signature: '4/4',
      notes: '',
      rating: 3
    });
    
    const linkForm = ref({
      exercise_id: props.exerciseId,
      url: '',
      description: ''
    });
    
    const linkToDelete = ref(null);
    
    // Modal references
    let practiceModal = null;
    let addLinkModal = null;
    let deleteLinkModal = null;
    
    // Load data on component mount
    onMounted(async () => {
      // Initialize modals
      const practiceModalEl = document.getElementById('practiceModal')
      if (practiceModalEl) {
          practiceModal = new Modal(practiceModalEl)
      }
      const addLinkModalEl = document.getElementById('addLinkModal')
      if (addLinkModalEl) {
          addLinkModal = new Modal(addLinkModalEl)
      }
      const deleteLinkEl = document.getElementById('deleteLinkModal')
      if (deleteLinkEl) {
          deleteLinkModal = new Modal(deleteLinkEl)
      }
      
      // Fetch exercise details
      await exercisesStore.fetchExerciseById(props.exerciseId);
      
      // Load exercise history
      await loadExerciseHistory();
      
      // Initialize chart if we have history data
      if (exerciseHistory.value.length > 1) {
        initProgressChart();
      }
    });
    
    // Watch for exercise ID changes
    watch(() => props.exerciseId, async (newId) => {
      if (newId) {
        await exercisesStore.fetchExerciseById(newId);
        await loadExerciseHistory();
        
        // Update form exercise IDs
        practiceForm.value.exercise_id = newId;
        linkForm.value.exercise_id = newId;
        
        // Reinitialize chart if needed
        if (exerciseHistory.value.length > 1) {
          if (chart) {
            chart.destroy();
          }
          initProgressChart();
        }
      }
    });
    
    // Access the current exercise
    const exercise = computed(() => {
      return exercisesStore.currentExercise;
    });
    
    // Load exercise history
    const loadExerciseHistory = async () => {
      loadingHistory.value = true;
      historyStore.setExerciseId(props.exerciseId);
      
      try {
        await historyStore.fetchExerciseHistory({
          page_size: 100 // Get a large number to calculate stats
        });
      } catch (error) {
        console.error('Error loading exercise history:', error);
      } finally {
        loadingHistory.value = false;
      }
    };
    
    // Access the exercise history
    const exerciseHistory = computed(() => {
      return historyStore.exerciseHistory.sort((a, b) => 
        new Date(b.date) - new Date(a.date)
      );
    });
    
    // Calculate stats
    const averageBpm = computed(() => {
      if (exerciseHistory.value.length === 0) return 0;
      
      const sum = exerciseHistory.value.reduce((total, entry) => total + entry.bpm, 0);
      return Math.round(sum / exerciseHistory.value.length);
    });
    
    const maxBpm = computed(() => {
      if (exerciseHistory.value.length === 0) return 0;
      
      return Math.max(...exerciseHistory.value.map(entry => entry.bpm));
    });
    
    const averageRating = computed(() => {
      if (exerciseHistory.value.length === 0) return 0;
      
      const sum = exerciseHistory.value.reduce((total, entry) => total + entry.rating, 0);
      return sum / exerciseHistory.value.length;
    });
    
    const lastPracticed = computed(() => {
      if (exerciseHistory.value.length === 0) return null;
      
      return exerciseHistory.value[0];
    });
    
    // Initialize chart with BPM progress data
    const initProgressChart = () => {
      if (!progressChart.value) return;
      
      // Process data for chart - get last 10 entries in chronological order
      const entries = [...exerciseHistory.value]
        .sort((a, b) => new Date(a.date) - new Date(b.date))
        .slice(-10);
      
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
          },
          plugins: {
            legend: {
              display: false
            }
          }
        }
      });
    };
    
    // Action handlers
    const editExercise = () => {
      router.push(`/exercises/${props.exerciseId}/edit`);
    };
    
    const practiceExercise = () => {
      // Set default values based on last practice if available
      if (lastPracticed.value) {
        practiceForm.value = {
          exercise_id: props.exerciseId,
          date: new Date().toISOString(),
          bpm: lastPracticed.value.bpm,
          time_signature: lastPracticed.value.time_signature,
          notes: '',
          rating: 3
        };
      } else {
        practiceForm.value = {
          exercise_id: props.exerciseId,
          date: new Date().toISOString(),
          bpm: 100,
          time_signature: '4/4',
          notes: '',
          rating: 3
        };
      }
      
      practiceModal.show();
    };
    
    const recordPractice = async () => {
      try {
        await historyStore.createExerciseHistory(practiceForm.value);
        
        // Close modal
        practiceModal.hide();
        
        // Reload history
        await loadExerciseHistory();
        
        // Reinitialize chart if needed
        if (exerciseHistory.value.length > 1) {
          if (chart) {
            chart.destroy();
          }
          initProgressChart();
        }
      } catch (error) {
        console.error('Error recording practice:', error);
      }
    };
    
    // Link management
    const openAddLinkModal = () => {
      linkForm.value = {
        exercise_id: props.exerciseId,
        url: '',
        description: ''
      };
      
      addLinkModal.show();
    };
    
    const addLink = async () => {
      try {
        await exercisesStore.addExerciseLink(linkForm.value);
        
        // Close modal and reset form
        addLinkModal.hide();
        linkForm.value = {
          exercise_id: props.exerciseId,
          url: '',
          description: ''
        };
      } catch (error) {
        console.error('Error adding link:', error);
      }
    };
    
    const confirmDeleteLink = (link) => {
      linkToDelete.value = link;
      deleteLinkModal.show();
    };
    
    const deleteLink = async () => {
      if (!linkToDelete.value) return;
      
      try {
        await exercisesStore.deleteExerciseLink(linkToDelete.value.id);
        
        // Close modal
        deleteLinkModal.hide();
        linkToDelete.value = null;
      } catch (error) {
        console.error('Error deleting link:', error);
      }
    };
    
    return {
      exercisesStore,
      historyStore,
      exercise,
      exerciseHistory,
      loadingHistory,
      practiceForm,
      linkForm,
      linkToDelete,
      progressChart,
      averageBpm,
      maxBpm,
      averageRating,
      lastPracticed,
      editExercise,
      practiceExercise,
      recordPractice,
      openAddLinkModal,
      addLink,
      confirmDeleteLink,
      deleteLink,
      formatDate,
      getRatingDescription
    };
  }
};
</script>
