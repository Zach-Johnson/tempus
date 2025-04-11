<template>
  <div class="exercise-filters card mb-4">
    <div class="card-body">
      <div class="d-flex justify-content-between align-items-center mb-3">
        <h5 class="card-title mb-0">Filters</h5>
        <button 
          class="btn btn-sm btn-outline-secondary" 
          @click="clearAll"
          :disabled="!hasFilters"
        >
          <i class="bi bi-x-circle me-1"></i> Clear All
        </button>
      </div>
      
      <div class="row g-3">
        <!-- Search input -->
        <div class="col-md-4">
          <div class="form-group">
            <label for="searchTerm" class="form-label">Search</label>
            <div class="input-group">
              <input 
                type="text" 
                class="form-control" 
                id="searchTerm" 
                v-model="filters.searchTerm"
                placeholder="Search exercises..."
              >
              <button 
                v-if="filters.searchTerm" 
                class="btn btn-outline-secondary" 
                type="button"
                @click="filters.searchTerm = ''"
              >
                <i class="bi bi-x"></i>
              </button>
            </div>
          </div>
        </div>
        
        <!-- Category filter -->
        <div class="col-md-4">
          <div class="form-group">
            <label for="categoryFilter" class="form-label">Categories</label>
            <select 
              multiple 
              class="form-select" 
              id="categoryFilter" 
              v-model="filters.categoryIds"
              size="1"
            >
              <option v-for="category in categoriesStore.categories" :key="category.id" :value="category.id">
                {{ category.name }}
              </option>
            </select>
          </div>
        </div>
        
        <!-- Tag filter -->
        <div class="col-md-4">
          <div class="form-group">
            <label for="tagFilter" class="form-label">Tags</label>
            <select 
              multiple 
              class="form-select" 
              id="tagFilter" 
              v-model="filters.tagIds"
              size="1"
            >
              <option v-for="tag in tagsStore.tags" :key="tag.id" :value="tag.id">
                {{ tag.name }}
              </option>
            </select>
          </div>
        </div>
      </div>
      
      <!-- Selected filters -->
      <div v-if="hasFilters" class="selected-filters mt-3">
        <div class="selected-filters-label mb-2">Applied filters:</div>
        <div class="filter-badges">
          <!-- Search term badge -->
          <span v-if="filters.searchTerm" class="badge bg-light text-dark me-2 mb-2">
            Search: {{ filters.searchTerm }}
            <button class="btn-close btn-close-sm ms-1" @click="filters.searchTerm = ''" aria-label="Remove search filter"></button>
          </span>
          
          <!-- Category badges -->
          <span 
            v-for="categoryId in filters.categoryIds" 
            :key="`cat-${categoryId}`" 
            class="badge bg-primary me-2 mb-2"
          >
            {{ getCategoryName(categoryId) }}
            <button 
              class="btn-close btn-close-white btn-close-sm ms-1" 
              @click="removeCategory(categoryId)" 
              aria-label="Remove category filter"
            ></button>
          </span>
          
          <!-- Tag badges -->
          <span 
            v-for="tagId in filters.tagIds" 
            :key="`tag-${tagId}`" 
            class="badge bg-secondary me-2 mb-2"
          >
            {{ getTagName(tagId) }}
            <button 
              class="btn-close btn-close-white btn-close-sm ms-1" 
              @click="removeTag(tagId)" 
              aria-label="Remove tag filter"
            ></button>
          </span>
        </div>
      </div>
      
      <!-- Apply filters button -->
      <div class="d-flex justify-content-end mt-3">
        <button 
          class="btn btn-primary" 
          @click="applyFilters"
          :disabled="!filtersChanged"
        >
          <i class="bi bi-filter me-1"></i> Apply Filters
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue';
import { useCategoriesStore } from '@/store/modules/categories';
import { useTagsStore } from '@/store/modules/tags';

