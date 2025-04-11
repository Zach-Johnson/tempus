<template>
  <div class="exercise-list">
    <!-- Exercises list header with action buttons -->
    <div class="list-header d-flex justify-content-between align-items-center mb-4">
      <h2>Exercises</h2>
      <button class="btn btn-primary" @click="openCreateModal">
        <i class="bi bi-plus-circle me-1"></i> New Exercise
      </button>
    </div>
    
    <!-- Filters -->
    <ExerciseFilters 
      @apply-filters="applyFilters" 
      @clear-filters="clearFilters"
      :initial-filters="exercisesStore.filters"
    />
    
    <!-- Loading indicator -->
    <div v-if="exercisesStore.loading" class="text-center my-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
      <p class="mt-2">Loading exercises...</p>
    </div>
    
    <!-- Error message -->
    <div v-else-if="exercisesStore.error" class="alert alert-danger" role="alert">
      <i class="bi bi-exclamation-triangle-fill me-2"></i>
      {{ exercisesStore.error }}
    </div>
    
    <!-- Empty state -->
    <div 
      v-else-if="exercisesStore.exercises.length === 0" 
      class="empty-state text-center my-5"
    >
      <i class="bi bi-music-note-list display-1 text-muted"></i>
      <h3 class="mt-3">No Exercises Found</h3>
      <p v-if="hasActiveFilters" class="text-muted">
        No exercises match your current filters. Try adjusting your search criteria.
        <button class="btn btn-link p-0" @click="clearFilters">Clear filters</button>
      </p>
      <p v-else class="text-muted">
        Start by creating a new exercise to add to your practice routine.
      </p>
      <button v-if="!hasActiveFilters" class="btn btn-primary mt-3" @click="openCreateModal">
        <i class="bi bi-plus-circle me-1"></i> Create First Exercise
      </button>
    </div>
    
    <!-- Exercise grid -->
    <div v-else class="row row-cols-1 row-cols-md-2 row-cols-lg-3 g-4">
      <div 
        v-for="exercise in exercisesStore.exercises" 
        :key="exercise.id" 
        class="col"
      >
        <div class="card h-100 exercise-card">
          <!-- Exercise image if available -->
          <div v-if="exercise.music_image_path" class="exercise-image">
            <img :src="exercise.music_image_path" class="card-img-top" alt="Exercise sheet music">
          </div>
          <div v-else class="exercise-image-placeholder">
            <i class="bi bi-music-note-beamed"></i>
          </div>
          
          <div class="card-body">
            <h5 class="card-title">{{ exercise.name }}</h5>
            <p class="card-text">{{ truncateText(exercise.description, 100) }}</p>
            
            <!-- Categories and tags -->
            <div class="mb-2">
              <span 
                v-for="category in exercise.categories" 
                :key="category.id" 
                class="badge bg-primary me-1 mb-1"
              >
                {{ category.name }}
              </span>
              <span 
                v-for="tag in exercise.tags" 
                :key="tag.id" 
                class="badge bg-secondary me-1 mb-1"
              >
                {{ tag.name }}
              </span>
            </div>
          </div>
          
          <div class="card-footer">
            <div class="d-flex justify-content-between align-items-center">
              <small class="text-muted">
                Created: {{ formatDate(exercise.created_at) }}
              </small>
              <div class="btn-group">
                <router-link 
                  :to="`/exercises/${exercise.id}`"
                  class="btn btn-sm btn-outline-primary"
                  title="View Exercise"
                >
                  <i class="bi bi-eye"></i>
                </router-link>
                <button 
                  class="btn btn-sm btn-outline-success" 
                  @click="practiceExercise(exercise)"
                  title="Practice Now"
                >
                  <i class="bi bi-play-circle"></i>
                </button>
                <button 
                  class="btn btn-sm btn-outline-danger" 
                  @click="confirmDelete(exercise)"
                  title="Delete Exercise"
                >
                  <i class="bi bi-trash"></i>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Pagination -->
    <Pagination
      v-if="exercisesStore.totalCount > 0"
      :total-items="exercisesStore.totalCount"
      :current-page="currentPage"
      :page-size="pageSize"
      @page-change="handlePageChange"
      @page-size-change="handlePageSizeChange"
      class="mt-4"
    />
    
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
              Practice {{ exerciseToPractice?.name }}
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
    
    <!-- Delete Confirmation Modal -->
    <div 
      class="modal fade" 
      id="deleteModal" 
      tabindex="-1" 
      aria-labelledby="deleteModalLabel" 
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="deleteModalLabel">Confirm Delete</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <p>Are you sure you want to delete the exercise <strong>{{ exerciseToDelete?.name }}</strong>?</p>
            <p class="text-danger">This action cannot be undone. All practice data associated with this exercise will also be deleted.</p>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
            <button type="button" class="btn btn-danger" @click="deleteExercise" :disabled="exercisesStore.loading">
              <span v-if="exercisesStore.loading" class="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Create Exercise Modal -->
    <div 
      class="modal fade" 
      id="createExerciseModal" 
      tabindex="-1" 
      aria-labelledby="createExerciseModalLabel" 
      aria-hidden="true"
    >
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="createExerciseModalLabel">Create New Exercise</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="createExercise">
              <div class="mb-3">
                <label for="exerciseName" class="form-label">Exercise Name</label>
                <input 
                  type="text" 
                  class="form-control" 
                  id="exerciseName" 
                  v-model="exerciseForm.name" 
                  required
                  placeholder="Enter exercise name"
                >
              </div>
              
              <div class="mb-3">
                <label for="exerciseDescription" class="form-label">Description</label>
                <textarea 
                  class="form-control" 
                  id="exerciseDescription" 
                  v-model="exerciseForm.description" 
                  rows="3"
                  placeholder="Enter exercise description"
                ></textarea>
              </div>
              
              <div class="row">
                <div class="col-md-6 mb-3">
                  <label for="exerciseCategories" class="form-label">Categories</label>
                  <select 
                    multiple 
                    class="form-select" 
                    id="exerciseCategories" 
                    v-model="selectedCategories"
                    size="4"
                  >
                    <option 
                      v-for="category in categoriesStore.categories" 
                      :key="category.id" 
                      :value="category.id"
                    >
                      {{ category.name }}
                    </option>
                  </select>
                  <small class="form-text text-muted">Hold Ctrl/Cmd to select multiple</small>
                </div>
                
                <div class="col-md-6 mb-3">
                  <label for="exerciseTags" class="form-label">Tags</label>
                  <select 
                    multiple 
                    class="form-select" 
                    id="exerciseTags" 
                    v-model="selectedTags"
                    size="4"
                  >
                    <option 
                      v-for="tag in tagsStore.tags" 
                      :key="tag.id" 
                      :value="tag.id"
                    >
                      {{ tag.name }}
                    </option>
                  </select>
                  <small class="form-text text-muted">Hold Ctrl/Cmd to select multiple</small>
                </div>
              </div>
              
              <div class="text-end">
                <button type="button" class="btn btn-secondary me-2" data-bs-dismiss="modal">
                  Cancel
                </button>
                <button type="submit" class="btn btn-primary" :disabled="exercisesStore.loading">
                  <span v-if="exercisesStore.loading" class="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
                  Create Exercise
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import Pagination from '@/components/common/Pagination.vue';
import ExerciseFilters from '@/components/exercises/ExerciseFilters.vue';
import { useExercisesStore } from '@/store/modules/exercises';
import { useHistoryStore } from '@/store/modules/history';
import { useCategoriesStore } from '@/store/modules/categories';
import { useTagsStore } from '@/store/modules/tags';
import { formatDate } from '@/utils/dateUtils';
import { truncateText, getRatingDescription } from '@/utils/formatters';
import { Modal } from 'bootstrap'

