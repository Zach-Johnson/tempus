<template>
  <div>
    <v-row class="mb-4">
      <v-col cols="12" sm="8">
        <h1 class="text-h3">Exercises</h1>
        <p class="text-body-1">Manage your drumming exercises</p>
      </v-col>
      <v-col cols="12" sm="4" class="d-flex justify-end align-center">
        <v-btn color="primary" prepend-icon="mdi-plus" @click="openCreateDialog">
          New Exercise
        </v-btn>
      </v-col>
    </v-row>

    <!-- Search and filter -->
    <v-card class="mb-6">
      <v-card-text>
        <v-row>
          <v-col cols="12" sm="6" md="4">
            <v-text-field v-model="search" label="Search exercises" variant="outlined" density="compact"
              prepend-inner-icon="mdi-magnify" hide-details clearable></v-text-field>
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-select v-model="categoryFilter" :items="categoriesForSelect" label="Filter by category"
              variant="outlined" density="compact" hide-details clearable></v-select>
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-select v-model="tagFilter" :items="tagsForSelect" label="Filter by tag" variant="outlined"
              density="compact" hide-details clearable></v-select>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- Exercises list -->
    <v-card>
      <v-card-text v-if="exercisesStore.loading">
        <div class="d-flex justify-center align-center py-4">
          <v-progress-circular indeterminate color="primary"></v-progress-circular>
        </div>
      </v-card-text>

      <template v-else>
        <v-data-table v-if="filteredExercises.length > 0" :headers="headers" :items="filteredExercises" hover
          :items-per-page="50" class="elevation-1">
          <template v-slot:item.description="{ item }">
            <span v-if="item.description">{{ truncateText(item.description, 80) }}</span>
            <span v-else class="text-grey">No description</span>
          </template>

          <template v-slot:item.categories="{ item }">
            <v-chip-group>
              <template v-for="categoryId in exercisesStore.getCategoryIdsForExercise(item.id)" :key="categoryId">
                <category-chip v-if="getCategoryById(categoryId)" :category="getCategoryById(categoryId)" size="small"
                  class="mr-1"></category-chip>
              </template>
            </v-chip-group>
          </template>

          <template v-slot:item.tags="{ item }">
            <v-chip-group>
              <template v-for="tagId in item.tagIds || []" :key="tagId">
                <tag-chip v-if="getTagById(tagId)" :tag="getTagById(tagId)" size="small" class="mr-1"></tag-chip>
              </template>
            </v-chip-group>
          </template>

          <template v-slot:item.createdAt="{ item }">
            {{ formatDate(item.createdAt) }}
          </template>

          <template v-slot:item.actions="{ item }">
            <v-btn icon variant="text" size="small" color="primary" @click="openEditDialog(item)" class="mr-1">
              <v-icon>mdi-pencil</v-icon>
            </v-btn>
            <v-btn icon variant="text" size="small" color="primary"
              :to="{ name: 'exercise-detail', params: { id: item.id } }" class="mr-1">
              <v-icon>mdi-eye</v-icon>
            </v-btn>
            <v-btn icon variant="text" size="small" color="error" @click="confirmDelete(item)">
              <v-icon>mdi-delete</v-icon>
            </v-btn>
          </template>
        </v-data-table>

        <v-card-text v-else class="text-center py-8">
          <v-icon icon="mdi-music-note-off" size="64" color="grey-lighten-1" class="mb-4"></v-icon>
          <h3 class="text-h6 mb-2">No exercises found</h3>
          <p class="text-body-2 text-grey">
            {{
              exercisesStore.exercises.length === 0
                ? "You haven't created any exercises yet"
                : "No exercises match your search criteria"
            }}
          </p>
        </v-card-text>
      </template>
    </v-card>

    <!-- Create/Edit Dialog -->
    <exercise-form-dialog v-model="dialogVisible" :exercise="selectedExercise" :is-edit="isEdit" @save="saveExercise" />

    <!-- Delete Confirmation Dialog -->
    <v-dialog v-model="deleteDialog" max-width="500">
      <v-card>
        <v-card-title class="text-h5">Delete Exercise</v-card-title>
        <v-card-text>
          Are you sure you want to delete the exercise <strong>{{ selectedExercise?.name }}</strong>?
          This action cannot be undone and will remove this exercise from all practice sessions.
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="grey-darken-1" variant="text" @click="deleteDialog = false">Cancel</v-btn>
          <v-btn color="error" variant="flat" @click="deleteExercise" :loading="deleteLoading">
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
import { useExercisesStore } from '@/stores/exercises.js'
import { useCategoriesStore } from '@/stores/categories.js'
import { useTagsStore } from '@/stores/tags.js'
import { useAppStore } from '@/stores/app.js'
import ExerciseFormDialog from '@/components/exercises/ExerciseFormDialog.vue'
import CategoryChip from '@/components/categories/CategoryChip.vue'
import TagChip from '@/components/tags/TagChip.vue'

