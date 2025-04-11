<template>
  <div class="tags-view">
    <div class="container-fluid">
      <div class="page-header mb-4">
        <h1>Tags</h1>
        <p class="text-muted">Create and manage tags to help organize your exercises.</p>
      </div>
      
      <div class="card">
        <div class="card-header d-flex justify-content-between align-items-center">
          <h5 class="mb-0">Manage Tags</h5>
          <button class="btn btn-primary" @click="openCreateModal">
            <i class="bi bi-plus-circle me-1"></i> New Tag
          </button>
        </div>
        <div class="card-body">
          <!-- Loading indicator -->
          <div v-if="tagsStore.loading" class="text-center my-5">
            <div class="spinner-border text-primary" role="status">
              <span class="visually-hidden">Loading...</span>
            </div>
            <p class="mt-2">Loading tags...</p>
          </div>
          
          <!-- Error message -->
          <div v-else-if="tagsStore.error" class="alert alert-danger" role="alert">
            <i class="bi bi-exclamation-triangle-fill me-2"></i>
            {{ tagsStore.error }}
          </div>
          
          <!-- Empty state -->
          <div 
            v-else-if="tagsStore.tags.length === 0" 
            class="empty-state text-center my-5"
          >
            <i class="bi bi-tag display-1 text-muted"></i>
            <h3 class="mt-3">No Tags Found</h3>
            <p class="text-muted">Tags help you organize and filter your exercises.</p>
            <button class="btn btn-primary mt-3" @click="openCreateModal">
              <i class="bi bi-plus-circle me-1"></i> Create First Tag
            </button>
          </div>
          
          <!-- Tag list -->
          <div v-else>
            <div class="tag-list">
              <div 
                v-for="tag in tagsStore.tags" 
                :key="tag.id" 
                class="tag-item"
              >
                <div class="tag-content">
                  <span class="tag-name">{{ tag.name }}</span>
                  <div class="tag-actions">
                    <button 
                      class="btn btn-sm btn-outline-danger" 
                      @click="confirmDelete(tag)"
                      title="Delete Tag"
                    >
                      <i class="bi bi-trash"></i>
                    </button>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- Pagination -->
            <Pagination
              v-if="tagsStore.totalCount > 0"
              :total-items="tagsStore.totalCount"
              :current-page="currentPage"
              :page-size="pageSize"
              @page-change="handlePageChange"
              @page-size-change="handlePageSizeChange"
              class="mt-4"
            />
          </div>
        </div>
      </div>
    </div>
    
    <!-- Create Tag Modal -->
    <div 
      class="modal fade" 
      id="createTagModal" 
      tabindex="-1" 
      aria-labelledby="createTagModalLabel" 
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="createTagModalLabel">Create New Tag</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="createTag">
              <div class="mb-3">
                <label for="tagName" class="form-label">Tag Name</label>
                <input 
                  type="text" 
                  class="form-control" 
                  id="tagName" 
                  v-model="tagForm.name" 
                  required
                  placeholder="Enter tag name"
                >
              </div>
              <div class="text-end">
                <button type="button" class="btn btn-secondary me-2" data-bs-dismiss="modal">
                  Cancel
                </button>
                <button type="submit" class="btn btn-primary" :disabled="tagsStore.loading">
                  <span v-if="tagsStore.loading" class="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
                  Create
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
      id="deleteTagModal" 
      tabindex="-1" 
      aria-labelledby="deleteTagModalLabel" 
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="deleteTagModalLabel">Confirm Delete</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <p>Are you sure you want to delete the tag <strong>{{ tagToDelete?.name }}</strong>?</p>
            <p class="text-danger">This action cannot be undone. This will remove the tag from all exercises.</p>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
            <button type="button" class="btn btn-danger" @click="deleteTag" :disabled="tagsStore.loading">
              <span v-if="tagsStore.loading" class="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import Pagination from '@/components/common/Pagination.vue';
import { useTagsStore } from '@/store/modules/tags';
import { Modal } from 'bootstrap'

export default {
  name: 'TagsView',
  
  components: {
    Pagination
  },
  
  setup() {
    const tagsStore = useTagsStore();
    
    // Pagination state
    const currentPage = ref(1);
    const pageSize = ref(20);
    
    // Form state
    const tagForm = ref({
      name: ''
    });
    
    // Delete confirmation state
    const tagToDelete = ref(null);
    
    // Modal references
    let createModal = null;
    let deleteModal = null;
    
    onMounted(() => {
      // Initialize modals
      const createModalEl = document.getElementById('createTagModal')
      if (createModalEl) {
          createModal = new Modal(createModalEl)
      }
      const deleteModalEl = document.getElementById('deleteTagModal')
      if (deleteModalEl) {
          deleteModal = new Modal(deleteModalEl)
      }
      
      // Fetch tags
      loadTags();
    });
    
    // Load tags with pagination
    const loadTags = () => {
      tagsStore.fetchTags({
        page: currentPage.value,
        page_size: pageSize.value
      });
    };
    
    // Pagination handlers
    const handlePageChange = (page) => {
      currentPage.value = page;
      loadTags();
    };
    
    const handlePageSizeChange = (size) => {
      pageSize.value = size;
      currentPage.value = 1; // Reset to first page
      loadTags();
    };
    
    // Modal handlers
    const openCreateModal = () => {
      tagForm.value = {
        name: ''
      };
      createModal.show();
    };
    
    const confirmDelete = (tag) => {
      tagToDelete.value = tag;
      deleteModal.show();
    };
    
    // Form submission
    const createTag = async () => {
      try {
        await tagsStore.createTag(tagForm.value);
        
        // Close modal and reset form
        createModal.hide();
        tagForm.value = {
          name: ''
        };
        
        // Reload tags
        loadTags();
      } catch (error) {
        console.error('Error creating tag:', error);
      }
    };
    
    // Delete tag
    const deleteTag = async () => {
      if (!tagToDelete.value) return;
      
      try {
        await tagsStore.deleteTag(tagToDelete.value.id);
        
        // Close modal
        deleteModal.hide();
        tagToDelete.value = null;
        
        // Reload tags
        loadTags();
      } catch (error) {
        console.error('Error deleting tag:', error);
      }
    };
    
    return {
      tagsStore,
      currentPage,
      pageSize,
      tagForm,
      tagToDelete,
      handlePageChange,
      handlePageSizeChange,
      openCreateModal,
      confirmDelete,
      createTag,
      deleteTag
    };
  }
};
</script>

<style scoped>
.empty-state {
  padding: 2rem;
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

.tag-item {
  border: 1px solid #dee2e6;
  border-radius: 0.25rem;
  background-color: #f8f9fa;
  padding: 0.5rem 0.75rem;
  transition: background-color 0.2s;
}

.tag-item:hover {
  background-color: #e9ecef;
}

.tag-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.tag-name {
  font-weight: 500;
}

.tag-actions {
  display: flex;
  gap: 0.5rem;
}
</style>
