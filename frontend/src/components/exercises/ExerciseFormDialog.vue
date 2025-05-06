<template>
  <v-dialog v-model="dialogModel" max-width="800" persistent>
    <v-card>
      <v-card-title class="text-h5">
        {{ isEdit ? `Edit Exercise: ${originalExercise?.name}` : 'Create New Exercise' }}
      </v-card-title>

      <v-card-text>
        <v-form ref="form" @submit.prevent="save" v-model="formValid">
          <v-row>
            <v-col cols="12">
              <v-text-field v-model="formData.name" label="Exercise Name" :rules="nameRules" required variant="outlined"
                autofocus></v-text-field>
            </v-col>

            <v-col cols="12">
              <v-textarea v-model="formData.description" label="Description" variant="outlined" rows="4"
                placeholder="Provide a description for this exercise (optional)"></v-textarea>
            </v-col>

            <v-col cols="12">
              <div class="mb-2">
                <v-select v-model="formData.tagIds" :items="tagsStore.tags" item-title="name" item-value="id"
                  label="Tags" multiple chips variant="outlined" :loading="tagsStore.loading">
                  <template v-slot:selection="{ item }">
                    <v-chip :key="item.raw.id" title="item.raw.name" closable @click:close="removeTag(item.raw.id)">
                      {{ item.raw.name }}
                      <div v-if="item.raw.categoryIds?.length" class="d-flex align-center mt-1">
                        <v-chip v-for="categoryId in item.raw.categoryIds" :key="`cat-${item.raw.id}-${categoryId}`"
                          size="x-small" :color="getCategoryColor(categoryId)" class="ml-1">
                          {{ getCategoryName(categoryId) }}
                        </v-chip>
                      </div>
                    </v-chip>
                  </template>
                  <template v-slot:item="{ item, props }">
                    <v-list-item v-bind="props">
                      <template v-slot:prepend>
                        <v-checkbox-btn :model-value="isTagSelected(item.raw.id)"></v-checkbox-btn>
                      </template>
                      <template v-slot:title>
                        {{ item.raw.name }}
                      </template>
                      <template v-slot:subtitle>
                        <div v-if="item.raw.categoryIds?.length" class="d-flex align-center mt-1">
                          <v-chip v-for="categoryId in item.raw.categoryIds" :key="`cat-${item.raw.id}-${categoryId}`"
                            size="x-small" :color="getCategoryColor(categoryId)" class="mr-1">
                            {{ getCategoryName(categoryId) }}
                          </v-chip>
                        </div>
                      </template>
                    </v-list-item>
                  </template>
                  <template v-slot:append-inner>
                    <v-btn icon size="small" variant="text" class="ms-2" @click.stop="openTagForm">
                      <v-icon>mdi-plus</v-icon>
                    </v-btn>
                  </template>
                </v-select>
              </div>
            </v-col>

            <v-col cols="12" v-if="!isEdit">
              <div class="d-flex align-center mb-2">
                <div class="text-body-1 font-weight-medium">Exercise Images</div>
              </div>

              <div class="paste-area pa-4 text-center mb-4" tabindex="0" @paste="handlePaste"
                @focus="$event.target.classList.add('focused')" @blur="$event.target.classList.remove('focused')">
                <div v-if="imagePreviewUrls.length === 0">
                  <v-icon icon="mdi-image-plus" size="48" color="grey-lighten-1" class="mb-2"></v-icon>
                  <p class="text-body-2">Click here and paste an image from your clipboard</p>
                  <p class="text-caption text-grey">Press Ctrl+V or Cmd+V after clicking this area</p>
                </div>

                <v-row v-else>
                  <v-col v-for="(preview, index) in imagePreviewUrls" :key="index" cols="6" sm="4" md="3">
                    <v-card class="mb-2">
                      <v-img :src="preview.url" height="120" cover></v-img>
                      <v-card-text class="pa-2">
                        <div class="text-caption text-truncate">{{ preview.name }}</div>
                        <div class="text-caption text-grey">{{ preview.size }}</div>
                      </v-card-text>
                      <v-card-actions class="pa-2 pt-0">
                        <v-spacer></v-spacer>
                        <v-btn icon variant="text" size="small" color="error" @click="removeImage(index)">
                          <v-icon>mdi-delete</v-icon>
                        </v-btn>
                      </v-card-actions>
                    </v-card>
                  </v-col>
                </v-row>
              </div>
            </v-col>

            <!-- External Resources Section -->
            <v-col cols="12">
              <div class="d-flex align-center mb-2">
                <div class="text-body-1 font-weight-medium">External Resources</div>
                <v-btn variant="text" density="compact" icon="mdi-plus" size="small" color="primary" class="ml-2"
                  @click="addExternalLink" title="Add external resource"></v-btn>
              </div>

              <div v-if="formData.links.length === 0" class="text-body-2 text-grey mb-4">
                No external resources added yet
              </div>

              <div v-else>
                <v-list density="compact">
                  <v-list-item v-for="(link, index) in formData.links" :key="index">
                    <template v-slot:prepend>
                      <v-icon icon="mdi-link" size="small" class="mr-2"></v-icon>
                    </template>

                    <v-list-item-title>
                      <a :href="link.url" target="_blank" rel="noopener noreferrer">
                        {{ link.description || link.url }}
                      </a>
                    </v-list-item-title>

                    <template v-slot:append>
                      <v-btn icon variant="text" size="small" color="error" @click="removeLink(index)">
                        <v-icon>mdi-delete</v-icon>
                      </v-btn>
                    </template>
                  </v-list-item>
                </v-list>
              </div>
            </v-col>

            <!-- Images Section -->
            <v-col v-if="isEdit" cols="12">
              <div class="d-flex align-center mb-2">
                <div class="text-body-1 font-weight-medium">Images</div>
                <v-btn variant="text" density="compact" icon="mdi-plus" size="small" color="primary" class="ml-2"
                  @click="addImage" title="Add image"></v-btn>
              </div>
            </v-col>
          </v-row>
          <!-- Exercise images card -->
          <v-container v-if="exercise.images && exercise.images.length > 0">
            <v-row class="justify-center">
              <v-col cols="12" sm="8" md="6" v-for="(image, index) in exercise.images" :key="index"
                class="d-flex justify-center">
                <v-card class="mb-4" width="100%" max-width="1000">
                  <v-img :src="'data:image/jpeg;base64,' + image.imageData" class="bg-grey-lighten-4" contain></v-img>
                  <v-card-actions class="pa-2 pt-0">
                    <v-spacer></v-spacer>
                    <v-btn icon variant="text" size="small" color="error" @click="confirmDeleteImage(image)">
                      <v-icon>mdi-delete</v-icon>
                    </v-btn>
                  </v-card-actions>
                </v-card>
              </v-col>
            </v-row>
          </v-container>
        </v-form>
      </v-card-text>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="grey-darken-1" variant="text" @click="close">
          Cancel
        </v-btn>
        <v-btn color="primary" variant="flat" @click="save" :disabled="!formValid || !formData.name" :loading="saving">
          {{ isEdit ? 'Update' : 'Create' }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- External Link Dialog -->
  <v-dialog v-model="linkDialog" max-width="600">
    <v-card>
      <v-card-title>Add External Resource</v-card-title>
      <v-card-text>
        <v-form ref="linkForm" @submit.prevent="saveLink" v-model="linkFormValid">
          <v-text-field v-model="linkFormData.url" label="URL" :rules="urlRules" required variant="outlined"
            placeholder="https://example.com" class="mb-4"></v-text-field>

          <v-text-field v-model="linkFormData.description" label="Description" variant="outlined"
            placeholder="Description of the resource (optional)"></v-text-field>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="grey-darken-1" variant="text" @click="linkDialog = false">Cancel</v-btn>
        <v-btn color="primary" variant="flat" @click="saveLink" :disabled="!linkFormValid">
          Add
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Image Dialog -->
  <v-dialog v-model="imageDialog" max-width="600">
    <v-card>
      <v-card-title>Exercise Images</v-card-title>
      <v-card-text>
        <div class="paste-area pa-4 text-center mb-4" tabindex="0" @paste="handlePasteAndSave"
          @focus="$event.target.classList.add('focused')" @blur="$event.target.classList.remove('focused')">
          <v-icon icon="mdi-image-plus" size="48" color="grey-lighten-1" class="mb-2"></v-icon>
          <p class="text-body-2">Click here and paste an image from your clipboard</p>
          <p class="text-caption text-grey">Press Ctrl+V or Cmd+V after clicking this area</p>
        </div>
      </v-card-text>
    </v-card>
  </v-dialog>

  <!-- Tag Form Dialog -->
  <tag-form-dialog v-model="tagFormDialog" :tag="{ name: '', category_ids: [] }" :is-edit="false"
    @save="onTagCreated" />

  <!-- Delete Image Confirmation Dialog -->
  <v-dialog v-model="deleteImageDialog" max-width="500">
    <v-card>
      <v-card-title class="text-h5">Delete Image</v-card-title>
      <v-card-text>
        Are you sure?
        This action cannot be undone.
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="grey-darken-1" variant="text" @click="deleteImageDialog = false">Cancel</v-btn>
        <v-btn color="error" variant="flat" @click="deleteImage" :loading="deleteLoading">
          Delete
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useTagsStore } from '@/stores/tags.js'
import { useCategoriesStore } from '@/stores/categories.js'
import { useExercisesStore } from '@/stores/exercises.js'
import { useAppStore } from '@/stores/app.js'
import TagFormDialog from '@/components/tags/TagFormDialog.vue'

