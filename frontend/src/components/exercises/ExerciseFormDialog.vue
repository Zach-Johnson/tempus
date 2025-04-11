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
              <v-text-field
                v-model="formData.name"
                label="Exercise Name"
                :rules="nameRules"
                required
                variant="outlined"
                autofocus
              ></v-text-field>
            </v-col>
            
            <v-col cols="12">
              <v-textarea
                v-model="formData.description"
                label="Description"
                variant="outlined"
                rows="4"
                placeholder="Provide a description for this exercise (optional)"
              ></v-textarea>
            </v-col>
            
            <v-col cols="12" md="6">
              <category-select
                v-model="formData.categoryIds"
                label="Categories"
                multiple
                chips
              ></category-select>
            </v-col>
            
            <v-col cols="12" md="6">
              <v-select
                v-model="formData.tagIds"
                :items="tagsForSelect"
                label="Tags"
                multiple
                chips
                variant="outlined"
                :loading="tagsStore.loading"
              >
                <template v-slot:selection="{ item }">
                  <v-chip
                    :key="item.raw.value"
                    closable
                    @click:close="removeTag(item.raw.value)"
                  >
                    {{ item.title }}
                  </v-chip>
                </template>
                <template v-slot:append-inner>
                  <v-btn
                    icon
                    size="small"
                    variant="text"
                    class="ms-2"
                    @click.stop="openTagForm"
                  >
                    <v-icon>mdi-plus</v-icon>
                  </v-btn>
                </template>
              </v-select>
            </v-col>
            
            <!-- External Resources Section -->
            <v-col cols="12">
              <div class="d-flex align-center mb-2">
                <div class="text-body-1 font-weight-medium">External Resources</div>
                <v-btn
                  variant="text"
                  density="compact"
                  icon="mdi-plus"
                  size="small"
                  color="primary"
                  class="ml-2"
                  @click="addExternalLink"
                  title="Add external resource"
                ></v-btn>
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
                      <v-btn
                        icon
                        variant="text"
                        size="small"
                        color="error"
                        @click="removeLink(index)"
                      >
                        <v-icon>mdi-delete</v-icon>
                      </v-btn>
                    </template>
                  </v-list-item>
                </v-list>
              </div>
            </v-col>
          </v-row>
        </v-form>
      </v-card-text>
      
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          color="grey-darken-1"
          variant="text"
          @click="close"
        >
          Cancel
        </v-btn>
        <v-btn
          color="primary"
          variant="flat"
          @click="save"
          :disabled="!formValid || !formData.name"
          :loading="saving"
        >
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
          <v-text-field
            v-model="linkFormData.url"
            label="URL"
            :rules="urlRules"
            required
            variant="outlined"
            placeholder="https://example.com"
            class="mb-4"
          ></v-text-field>
          
          <v-text-field
            v-model="linkFormData.description"
            label="Description"
            variant="outlined"
            placeholder="Description of the resource (optional)"
          ></v-text-field>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="grey-darken-1" variant="text" @click="linkDialog = false">Cancel</v-btn>
        <v-btn 
          color="primary" 
          variant="flat" 
          @click="saveLink"
          :disabled="!linkFormValid"
        >
          Add
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
  
  <!-- Tag Form Dialog -->
  <tag-form-dialog
    v-model="tagFormDialog"
    :tag="{ name: '', category_ids: [] }"
    :is-edit="false"
    @save="onTagCreated"
  />
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useTagsStore } from '@/stores/tags.js'
import { useCategoriesStore } from '@/stores/categories.js'
import { useAppStore } from '@/stores/app.js'
import CategorySelect from '@/components/categories/CategorySelect.vue'
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
const appStore = useAppStore()

// Data
const form = ref(null)
const formValid = ref(false)
const saving = ref(false)
const formData = ref({
  name: '',
  description: '',
  categoryIds: [],
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

// Tag Form Dialog
const tagFormDialog = ref(false)

// Computed
const dialogModel = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const tagsForSelect = computed(() => {
  return tagsStore.tags.map(tag => ({
    title: tag.name,
    value: tag.id
  }))
})

// Validation rules
const nameRules = [
  v => !!v || 'Exercise name is required',
  v => (v && v.length >= 2) || 'Exercise name must be at least 2 characters',
  v => (v && v.length <= 50) || 'Exercise name must be less than 50 characters'
]

const urlRules = [
  v => !!v || 'URL is required',
  v => /^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w.-]*)*\/?$/.test(v) || 'Enter a valid URL'
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
      categoryIds: formData.value.categoryIds.map(id => Number(id)),
      tagIds: formData.value.tagIds.map(id => Number(id)),
      links: formData.value.links
    }
    
    // Emit save event with the form data
    emit('save', exerciseData)
  } catch (error) {
    console.error('Error saving exercise:', error)
  } finally {
    saving.value = false
  }
}

function removeTag(tagId) {
  formData.value.tagIds = formData.value.tagIds.filter(id => id !== tagId)
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

// Reset form when dialog opens/closes
watch(() => props.modelValue, (isOpen) => {
  if (isOpen && props.exercise) {
    // Clone the exercise to avoid modifying the original
    originalExercise.value = { ...props.exercise }
    
    // Set form values
    formData.value = {
      name: props.exercise.name || '',
      description: props.exercise.description || '',
      categoryIds: props.exercise.categoryIds ? [...props.exercise.categoryIds] : [],
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
      categoryIds: [],
      tagIds: [],
      links: []
    }
    originalExercise.value = null
  }
})

// Load tags and categories if needed
onMounted(async () => {
  if (tagsStore.tags.length === 0) {
    await tagsStore.fetchTags()
  }
  
  if (categoriesStore.categories.length === 0) {
    await categoriesStore.fetchCategories()
  }
})
</script>