export default {
  name: 'ExerciseList',
  
  components: {
    Pagination,
    ExerciseFilters
  },
  
  setup() {
    const router = useRouter();
    const exercisesStore = useExercisesStore();
    const historyStore = useHistoryStore();
    const categoriesStore = useCategoriesStore();
    const tagsStore = useTagsStore();
    
    // Pagination state
    const currentPage = ref(1);
    const pageSize = ref(12);
    
    // Modal state
    const exerciseToPractice = ref(null);
    const exerciseToDelete = ref(null);
    
    // Form state
    const practiceForm = ref({
      exercise_id: null,
      date: new Date().toISOString(),
      bpm: 100,
      time_signature: '4/4',
      notes: '',
      rating: 3
    });
    
    // Exercise form for creating new exercises
    const exerciseForm = ref({
      name: '',
      description: '',
      music_image_path: ''
    });
    
    // Selected categories and tags for new exercise
    const selectedCategories = ref([]);
    const selectedTags = ref([]);
    
    // Modal references
    let practiceModal = null;
    let deleteModal = null;
    let createExerciseModal = null;
    
    onMounted(() => {
      // Initialize modals
      const practiceModalEl = document.getElementById('practiceModal')
      if (practiceModalEl) {
          practiceModal = new Modal(practiceModalEl)
      }
      const deleteModalEl = document.getElementById('deleteModal')
      if (deleteModalEl) {
          deleteModal = new Modal(deleteModalEl)
      }
      const createExerciseModalEl = document.getElementById('createExerciseModal')
      if (createExerciseModalEl) {
          createExerciseModal = new Modal(createExerciseModalEl)
      }
      
      // Fetch exercises
      loadExercises();
      
      // Load categories and tags if needed
      if (categoriesStore.categories.length === 0) {
        categoriesStore.fetchCategories({ page_size: 100 });
      }
      
      if (tagsStore.tags.length === 0) {
        tagsStore.fetchTags({ page_size: 100 });
      }
    });
    
    // Check if filters are active
    const hasActiveFilters = computed(() => {
      const filters = exercisesStore.filters;
      return (
        (filters.categoryIds && filters.categoryIds.length > 0) ||
        (filters.tagIds && filters.tagIds.length > 0) ||
        (filters.searchTerm && filters.searchTerm.trim() !== '')
      );
    });
    
    // Load exercises with pagination and filters
    const loadExercises = () => {
      exercisesStore.fetchExercises({
        page: currentPage.value,
        page_size: pageSize.value
      });
    };
    
    // Pagination handlers
    const handlePageChange = (page) => {
      currentPage.value = page;
      loadExercises();
    };
    
    const handlePageSizeChange = (size) => {
      pageSize.value = size;
      currentPage.value = 1; // Reset to first page
      loadExercises();
    };
    
    // Filter handlers
    const applyFilters = (filters) => {
      exercisesStore.setFilters(filters);
      currentPage.value = 1; // Reset to first page
      loadExercises();
    };
    
    const clearFilters = () => {
      exercisesStore.clearFilters();
      loadExercises();
    };
    
    // Navigation handlers
    const openCreateModal = () => {
      // Reset the form
      exerciseForm.value = {
        name: '',
        description: '',
        music_image_path: ''
      };
      selectedCategories.value = [];
      selectedTags.value = [];
      
      // Show the modal
      createExerciseModal.show();
    };
    
    // Create a new exercise
    const createExercise = async () => {
      try {
        // Create the exercise first
        const newExercise = await exercisesStore.createExercise(exerciseForm.value);
        
        // Add categories if selected
        for (const categoryId of selectedCategories.value) {
          await exercisesStore.addExerciseToCategory(newExercise.id, categoryId);
        }
        
        // Add tags if selected
        for (const tagId of selectedTags.value) {
          await exercisesStore.addTagToExercise(newExercise.id, tagId);
        }
        
        // Hide the modal
        createExerciseModal.hide();
        
        // Reload exercises
        loadExercises();
        
        // Optionally navigate to the new exercise details
        router.push(`/exercises/${newExercise.id}`);
      } catch (error) {
        console.error('Failed to create exercise:', error);
      }
    };
    
    // Modal handlers
    const practiceExercise = (exercise) => {
      exerciseToPractice.value = exercise;
      practiceForm.value = {
        exercise_id: exercise.id,
        date: new Date().toISOString(),
        bpm: 100,
        time_signature: '4/4',
        notes: '',
        rating: 3
      };
      practiceModal.show();
    };
    
    const confirmDelete = (exercise) => {
      exerciseToDelete.value = exercise;
      deleteModal.show();
    };
    
    // Record practice
    const recordPractice = async () => {
      try {
        await historyStore.createExerciseHistory(practiceForm.value);
        
        // Close modal and reset form
        practiceModal.hide();
        exerciseToPractice.value = null;
        practiceForm.value = {
          exercise_id: null,
          date: new Date().toISOString(),
          bpm: 100,
          time_signature: '4/4',
          notes: '',
          rating: 3
        };
      } catch (error) {
        console.error('Error recording practice:', error);
      }
    };
    
    // Delete exercise
    const deleteExercise = async () => {
      if (!exerciseToDelete.value) return;
      
      try {
        await exercisesStore.deleteExercise(exerciseToDelete.value.id);
        
        // Close modal
        deleteModal.hide();
        exerciseToDelete.value = null;
        
        // Reload exercises
        loadExercises();
      } catch (error) {
        console.error('Error deleting exercise:', error);
      }
    };
    
    return {
      exercisesStore,
      historyStore,
      categoriesStore,
      tagsStore,
      currentPage,
      pageSize,
      exerciseToPractice,
      exerciseToDelete,
      practiceForm,
      exerciseForm,
      selectedCategories,
      selectedTags,
      hasActiveFilters,
      handlePageChange,
      handlePageSizeChange,
      applyFilters,
      clearFilters,
      openCreateModal,
      createExercise,
      practiceExercise,
      confirmDelete,
      recordPractice,
      deleteExercise,
      formatDate,
      truncateText,
      getRatingDescription
    };
  }
};
</script>

<style scoped>
.empty-state {
  padding: 2rem;
}

.exercise-card {
  transition: transform 0.2s, box-shadow 0.2s;
}

.exercise-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.exercise-image {
  height: 160px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f8f9fa;
}

.exercise-image img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.exercise-image-placeholder {
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f8f9fa;
  color: #adb5bd;
  font-size: 2rem;
}

.badge {
  font-weight: 500;
}

.rating-input {
  display: flex;
  align-items: center;
}

.stars {
  font-size: 1.5rem;
  color: #ffc107;
  cursor: pointer;
  margin-right: 1rem;
}

.rating-text {
  font-size: 0.875rem;
  color: #495057;
}
</style>
