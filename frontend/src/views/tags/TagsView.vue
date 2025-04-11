<template>
  <div>
    <v-row class="mb-4">
      <v-col cols="12" sm="8">
        <h1 class="text-h3">Tags</h1>
        <p class="text-body-1">Manage tags for categorizing exercises</p>
      </v-col>
      <v-col cols="12" sm="4" class="d-flex justify-end align-center">
        <v-btn
          color="primary"
          prepend-icon="mdi-plus"
          @click="openCreateDialog"
        >
          New Tag
        </v-btn>
      </v-col>
    </v-row>

    <!-- Search and filter -->
    <v-card class="mb-6">
      <v-card-text>
        <v-row>
          <v-col cols="12" sm="6" md="4">
            <v-text-field
              v-model="search"
              label="Search tags"
              variant="outlined"
              density="compact"
              prepend-inner-icon="mdi-magnify"
              hide-details
              clearable
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-select
              v-model="categoryFilter"
              :items="categoriesForSelect"
              label="Filter by category"
              variant="outlined"
              density="compact"
              hide-details
              clearable
            ></v-select>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- Tags list -->
    <v-card>
      <v-card-text v-if="tagsStore.loading">
        <div class="d-flex justify-center align-center py-4">
          <v-progress-circular indeterminate color="primary"></v-progress-circular>
        </div>
      </v-card-text>
      
      <template v-else>
        <v-data-table
          v-if="filteredTags.length > 0"
          :headers="headers"
          :items="filteredTags"
          :items-per-page="10"
          class="elevation-1            "
        >
        <template #item.createdAt="{ item }">
            {{ appStore.formatDate(item.createdAt) }}
        </template>
          <template v-slot:item.categories="{ item }">
            <v-chip-group>
              <v-chip
                v-for="categoryId in item.categoryIds"
                :key="categoryId"
                size="small"
                color="primary"
                variant="outlined"
                class="mr-1"
              >
                {{ getCategoryName(categoryId) }}
              </v-chip>
            </v-chip-group>
          </template>
          
          <template v-slot:item.actions="{ item }">
            <v-btn
              icon
              variant="text"
              size="small"
              color="primary"
              @click="openEditDialog(item)"
            >
              <v-icon>mdi-pencil</v-icon>
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
          <v-icon icon="mdi-tag-off" size="64" color="grey-lighten-1" class="mb-4"></v-icon>
          <h3 class="text-h6 mb-2">No tags found</h3>
          <p class="text-body-2 text-grey">
            {{ 
              tagsStore.tags.length === 0 
                ? "You haven't created any tags yet" 
                : "No tags match your search criteria" 
            }}
          </p>
        </v-card-text>
      </template>
    </v-card>

    <!-- Create/Edit Dialog -->
    <tag-form-dialog
      v-model="dialogVisible"
      :tag="selectedTag"
      :is-edit="isEdit"
      @save="saveTag"
    />

    <!-- Delete Confirmation Dialog -->
    <v-dialog v-model="deleteDialog" max-width="500">
      <v-card>
        <v-card-title class="text-h5">Delete Tag</v-card-title>
        <v-card-text>
          Are you sure you want to delete the tag <strong>{{ selectedTag?.name }}</strong>? 
          This will remove it from all exercises. This action cannot be undone.
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="grey-darken-1" variant="text" @click="deleteDialog = false">Cancel</v-btn>
          <v-btn 
            color="error" 
            variant="flat" 
            @click="deleteTag" 
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
import { useTagsStore } from '@/stores/tags.js'
import { useCategoriesStore } from '@/stores/categories.js'
import { useAppStore } from '@/stores/app.js'
import TagFormDialog from '@/components/tags/TagFormDialog.vue'

const tagsStore = useTagsStore()
const categoriesStore = useCategoriesStore()
const appStore = useAppStore()

// Data
const search = ref('')
const categoryFilter = ref(null)
const dialogVisible = ref(false)
const isEdit = ref(false)
const selectedTag = ref(null)
const deleteDialog = ref(false)
const deleteLoading = ref(false)

// Table headers
const headers = [
  { title: 'Name', key: 'name' },
  { title: 'Categories', key: 'categories' },
  { title: 'Created', key: 'createdAt' },
  { title: 'Actions', key: 'actions', sortable: false }
]

// Computed
const filteredTags = computed(() => {
  let result = tagsStore.tags
  
  // Search filter
  if (search.value) {
    const searchLower = search.value.toLowerCase()
    result = result.filter(tag => 
      tag.name.toLowerCase().includes(searchLower)
    )
  }
  
  // Category filter
  if (categoryFilter.value) {
    result = result.filter(tag => 
      tag.categoryIds && tag.categoryIds.includes(parseInt(categoryFilter.value))
    )
  }
  
  return result
})

const categoriesForSelect = computed(() => {
  return categoriesStore.categories.map(category => ({
    title: category.name,
    value: category.id
  }))
})

// Methods
function getCategoryName(categoryId) {
  const category = categoriesStore.categoryById(categoryId)
  return category ? category.name : 'Unknown'
}

function openCreateDialog() {
  selectedTag.value = {
    name: '',
    categoryIds: []
  }
  isEdit.value = false
  dialogVisible.value = true
}

function openEditDialog(tag) {
  selectedTag.value = { ...tag }
  isEdit.value = true
  dialogVisible.value = true
}

async function saveTag(tagData) {
  try {
    if (isEdit.value) {
      const id = selectedTag.value.id
      await tagsStore.updateTag(id, tagData, 'name,categoryIds')
      appStore.showSuccessMessage(`Tag "${tagData.name}" updated successfully`)
    } else {
      await tagsStore.createTag(tagData)
      appStore.showSuccessMessage(`Tag "${tagData.name}" created successfully`)
    }
    dialogVisible.value = false
  } catch (error) {
    appStore.showErrorMessage(`Error saving tag: ${error.message}`)
  }
}

function confirmDelete(tag) {
  selectedTag.value = tag
  deleteDialog.value = true
}

async function deleteTag() {
  deleteLoading.value = true
  try {
    await tagsStore.deleteTag(selectedTag.value.id)
    appStore.showSuccessMessage(`Tag "${selectedTag.value.name}" deleted successfully`)
    deleteDialog.value = false
  } catch (error) {
    appStore.showErrorMessage(`Error deleting tag: ${error.message}`)
  } finally {
    deleteLoading.value = false
  }
}

// Lifecycle
onMounted(async () => {
  if (tagsStore.tags.length === 0) {
    await tagsStore.fetchTags()
  }
  
  if (categoriesStore.categories.length === 0) {
    await categoriesStore.fetchCategories()
  }
})
</script>