// Props
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  exercise: {
    type: Object,
    default: null
  },
  isEdit: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits(['update:modelValue', 'save'])

// Stores
const tagsStore = useTagsStore()
const categoriesStore = useCategoriesStore()
const exerciseStore = useExercisesStore()
const appStore = useAppStore()

// Data
const form = ref(null)
const formValid = ref(false)
const saving = ref(false)
const formData = ref({
  name: '',
  description: '',
  tagIds: [],
  links: []
})
const originalExercise = ref(null)

// External Link Dialog
const linkDialog = ref(false)
const linkForm = ref(null)
const linkFormValid = ref(false)
const linkFormData = ref({
  url: '',
  description: ''
})

// Image creation
const imageDialog = ref(false)
const imageFiles = ref([])
const imagePreviewUrls = ref([])
const selectedImage = ref(null) // Used for deleting images
const deleteImageDialog = ref(false)
const deleteLoading = ref(false)

// Tag Form Dialog
const tagFormDialog = ref(false)

// Category color mapping
const categoryColors = {
  1: '#1976D2', // Primary blue
  2: '#E53935', // Red
  3: '#43A047', // Green
  4: '#FB8C00', // Orange
  5: '#8E24AA', // Purple
  6: '#00ACC1', // Cyan
  7: '#FFB300', // Amber
  8: '#5E35B1', // Deep Purple
  9: '#1E88E5', // Blue
  10: '#00897B', // Teal
}

