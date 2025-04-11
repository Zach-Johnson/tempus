<script>
import { ref, onMounted } from 'vue';
import Pagination from '@/components/common/Pagination.vue';
import { useCategoriesStore } from '@/store/modules/categories';
import { formatDate } from '@/utils/dateUtils';
import { truncateText } from '@/utils/formatters';
import { Modal } from 'bootstrap'

export default {
  name: 'CategoryList',
  
  components: {
    Pagination
  },
  
  setup() {
    const categoriesStore = useCategoriesStore();
    
    // Pagination state
    const currentPage = ref(1);
    const pageSize = ref(10);
    
    // Form state
    const editMode = ref(false);
    const categoryForm = ref({
      id: null,
      name: '',
      description: ''
    });
    
    // Delete confirmation state
    const categoryToDelete = ref(null);
    
    // Modal references
    let categoryModal = null;
    let deleteModal = null;
    
    onMounted(() => {
      // Initialize modals
      const catModalEl = document.getElementById('categoryModal')
      if (catModalEl) {
          categoryModal = new Modal(catModalEl)
      }
      const deleteModalEl = document.getElementById('deleteModal')
      if (deleteModalEl) {
          deleteModal = new Modal(deleteModalEl)
      }
      
      // Fetch categories
      loadCategories();
    });
    
    // Load categories with pagination
    const loadCategories = () => {
      categoriesStore.fetchCategories({
        page: currentPage.value,
        page_size: pageSize.value
      });
    };
    
    // Pagination handlers
    const handlePageChange = (page) => {
      currentPage.value = page;
      loadCategories();
    };
    
    const handlePageSizeChange = (size) => {
      pageSize.value = size;
      currentPage.value = 1; // Reset to first page
      loadCategories();
    };
    
    // Modal handlers
    const openCreateModal = () => {
      editMode.value = false;
      categoryForm.value = {
        id: null,
        name: '',
        description: ''
      };
      categoryModal.show();
    };
    
    const editCategory = (category) => {
      editMode.value = true;
      categoryForm.value = { ...category };
      categoryModal.show();
    };
    
    const confirmDelete = (category) => {
      categoryToDelete.value = category;
      deleteModal.show();
    };
    
    // Form submission
    const saveCategory = async () => {
      try {
        if (editMode.value) {
          await categoriesStore.updateCategory(categoryForm.value);
        } else {
          await categoriesStore.createCategory(categoryForm.value);
        }
        
        // Close modal and reset form
        categoryModal.hide();
        categoryForm.value = {
          id: null,
          name: '',
          description: ''
        };
        
        // Reload categories
        loadCategories();
      } catch (error) {
        console.error('Error saving category:', error);
      }
    };
    
    // Delete category
    const deleteCategory = async () => {
      if (!categoryToDelete.value) return;
      
      try {
        await categoriesStore.deleteCategory(categoryToDelete.value.id);
        
        // Close modal
        deleteModal.hide();
        categoryToDelete.value = null;
        
        // Reload categories
        loadCategories();
      } catch (error) {
        console.error('Error deleting category:', error);
      }
    };
    
    return {
      categoriesStore,
      currentPage,
      pageSize,
      editMode,
      categoryForm,
      categoryToDelete,
      handlePageChange,
      handlePageSizeChange,
      openCreateModal,
      editCategory,
      confirmDelete,
      saveCategory,
      deleteCategory,
      formatDate,
      truncateText
    };
  }
};
</script>

<style scoped>
.empty-state {
  padding: 2rem;
}

.card {
  transition: transform 0.2s, box-shadow 0.2s;
}

.card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.card-title {
  font-weight: 600;
}