const router = useRouter()
const route = useRoute()
const exercisesStore = useExercisesStore()
const categoriesStore = useCategoriesStore()
const tagsStore = useTagsStore()
const appStore = useAppStore()

// Data
const search = ref('')
const categoryFilter = ref(null)
const tagFilter = ref(null)
const dialogVisible = ref(false)
const isEdit = ref(false)
const selectedExercise = ref(null)
const deleteDialog = ref(false)
const deleteLoading = ref(false)

// Table headers
const headers = [
  { title: 'Name', key: 'name' },
  { title: 'Description', key: 'description' },
  { title: 'Categories', key: 'categories' },
  { title: 'Tags', key: 'tags' },
  { title: 'Created', key: 'createdAt' },
  { title: 'Actions', key: 'actions', sortable: false }
]

// Computed
const filteredExercises = computed(() => {
  let result = exercisesStore.exercises

  // Search filter
  if (search.value) {
    const searchLower = search.value.toLowerCase()
    result = result.filter(exercise =>
      exercise.name.toLowerCase().includes(searchLower) ||
      (exercise.description && exercise.description.toLowerCase().includes(searchLower))
    )
  }

  // Category filter - using derived categories through tags
  if (categoryFilter.value) {
    result = result.filter(exercise => {
      const categoryIds = exercisesStore.getCategoryIdsForExercise(exercise.id)
      return categoryIds.includes(parseInt(categoryFilter.value))
    })
  }

  // Tag filter
  if (tagFilter.value) {
    result = result.filter(exercise =>
      exercise.tagIds && exercise.tagIds.includes(parseInt(tagFilter.value))
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

const tagsForSelect = computed(() => {
  return tagsStore.tags.map(tag => ({
    title: tag.name,
    value: tag.id
  }))
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

function getCategoryById(categoryId) {
  return categoriesStore.categoryById(categoryId)
}

function getTagById(tagId) {
  return tagsStore.tagById(tagId)
}

function openCreateDialog() {
  selectedExercise.value = {
    name: '',
    description: '',
    tagIds: [],
    links: []
  }
  isEdit.value = false
  dialogVisible.value = true
}

async function openEditDialog(exercise) {
  await exercisesStore.fetchExercise(exercise.id)
  selectedExercise.value = exercisesStore.currentExercise
  isEdit.value = true
  dialogVisible.value = true
}

async function saveExercise(exerciseData) {
  try {
    if (isEdit.value) {
      const id = selectedExercise.value.id
      await exercisesStore.updateExercise(id, exerciseData, 'name,description,tagIds')
      appStore.showSuccessMessage(`Exercise "${exerciseData.name}" updated successfully`)
    } else {
      await exercisesStore.createExercise(exerciseData)
      appStore.showSuccessMessage(`Exercise "${exerciseData.name}" created successfully`)
    }
    dialogVisible.value = false
  } catch (error) {
    appStore.showErrorMessage(`Error saving exercise: ${error.message}`)
  }
}

function confirmDelete(exercise) {
  selectedExercise.value = exercise
  deleteDialog.value = true
}

async function deleteExercise() {
  deleteLoading.value = true
  try {
    await exercisesStore.deleteExercise(selectedExercise.value.id)
    appStore.showSuccessMessage(`Exercise "${selectedExercise.value.name}" deleted successfully`)
    deleteDialog.value = false
  } catch (error) {
    appStore.showErrorMessage(`Error deleting exercise: ${error.message}`)
  } finally {
    deleteLoading.value = false
  }
}

// Watch for query parameters to apply filters
watch(() => route.query, (query) => {
  if (query.category) {
    categoryFilter.value = parseInt(query.category)
  }

  if (query.tag) {
    tagFilter.value = parseInt(query.tag)
  }
}, { immediate: true })

// Lifecycle
onMounted(async () => {
  // Load data if not already loaded
  if (exercisesStore.exercises.length === 0) {
    await exercisesStore.fetchExercises()
  }

  if (categoriesStore.categories.length === 0) {
    await categoriesStore.fetchCategories()
  }

  if (tagsStore.tags.length === 0) {
    await tagsStore.fetchTags()
  }
})
</script>