// Computed
const dialogModel = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// Validation rules
const nameRules = [
  v => !!v || 'Exercise name is required',
  v => (v && v.length >= 2) || 'Exercise name must be at least 2 characters',
  v => (v && v.length <= 50) || 'Exercise name must be less than 50 characters'
]

const urlRules = [
  v => !!v || 'URL is required',
  v => /^(https?:\/\/)?[\da-z.-]+\.[a-z.]{2,6}.*\/?$/.test(v) || 'Enter a valid URL'
]

// Methods
function close() {
  dialogModel.value = false
}

async function save() {
  if (!formValid.value) return

  saving.value = true
  try {
    // Prepare the data to emit
    const exerciseData = {
      name: formData.value.name,
      description: formData.value.description || '',
      tagIds: formData.value.tagIds.map(id => Number(id)),
      links: formData.value.links
    }

    // Add images if there are any
    if (imageFiles.value.length > 0) {
      exerciseData.images = await Promise.all(imageFiles.value.map(async (file) => {
        const arrayBuffer = await file.arrayBuffer()
        const base64 = btoa(
          new Uint8Array(arrayBuffer)
            .reduce((data, byte) => data + String.fromCharCode(byte), '')
        )
        return {
          imageData: base64,
          filename: file.name || 'pasted-image.png',
          mimeType: file.type || 'image/png',
          description: ''
        }
      }))
    }

    // Emit save event with the complete data
    emit('save', exerciseData)
  } catch (error) {
    console.error('Error saving exercise:', error)
  } finally {
    saving.value = false
  }
}

function isTagSelected(tagId) {
  return formData.value.tagIds.includes(tagId)
}

function removeTag(tagId) {
  formData.value.tagIds = formData.value.tagIds.filter(id => id !== tagId)
}

function getCategoryName(categoryId) {
  const category = categoriesStore.categoryById(categoryId)
  return category ? category.name : ''
}

function getCategoryColor(categoryId) {
  return categoryColors[categoryId] || '#9e9e9e'
}

function addExternalLink() {
  linkFormData.value = {
    url: '',
    description: ''
  }
  linkDialog.value = true
}

function saveLink() {
  if (!linkFormValid.value) return

  formData.value.links.push({
    url: linkFormData.value.url,
    description: linkFormData.value.description
  })

  linkDialog.value = false
}

function removeLink(index) {
  formData.value.links.splice(index, 1)
}

function addImage() {
  imageDialog.value = true
}

