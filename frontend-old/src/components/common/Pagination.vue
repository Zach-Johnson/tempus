<template>
  <div class="pagination-container">
    <nav v-if="totalPages > 1" aria-label="Pagination">
      <ul class="pagination">
        <!-- Previous button -->
        <li class="page-item" :class="{ disabled: currentPage === 1 }">
          <button 
            @click="changePage(currentPage - 1)" 
            class="page-link" 
            :disabled="currentPage === 1"
            aria-label="Previous"
          >
            <span aria-hidden="true">&laquo;</span>
          </button>
        </li>
        
        <!-- First page -->
        <li v-if="showFirstPage" class="page-item" :class="{ active: currentPage === 1 }">
          <button @click="changePage(1)" class="page-link">1</button>
        </li>
        
        <!-- Ellipsis before middle pages -->
        <li v-if="showStartEllipsis" class="page-item disabled">
          <span class="page-link">...</span>
        </li>
        
        <!-- Middle pages -->
        <li 
          v-for="page in middlePages" 
          :key="page" 
          class="page-item" 
          :class="{ active: currentPage === page }"
        >
          <button @click="changePage(page)" class="page-link">{{ page }}</button>
        </li>
        
        <!-- Ellipsis after middle pages -->
        <li v-if="showEndEllipsis" class="page-item disabled">
          <span class="page-link">...</span>
        </li>
        
        <!-- Last page -->
        <li v-if="showLastPage" class="page-item" :class="{ active: currentPage === totalPages }">
          <button @click="changePage(totalPages)" class="page-link">{{ totalPages }}</button>
        </li>
        
        <!-- Next button -->
        <li class="page-item" :class="{ disabled: currentPage === totalPages }">
          <button 
            @click="changePage(currentPage + 1)" 
            class="page-link" 
            :disabled="currentPage === totalPages"
            aria-label="Next"
          >
            <span aria-hidden="true">&raquo;</span>
          </button>
        </li>
      </ul>
    </nav>
    
    <!-- Page size selector -->
    <div v-if="showPageSize" class="page-size-selector">
      <label for="pageSize" class="form-label me-2">Items per page:</label>
      <select 
        id="pageSize" 
        class="form-select form-select-sm" 
        v-model="selectedPageSize"
        @change="changePageSize"
      >
        <option v-for="size in pageSizeOptions" :key="size" :value="size">{{ size }}</option>
      </select>
    </div>
  </div>
</template>

<script>
import { computed, ref, watch } from 'vue';

export default {
  name: 'Pagination',
  
  props: {
    totalItems: {
      type: Number,
      required: true
    },
    pageSize: {
      type: Number,
      default: 10
    },
    currentPage: {
      type: Number,
      default: 1
    },
    maxVisiblePages: {
      type: Number,
      default: 5
    },
    showPageSize: {
      type: Boolean,
      default: true
    },
    pageSizeOptions: {
      type: Array,
      default: () => [10, 25, 50, 100]
    }
  },
  
  emits: ['page-change', 'page-size-change'],
  
  setup(props, { emit }) {
    // Store selected page size
    const selectedPageSize = ref(props.pageSize);
    
    // Watch for external page size changes
    watch(() => props.pageSize, (newPageSize) => {
      selectedPageSize.value = newPageSize;
    });
    
    // Calculate total pages
    const totalPages = computed(() => {
      return Math.max(1, Math.ceil(props.totalItems / props.pageSize));
    });
    
    // Calculate which pages to show
    const middlePages = computed(() => {
      // If total pages is less than max visible pages, show all pages
      if (totalPages.value <= props.maxVisiblePages) {
        return Array.from({ length: totalPages.value }, (_, i) => i + 1);
      }
      
      // Calculate range of pages to show
      const halfMaxPages = Math.floor(props.maxVisiblePages / 2);
      let startPage = Math.max(props.currentPage - halfMaxPages, 2);
      let endPage = Math.min(props.currentPage + halfMaxPages, totalPages.value - 1);
      
      // Adjust if range is not full
      if (endPage - startPage + 1 < props.maxVisiblePages - 2) {
        if (props.currentPage < totalPages.value / 2) {
          endPage = Math.min(startPage + props.maxVisiblePages - 3, totalPages.value - 1);
        } else {
          startPage = Math.max(endPage - (props.maxVisiblePages - 3), 2);
        }
      }
      
      return Array.from({ length: endPage - startPage + 1 }, (_, i) => startPage + i);
    });
    
    // Show first page button?
    const showFirstPage = computed(() => {
      return totalPages.value > 1;
    });
    
    // Show last page button?
    const showLastPage = computed(() => {
      return totalPages.value > 1 && totalPages.value !== 1;
    });
    
    // Show ellipsis at start?
    const showStartEllipsis = computed(() => {
      return middlePages.value.length > 0 && middlePages.value[0] > 2;
    });
    
    // Show ellipsis at end?
    const showEndEllipsis = computed(() => {
      return middlePages.value.length > 0 && middlePages.value[middlePages.value.length - 1] < totalPages.value - 1;
    });
    
    // Change page event handler
    const changePage = (page) => {
      if (page < 1 || page > totalPages.value || page === props.currentPage) {
        return;
      }
      emit('page-change', page);
    };
    
    // Change page size event handler
    const changePageSize = () => {
      emit('page-size-change', selectedPageSize.value);
    };
    
    return {
      totalPages,
      middlePages,
      showFirstPage,
      showLastPage,
      showStartEllipsis,
      showEndEllipsis,
      changePage,
      selectedPageSize,
      changePageSize
    };
  }
};
</script>

<style scoped>
.pagination-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 1.5rem;
  margin-bottom: 1rem;
}

.pagination {
  margin-bottom: 0;
}

.page-size-selector {
  display: flex;
  align-items: center;
}

.page-size-selector select {
  width: auto;
}
</style>
