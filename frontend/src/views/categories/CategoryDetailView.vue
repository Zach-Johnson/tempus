<template>
  <div>
    <!-- Loading state -->
    <div v-if="loading" class="d-flex justify-center my-8">
      <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
    </div>

    <!-- Error state -->
    <v-alert v-else-if="error" type="error" class="my-4">
      {{ error }}
    </v-alert>

    <!-- Category details -->
    <template v-else>
      <div v-if="!categoriesStore.currentCategory">
        <v-alert type="info" class="my-4">
          Category not found. It may have been deleted.
          <div class="mt-4">
            <v-btn color="primary" to="/categories">
              Back to Categories
            </v-btn>
          </div>
        </v-alert>
      </div>
      <template v-else>
        <!-- Header section with actions -->
        <v-row class="mb-4">
          <v-col cols="12" sm="8">
            <v-btn
              variant="text"
              to="/categories"
              prepend-icon="mdi-arrow-left"
              class="mb-2"
            >
              Back to Categories
            </v-btn>
            <h1 class="text-h3">{{ categoriesStore.currentCategory.name }}</h1>
          </v-col>
          <v-col cols="12" sm="4" class="d-flex justify-end align-center">
            <v-btn
              color="primary"
              variant="outlined"
              class="mr-2"
              prepend-icon="mdi-pencil"
              @click="openEditDialog"
            >
              Edit
            </v-btn>
            <v-btn
              color="error"
              variant="outlined"
              prepend-icon="mdi-delete"
              @click="confirmDelete"
            >
              Delete
            </v-btn>
          </v-col>
        </v-row>

        <!-- Category details card -->
        <v-card class="mb-6">
          <v-card-text>
            <v-row>
              <v-col cols="12" md="6">
                <div class="text-body-1 font-weight-medium mb-1">Description</div>
                <div class="text-body-2 mb-4">{{ categoriesStore.currentCategory.description || 'No description provided' }}</div>
                
                <div class="text-body-1 font-weight-medium mb-1">Created</div>
                <div class="text-body-2 mb-4">{{ formatDate(categoriesStore.currentCategory.created_at) }}</div>
                
                <div class="text-body-1 font-weight-medium mb-1">Last Updated</div>
                <div class="text-body-2">{{ formatDate(categoriesStore.currentCategory.updated_at) }}</div>
              </v-col>
              <v-col cols="12" md="6">
                <div class="d-flex align-center mb-2">
                  <div class="text-body-1 font-weight-medium">Tags</div>
                  <v-btn
                    variant="text"
                    density="compact"
                    icon="mdi-plus"
                    size="small"
                    color="primary"
                    class="ml-2"
                    @click="goToTags"
                    title="Add tags"
                  ></v-btn>
                </div>
                <div v-if="loading" class="d-flex align-center mb-4">
                  <v-progress-circular indeterminate size="20" width="2" class="mr-2"></v-progress-circular>
                  <span class="text-body-2">Loading tags...</span>
                </div>
                <div v-else-if="categoryTags.length === 0" class="text-body-2 mb-4 text-grey">
                  No tags associated with this category
                </div>
                <div v-else class="mb-4">
                  <tag-list :tags="categoryTags" />
                </div>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>

        <!-- Exercises in this category -->
        <h2 class="text-h5 mb-4">Exercises in this Category</h2>
        <v-card>
          <v-card-text v-if="loadingExercises">
            <div class="d-flex justify-center py-4">
              <v-progress-circular indeterminate color="primary"></v-progress-circular>
            </div>
          </v-card-text>
          
          <v-card-text v-else-if="categoryExercises.length === 0" class="text-center py-8">
            <v-icon icon="mdi-music-note-off" size="64" color="grey-lighten-1" class="mb-4"></v-icon>
            <h3 class="text-h6 mb-2">No exercises found</h3>
            <p class="text-body-2 text-grey mb-4">
              This category doesn't have any exercises yet
            </p>
            <v-btn 
              color="primary" 
              prepend-icon="mdi-plus" 
              :to="{ name: 'exercises', query: { create: true, category: categoryId } }"
            >
              Create Exercise
            </v-btn>
          </v-card-text>
          
          <template v-else>
            <v-data-table
              :headers="exerciseHeaders"
              :items="categoryExercises"
              :items-per-page="5"
            >
              <template v-slot:item.actions="{ item }">
                <v-btn
                  icon
                  variant="text"
                  size="small"
                  color="primary"
                  :to="{ name: 'exercise-detail', params: { id: item.id }}"
                >
                  <v-icon>mdi-eye</v-icon>
                </v-btn>
              </template>
            </v-data-table>
            
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn 
                color="primary" 
                variant="text"
                :to="{ name: 'exercises', query: { category: categoryId } }"
              >
                View All Exercises
              </v-btn>
            </v-card-actions>
          </template>
        </v-card>
      </template>
    </template>

    <!-- Edit Dialog -->
    <category-form-dialog
      v-if="categoriesStore.currentCategory"
      v-model="dialogVisible"
      :category="categoriesStore.currentCategory"
      :is-edit="true"
      @save="saveCategory"
    />

    <!-- Delete Confirmation Dialog -->
    <v-dialog v-model="deleteDialog" max-width="500">
      <v-card>
        <v-card-title class="text-h5">Delete Category</v-card-title>
        <v-card-text>
          Are you sure you want to delete the category <strong>{{ categoriesStore.currentCategory?.name }}</strong>? 
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
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useCategoriesStore } from '@/stores/categories.js'
import { useTagsStore } from '@/stores/tags.js'
import { useExercisesStore } from '@/stores/exercises.js'
import { useAppStore } from '@/stores/app.js'
import { categoriesAPI } from '@/services/api.js'
import CategoryFormDialog from '@/components/categories/CategoryFormDialog.vue'
import TagList from '@/components/tags/TagList.vue'