.list-header {
  margin-bottom: 1.5rem;
}
</style>
<template>
  <div class="category-list">
    <!-- Category list header with action buttons -->
    <div class="list-header d-flex justify-content-between align-items-center mb-4">
      <h2>Categories</h2>
      <button class="btn btn-primary" @click="openCreateModal">
        <i class="bi bi-plus-circle me-1"></i> New Category
      </button>
    </div>
    
    <!-- Loading indicator -->
    <div v-if="categoriesStore.loading" class="text-center my-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
      <p class="mt-2">Loading categories...</p>
    </div>
    
    <!-- Error message -->
    <div v-else-if="categoriesStore.error" class="alert alert-danger" role="alert">
      <i class="bi bi-exclamation-triangle-fill me-2"></i>
      {{ categoriesStore.error }}
    </div>
    
    <!-- Empty state -->
    <div 
      v-else-if="categoriesStore.categories.length === 0" 
      class="empty-state text-center my-5"
    >
      <i class="bi bi-folder-x display-1 text-muted"></i>
      <h3 class="mt-3">No Categories Found</h3>
      <p class="text-muted">Start by creating a new category to organize your exercises.</p>
      <button class="btn btn-primary mt-3" @click="openCreateModal">
        <i class="bi bi-plus-circle me-1"></i> Create First Category
      </button>
    </div>
    
    <!-- Category list -->
    <div v-else class="row row-cols-1 row-cols-md-2 row-cols-lg-3 g-4">
      <div 
        v-for="category in categoriesStore.categories" 
        :key="category.id" 
        class="col"
      >
        <div class="card h-100">
          <div class="card-body">
            <h5 class="card-title">{{ category.name }}</h5>
            <p class="card-text">{{ truncateText(category.description, 120) }}</p>
          </div>
          <div class="card-footer">
            <div class="d-flex justify-content-between align-items-center">
              <small class="text-muted">
                Created: {{ formatDate(category.created_at) }}
              </small>
              <div class="btn-group">
                <button 
                  class="btn btn-sm btn-outline-primary" 
                  @click="editCategory(category)"
                  title="Edit Category"
                >
                  <i class="bi bi-pencil"></i>
                </button>
                <button 
                  class="btn btn-sm btn-outline-danger" 
                  @click="confirmDelete(category)"
                  title="Delete Category"
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
      v-if="categoriesStore.totalCount > 0"
      :total-items="categoriesStore.totalCount"
      :current-page="currentPage"
      :page-size="pageSize"
      @page-change="handlePageChange"
      @page-size-change="handlePageSizeChange"
      class="mt-4"
    />
    
    <!-- Category Form Modal -->
    <div 
      class="modal fade" 
      id="categoryModal" 
      tabindex="-1" 
      aria-labelledby="categoryModalLabel" 
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="categoryModalLabel">
              {{ editMode ? 'Edit Category' : 'Create New Category' }}
            </h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveCategory">
              <div class="mb-3">
                <label for="categoryName" class="form-label">Category Name</label>
                <input 
                  type="text" 
                  class="form-control" 
                  id="categoryName" 
                  v-model="categoryForm.name" 
                  required
                  placeholder="Enter category name"
                >
              </div>
              <div class="mb-3">
                <label for="categoryDescription" class="form-label">Description</label>
                <textarea 
                  class="form-control" 
                  id="categoryDescription" 
                  v-model="categoryForm.description" 
                  rows="3"
                  placeholder="Enter category description"
                ></textarea>
              </div>
              <div class="text-end">
                <button type="button" class="btn btn-secondary me-2" data-bs-dismiss="modal">
                  Cancel
                </button>
                <button type="submit" class="btn btn-primary" :disabled="categoriesStore.loading">
                  <span v-if="categoriesStore.loading" class="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
                  {{ editMode ? 'Update' : 'Create' }}
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
            <p>Are you sure you want to delete the category <strong>{{ categoryToDelete?.name }}</strong>?</p>
            <p class="text-danger">This action cannot be undone. Exercises will be removed from this category but will not be deleted.</p>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
            <button type="button" class="btn btn-danger" @click="deleteCategory" :disabled="categoriesStore.loading">
              <span v-if="categoriesStore.loading" class="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
