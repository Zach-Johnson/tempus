<template>
  <div class="category-list">
    <div v-if="loading" class="d-flex align-center">
      <v-progress-circular indeterminate size="20" width="2" color="primary" class="mr-2"></v-progress-circular>
      <span class="text-body-2">Loading categories...</span>
    </div>
    
    <div v-else-if="!categories || categories.length === 0" class="text-body-2 text-grey">
      {{ noCategoriesMessage }}
    </div>
    
    <div v-else class="category-list-container">
      <category-chip
        v-for="category in categories"
        :key="category.id"
        :category="category"
        :color="color"
        :variant="variant"
        :size="size"
        :closable="selectable && modelValue && modelValue.includes(category.id)"
        :clickable="selectable"
        :disable-navigation="disableNavigation"
        @click="toggleCategory(category)"
        @close="removeCategory(category)"
      ></category-chip>
      
      <div v-if="showAddButton && selectable" class="d-inline-block">
        <v-btn
          variant="outlined"
          :size="size"
          icon
          class="category-add-btn"
          @click="openCategorySelect"
        >
          <v-icon>mdi-plus</v-icon>
        </v-btn>
      </div>
    </div>
    
    <!-- Category Selection Dialog -->
    <v-dialog v-model="categorySelectDialog" max-width="500">
      <v-card>
        <v-card-title>Add Categories</v-card-title>
        <v-card-text>
          <v-text-field
            v-model="categorySearch"
            label="Search categories"
            variant="outlined"
            density="compact"
            prepend-inner-icon="mdi-magnify"
            clearable
            class="mb-4"
          ></v-text-field>
          
          <v-list v-if="availableCategories.length > 0" density="compact">
            <v-list-item
              v-for="category in filteredAvailableCategories"
              :key="category.id"
              :value="category.id"
              @click="addCategory(category)"
            >
              <template v-slot:prepend>
                <v-checkbox-btn
                  :model-value="isCategorySelected(category)"
                  color="primary"
                ></v-checkbox-btn>
              </template>
              <v-list-item-title>{{ category.name }}</v-list-item-title>
              <v-list-item-subtitle v-if="category.description">
                {{ truncateText(category.description, 60) }}
              </v-list-item-subtitle>
            </v-list-item>
          </v-list>
          
          <div v-else class="text-center py-4">
            <p class="text-body-2 text-grey mb-4">No categories available</p>
            <v-btn 
              color="primary"
              prepend-icon="mdi-plus"
              @click="openCreateDialog"
            >
              Create New Category
            </v-btn>
          </div>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn 
            color="primary" 
            variant="text" 
            @click="categorySelectDialog = false"
          >
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    
    <!-- Create Category Dialog -->
    <category-form-dialog
      v-model="createDialogVisible"
      :category="{ name: '', description: '' }"
      :is-edit="false"
      @save="createCategory"
    />
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useCategoriesStore } from '@/stores/categories.js'
import { useAppStore } from '@/stores/app.js'
import CategoryChip from './CategoryChip.vue'
import CategoryFormDialog from './CategoryFormDialog.vue'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  },
  categories: {
    type: Array,
    default: () => []
  },
  color: {
    type: String,
    default: 'secondary'
  },
  variant: {
    type: String,
    default: 'tonal'
  },
  size: {
    type: String,
    default: 'small'
  },
  selectable: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  noCategoriesMessage: {
    type: String,
    default: 'No categories'
  },
  showAddButton: {
    type: Boolean,
    default: true
  },
  disableNavigation: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'add', 'remove', 'category-created'])
const categoriesStore = useCategoriesStore()
const appStore = useAppStore()

// Data
const categorySelectDialog = ref(false)
const createDialogVisible = ref(false)
const categorySearch = ref('')

// Computed
const availableCategories = computed(() => {
  // Get all categories from the store
  return categoriesStore.categoriesSorted
})

const filteredAvailableCategories = computed(() => {
  if (!categorySearch.value) return availableCategories.value
  
  const search = categorySearch.value.toLowerCase()
  return availableCategories.value.filter(category => 
    category.name.toLowerCase().includes(search) ||
    (category.description && category.description.toLowerCase().includes(search))
  )
})

const selectedCategories = computed(() => {
  if (!props.modelValue) return []
  
  // Convert category IDs to category objects
  return props.modelValue
    .map(categoryId => {
      const category = categoriesStore.categoryById(categoryId)
      return category || null
    })
    .filter(category => category !== null)
})

// Methods
function truncateText(text, maxLength) {
  if (!text) return ''
  return text.length > maxLength 
    ? text.substring(0, maxLength) + '...' 
    : text
}

function isCategorySelected(category) {
  return props.modelValue && props.modelValue.includes(category.id)
}

function openCategorySelect() {
  categorySearch.value = ''
  categorySelectDialog.value = true
}

function openCreateDialog() {
  categorySelectDialog.value = false
  createDialogVisible.value = true
}

function toggleCategory(category) {
  if (!props.selectable) return
  
  if (isCategorySelected(category)) {
    removeCategory(category)
  } else {
    addCategory(category)
  }
}

function addCategory(category) {
  if (!props.selectable) return
  
  // Skip if already selected
  if (isCategorySelected(category)) return
  
  const newValue = [...(props.modelValue || []), category.id]
  emit('update:modelValue', newValue)
  emit('add', category)
}

function removeCategory(category) {
  if (!props.selectable) return
  
  const newValue = (props.modelValue || []).filter(id => id !== category.id)
  emit('update:modelValue', newValue)
  emit('remove', category)
}

async function createCategory(categoryData) {
  try {
    const newCategory = await categoriesStore.createCategory(categoryData)
    appStore.showSuccessMessage(`Category "${categoryData.name}" created successfully`)
    createDialogVisible.value = false
    
    // Add the new category to the selection
    if (props.selectable) {
      addCategory(newCategory)
    }
    
    emit('category-created', newCategory)
  } catch (error) {
    appStore.showErrorMessage(`Error creating category: ${error.message}`)
  }
}

// Fetch categories if needed
onMounted(async () => {
  if (categoriesStore.categories.length === 0) {
    await categoriesStore.fetchCategories()
  }
})
</script>

<style scoped>
.category-list-container {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
}

.category-add-btn {
  margin-right: 4px;
  margin-bottom: 4px;
}
</style>