export default {
  name: 'ExerciseFilters',
  
  props: {
    initialFilters: {
      type: Object,
      default: () => ({
        categoryIds: [],
        tagIds: [],
        searchTerm: ''
      })
    }
  },
  
  emits: ['apply-filters', 'clear-filters'],
  
  setup(props, { emit }) {
    const categoriesStore = useCategoriesStore();
    const tagsStore = useTagsStore();
    
    // Local filters state
    const filters = ref({
      categoryIds: [...(props.initialFilters.categoryIds || [])],
      tagIds: [...(props.initialFilters.tagIds || [])],
      searchTerm: props.initialFilters.searchTerm || ''
    });
    
    // Track original filters to detect changes
    const originalFilters = ref({
      categoryIds: [...(props.initialFilters.categoryIds || [])],
      tagIds: [...(props.initialFilters.tagIds || [])],
      searchTerm: props.initialFilters.searchTerm || ''
    });
    
    // Watch for external filter changes
    watch(() => props.initialFilters, (newFilters) => {
      filters.value = {
        categoryIds: [...(newFilters.categoryIds || [])],
        tagIds: [...(newFilters.tagIds || [])],
        searchTerm: newFilters.searchTerm || ''
      };
      
      originalFilters.value = {
        categoryIds: [...(newFilters.categoryIds || [])],
        tagIds: [...(newFilters.tagIds || [])],
        searchTerm: newFilters.searchTerm || ''
      };
    }, { deep: true });
    
    // Load categories and tags on mount
    onMounted(() => {
      if (categoriesStore.categories.length === 0) {
        categoriesStore.fetchCategories({ page_size: 100 });
      }
      
      if (tagsStore.tags.length === 0) {
        tagsStore.fetchTags({ page_size: 100 });
      }
    });
    
    // Check if any filters are applied
    const hasFilters = computed(() => {
      return (
        filters.value.categoryIds.length > 0 ||
        filters.value.tagIds.length > 0 ||
        filters.value.searchTerm.trim() !== ''
      );
    });
    
    // Check if filters have changed from original
    const filtersChanged = computed(() => {
      return (
        JSON.stringify(filters.value.categoryIds) !== JSON.stringify(originalFilters.value.categoryIds) ||
        JSON.stringify(filters.value.tagIds) !== JSON.stringify(originalFilters.value.tagIds) ||
        filters.value.searchTerm !== originalFilters.value.searchTerm
      );
    });
    
    // Helper methods for display
    const getCategoryName = (id) => {
      const category = categoriesStore.categories.find(c => c.id === id);
      return category ? category.name : 'Unknown Category';
    };
    
    const getTagName = (id) => {
      const tag = tagsStore.tags.find(t => t.id === id);
      return tag ? tag.name : 'Unknown Tag';
    };
    
    // Remove individual filters
    const removeCategory = (id) => {
      filters.value.categoryIds = filters.value.categoryIds.filter(
        categoryId => categoryId !== id
      );
    };
    
    const removeTag = (id) => {
      filters.value.tagIds = filters.value.tagIds.filter(
        tagId => tagId !== id
      );
    };
    
    // Apply all filters
    const applyFilters = () => {
      emit('apply-filters', { ...filters.value });
      
      // Update original filters
      originalFilters.value = {
        categoryIds: [...filters.value.categoryIds],
        tagIds: [...filters.value.tagIds],
        searchTerm: filters.value.searchTerm
      };
    };
    
    // Clear all filters
    const clearAll = () => {
      filters.value = {
        categoryIds: [],
        tagIds: [],
        searchTerm: ''
      };
      
      emit('clear-filters');
      
      // Update original filters
      originalFilters.value = {
        categoryIds: [],
        tagIds: [],
        searchTerm: ''
      };
    };
    
    return {
      categoriesStore,
      tagsStore,
      filters,
      hasFilters,
      filtersChanged,
      getCategoryName,
      getTagName,
      removeCategory,
      removeTag,
      applyFilters,
      clearAll
    };
  }
};
</script>

<style scoped>
.exercise-filters {
  background-color: #f8f9fa;
  border: none;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.card-title {
  font-weight: 600;
}

.selected-filters-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #6c757d;
}

.filter-badges {
  display: flex;
  flex-wrap: wrap;
}

.badge {
  display: inline-flex;
  align-items: center;
  padding: 0.5rem 0.75rem;
  font-weight: 500;
}

.btn-close-sm {
  font-size: 0.625rem;
  padding: 0.125rem;
}
</style>
