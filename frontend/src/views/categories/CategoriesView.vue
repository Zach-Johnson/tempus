<template>
  <div>
    <v-row class="mb-4">
      <v-col cols="12" sm="8">
        <h1 class="text-h3">Categories</h1>
        <p class="text-body-1">Manage categories for organizing your exercises</p>
      </v-col>
      <v-col cols="12" sm="4" class="d-flex justify-end align-center">
        <v-btn
          color="primary"
          prepend-icon="mdi-plus"
          @click="openCreateDialog"
        >
          New Category
        </v-btn>
      </v-col>
    </v-row>

    <!-- Search -->
    <v-card class="mb-6">
      <v-card-text>
        <v-row>
          <v-col cols="12" sm="6" md="4">
            <v-text-field
              v-model="search"
              label="Search categories"
              variant="outlined"
              density="compact"
              prepend-inner-icon="mdi-magnify"
              hide-details
              clearable
            ></v-text-field>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- Categories list -->
    <v-card>
      <v-card-text v-if="categoriesStore.loading">
        <div class="d-flex justify-center align-center py-4">
          <v-progress-circular indeterminate color="primary"></v-progress-circular>
        </div>
      </v-card-text>
      
      <template v-else>
        <v-data-table
          v-if="filteredCategories.length > 0"
          :headers="headers"
          :items="filteredCategories"
          :items-per-page="10"
          class="elevation-1"
        >
          <template v-slot:item.description="{ item }">
            <span v-if="item.description">{{ truncateText(item.description, 80) }}</span>
            <span v-else class="text-grey">No description</span>
          </template>
          
          <template v-slot:item.created_at="{ item }">
            {{ formatDate(item.created_at) }}
          </template>

          <template v-slot:item.actions="{ item }">
            <v-btn
              icon
              variant="text"
              size="small"
              color="primary"
              @click="openEditDialog(item)"
              class="mr-1"
            >
              <v-icon>mdi-pencil</v-icon>
            </v-btn>
            <v-btn
              icon
              variant="text"
              size="small"
              color="primary"
              :to="{ name: 'category-detail', params: { id: item.id }}"
              class="mr-1"
            >
              <v-icon>mdi-eye</v-icon>
            </v-btn>
            <v-btn
              icon
              variant="text"
              size="small"
              color="error"
              @click="confirmDelete(item)"
            >
              <v-icon>mdi-delete</v-icon>
            </v-btn>
          </template>
        </v-data-table>
        
        <v-card-text v-else class="text-center py-8">
          <v-icon icon="mdi-folder-off" size="64" color="grey-lighten-1" class="mb-4"></v-icon>
          <h3 class="text-h6 mb-2">No categories found</h3>
          <p class="text-body-2 text-grey">
            {{ 
              categoriesStore.categories.length === 0 
                ? "You haven't created any categories yet" 
                : "No categories match your search criteria" 
            }}
          </p>
        </v-card-text>
      </template>
    </v-card>

    <!-- Create/Edit Dialog -->
    <category-form-dialog
      v-model="dialogVisible"
      :category="selectedCategory"
      :is-edit="isEdit"
      @save="saveCategory"
    />

    <!-- Delete Confirmation Dialog -->
    <v-dialog v-model="deleteDialog" max-width="500">
      <v-card>
        <v-card-title class="text-h5">Delete Category</v-card-title>
        <v-card-text>
          Are you sure you want to delete the category <strong>{{ selectedCategory?.name }}</strong>? 
          This will also remove it from all associated tags and exercises. This action cannot be undone.
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="grey-darken-1" variant="text" @click="deleteDialog = false">Cancel</v-btn>
          <v-btn 
            color="error" 
            variant="flat" 
            @click="deleteCategory" 
            :loading="deleteLoading"
          >
            Delete
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useCategoriesStore } from '@/stores/categories.js'
import { useAppStore } from '@/stores/app.js'
import CategoryFormDialog from '@/components/categories/CategoryFormDialog.vue'

const categoriesStore = useCategoriesStore()
const appStore = useAppStore()

// Data
const search = ref('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const selectedCategory = ref(null)
const deleteDialog = ref(false)
const deleteLoading = ref(false)

// Table headers
const headers = [
  { title: 'Name', key: 'name' },
  { title: 'Description', key: 'description' },
  { title: 'Created', key: 'created_at' },
  { title: 'Actions', key: 'actions', sortable: false }
]

// Computed
const filteredCategories = computed(() => {
  // First deduplicate the categories by ID
  const uniqueCategories = []
  const seenIds = new Set()
  
  categoriesStore.categories.forEach(category => {
    if (!seenIds.has(category.id)) {
      seenIds.add(category.id)
      uniqueCategories.push(category)
    }
  })
  
  let result = uniqueCategories
  
  // Search filter
  if (search.value) {
    const searchLower = search.value.toLowerCase()
    result = result.filter(category => 
      category.name.toLowerCase().includes(searchLower) || 
      (category.description && category.description.toLowerCase().includes(searchLower))
    )
  }
  
  return result
})

// Methods
function formatDate(dateString) {
  return appStore.formatDate(dateString)
}

function truncateText(text, maxLength) {
  if (!text) return ''
  return text.length > maxLength 
    ? text.substring(0, maxLength) + '...' 
    : text
}

function openCreateDialog() {
  selectedCategory.value = {
    name: '',
    description: ''
  }
  isEdit.value = false
  dialogVisible.value = true
}

function openEditDialog(category) {
  selectedCategory.value = { ...category }
  isEdit.value = true
  dialogVisible.value = true
}

async function saveCategory(categoryData) {
  try {
    if (isEdit.value) {
      const id = selectedCategory.value.id
      await categoriesStore.updateCategory(id, categoryData, 'name,description')
      appStore.showSuccessMessage(`Category "${categoryData.name}" updated successfully`)
    } else {
      await categoriesStore.createCategory(categoryData)
      appStore.showSuccessMessage(`Category "${categoryData.name}" created successfully`)
    }
    dialogVisible.value = false
  } catch (error) {
    appStore.showErrorMessage(`Error saving category: ${error.message}`)
  }
}

function confirmDelete(category) {
  selectedCategory.value = category
  deleteDialog.value = true
}

async function deleteCategory() {
  deleteLoading.value = true
  try {
    await categoriesStore.deleteCategory(selectedCategory.value.id)
    appStore.showSuccessMessage(`Category "${selectedCategory.value.name}" deleted successfully`)
    deleteDialog.value = false
  } catch (error) {
    appStore.showErrorMessage(`Error deleting category: ${error.message}`)
  } finally {
    deleteLoading.value = false
  }
}

// Lifecycle
onMounted(async () => {
  if (categoriesStore.categories.length === 0) {
    await categoriesStore.fetchCategories()
  }
})
</script>