const router = useRouter()
const route = useRoute()
const categoriesStore = useCategoriesStore()
const tagsStore = useTagsStore()
const exercisesStore = useExercisesStore()
const appStore = useAppStore()

// Data
const loading = ref(true)
const loadingExercises = ref(true)
const error = ref(null)
const dialogVisible = ref(false)
const deleteDialog = ref(false)
const deleteLoading = ref(false)

// Computed
const categoryId = computed(() => parseInt(route.params.id))

const category = computed(() => {
  // First check if the category is already in the store
  const categoryFromStore = categoriesStore.categoryById(categoryId.value)
  
  // Then check if we have a currentCategory from a direct fetch
  const currentCategory = categoriesStore.currentCategory
  
  // Return the category if it's found in either place
  if (categoryFromStore) {
    return categoryFromStore
  } else if (currentCategory && currentCategory.id === categoryId.value) {
    return currentCategory
  }
  
  return null
})

const categoryTags = computed(() => {
  return tagsStore.tagsByCategory(categoryId.value) || []
})

const categoryExercises = computed(() => {
  return exercisesStore.exercisesByCategory(categoryId.value) || []
})

// Table headers
const exerciseHeaders = [
  { title: 'Name', key: 'name' },
  { title: 'Description', key: 'description' },
  { title: 'Actions', key: 'actions', sortable: false, align: 'end' }
]

// Methods
function formatDate(dateString) {
  return appStore.formatDate(dateString)
}

function openEditDialog() {
  dialogVisible.value = true
}

function goToTags() {
  router.push({ name: 'tags', query: { category: categoryId.value } })
}

async function saveCategory(categoryData) {
  try {
    await categoriesStore.updateCategory(categoryId.value, categoryData, 'name,description')
    appStore.showSuccessMessage(`Category "${categoryData.name}" updated successfully`)
    dialogVisible.value = false
  } catch (error) {
    appStore.showErrorMessage(`Error updating category: ${error.message}`)
  }
}

function confirmDelete() {
  deleteDialog.value = true
}

async function deleteCategory() {
  deleteLoading.value = true
  try {
    await categoriesStore.deleteCategory(categoryId.value)
    appStore.showSuccessMessage(`Category "${categoriesStore.currentCategory.name}" deleted successfully`)
    deleteDialog.value = false
    router.push({ name: 'categories' })
  } catch (error) {
    appStore.showErrorMessage(`Error deleting category: ${error.message}`)
  } finally {
    deleteLoading.value = false
  }
}

// Load data
async function loadData() {
  loading.value = true
  error.value = null
  
  try {
    console.log(`Loading category with ID: ${categoryId.value}`)
    
    // Always fetch the category to ensure we have the latest data
    const fetchedCategory = await categoriesStore.fetchCategory(categoryId.value)
    console.log('Fetched category:', fetchedCategory)
    console.log('Current category from store:', categoriesStore.currentCategory)
    
    // If the category is still not found after fetching, show an error
    if (!categoriesStore.currentCategory) {
      console.warn(`Category with ID ${categoryId.value} not found after fetching`)
      error.value = `Category with ID ${categoryId.value} could not be found`
    } else {
      console.log('Category found:', categoriesStore.currentCategory)
    }
    
    if (tagsStore.tags.length === 0) {
      await tagsStore.fetchTags()
    }
  } catch (err) {
    console.error(`Error loading category: ${err.message}`, err)
    error.value = `Error loading category: ${err.message}`
  } finally {
    loading.value = false
  }
  
  loadingExercises.value = true
  try {
    if (exercisesStore.exercises.length === 0) {
      await exercisesStore.fetchExercises({ category_id: categoryId.value })
    }
  } catch (err) {
    console.error('Error loading exercises:', err)
  } finally {
    loadingExercises.value = false
  }
}

// Watch for route param changes to reload data
watch(() => route.params.id, () => {
  loadData()
})

// Lifecycle
onMounted(async () => {
  await loadData()
  
  // One final check to force a direct category load if still not found
  if (!categoriesStore.currentCategory && categoryId.value) {
    try {
      // Try to manually get the category by ID
      const response = await categoriesAPI.get(categoryId.value)
      const manualCategory = response.data
      
      // Force update the store
      categoriesStore.$patch(state => {
        state.currentCategory = manualCategory
        
        const index = state.categories.findIndex(c => c.id === categoryId.value)
        if (index !== -1) {
          state.categories[index] = manualCategory
        } else {
          state.categories.push(manualCategory)
        }
      })
      
      console.log('Manually loaded category:', manualCategory)
    } catch (err) {
      console.error('Final attempt to load category failed:', err)
    }
  }
})
</script>