function openTagForm() {
  tagFormDialog.value = true
}

function onTagCreated(tagData) {
  // When a new tag is created, refresh tags and add the new tag to the selected tags
  tagsStore.fetchTags().then(() => {
    // Find the newly created tag by name
    const newTag = tagsStore.tags.find(tag => tag.name === tagData.name)
    if (newTag && !formData.value.tagIds.includes(newTag.id)) {
      formData.value.tagIds.push(newTag.id)
    }
  })
}

function handlePaste(event) {
  const items = (event.clipboardData || event.originalEvent.clipboardData).items

  for (const item of items) {
    if (item.kind === 'file' && item.type.startsWith('image/')) {
      const file = item.getAsFile()
      if (!file) continue

      // Add file to the array
      imageFiles.value.push(file)

      // Create a preview
      const reader = new FileReader()
      reader.onload = (e) => {
        imagePreviewUrls.value.push({
          url: e.target.result,
          name: file.name || `pasted-image.${file.type.split('/')[1]}`,
          size: formatFileSize(file.size),
          type: file.type
        })
      }
      reader.readAsDataURL(file)

      break // Only process one image at a time
    }
  }
}

async function handlePasteAndSave(event) {
  const items = (event.clipboardData || event.originalEvent.clipboardData).items

  saving.value = true
  try {
    for (const item of items) {
      if (item.kind === 'file' && item.type.startsWith('image/')) {
        const file = item.getAsFile()
        if (!file) continue

        const arrayBuffer = await file.arrayBuffer()
        const base64 = btoa(
          new Uint8Array(arrayBuffer)
            .reduce((data, byte) => data + String.fromCharCode(byte), '')
        )
        const image = {
          image_data: base64,
          filename: file.name || 'pasted-image.png',
          mime_type: file.type || 'image/png',
          description: ''
        }

        await exerciseStore.addExerciseImage(props.exercise.id, image)

        break // Only process one image at a time
      }
    }
  } catch (error) {
    appStore.showErrorMessage(`Error adding exercise image : ${error.message}`)
  } finally {
    imageDialog.value = false
    saving.value = false
  }
}

// Helper function to format file sizes
function formatFileSize(bytes) {
  if (bytes < 1024) return bytes + ' bytes'
  else if (bytes < 1048576) return (bytes / 1024).toFixed(1) + ' KB'
  else return (bytes / 1048576).toFixed(1) + ' MB'
}

// Remove an image
function removeImage(index) {
  imageFiles.value.splice(index, 1)
  imagePreviewUrls.value.splice(index, 1)
}

function confirmDeleteImage(image) {
  selectedImage.value = image
  deleteImageDialog.value = true
}

async function deleteImage() {
  deleteLoading.value = true
  try {
    await exerciseStore.deleteExerciseImage(selectedImage.value.id)
    appStore.showSuccessMessage(`Image deleted successfully`)
    deleteImageDialog.value = false
  } catch (error) {
    appStore.showErrorMessage(`Error deleting image: ${error.message}`)
  } finally {
    deleteLoading.value = false
  }
}

// Reset form when dialog opens/closes
// Extend your existing watch function
watch(() => props.modelValue, (isOpen) => {
  if (isOpen && props.exercise) {
    // Clone the exercise to avoid modifying the original
    originalExercise.value = { ...props.exercise }

    // Set form values
    formData.value = {
      name: props.exercise.name || '',
      description: props.exercise.description || '',
      tagIds: props.exercise.tagIds ? [...props.exercise.tagIds] : [],
      links: props.exercise.links ? JSON.parse(JSON.stringify(props.exercise.links)) : []
    }
  } else {
    // Reset form when dialog closes
    if (form.value) {
      form.value.reset()
    }
    formData.value = {
      name: '',
      description: '',
      tagIds: [],
      links: []
    }
    originalExercise.value = null
  }
})

// Load tags and categories if needed
onMounted(async () => {
  // await exerciseStore.fetchExercise(props.exercise.id)
  if (tagsStore.tags.length === 0) {
    await tagsStore.fetchTags()
  }

  if (categoriesStore.categories.length === 0) {
    await categoriesStore.fetchCategories()
  }
})
</script>

<style scoped>
.paste-area {
  border: 2px dashed rgba(0, 0, 0, 0.12);
  border-radius: 8px;
  min-height: 120px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.paste-area:hover,
.paste-area.focused {
  border-color: rgb(var(--v-theme-primary));
  outline: none;
  background-color: rgba(var(--v-theme-primary), 0.05);
}
</style>
