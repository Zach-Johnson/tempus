<template>
  <v-dialog v-model="dialogModel" max-width="600" persistent>
    <v-card>
      <v-card-title class="text-h5">
        {{ isEdit ? `Edit Tag: ${originalTag?.name}` : 'Create New Tag' }}
      </v-card-title>
      
      <v-card-text>
        <v-form ref="form" @submit.prevent="save" v-model="formValid">
          <v-row>
            <v-col cols="12">
              <v-text-field
                v-model="formData.name"
                label="Tag Name"
                :rules="nameRules"
                required
                variant="outlined"
                autofocus
              ></v-text-field>
            </v-col>
            
            <v-col cols="12">
              <v-select
                v-model="formData.category_ids"
                :items="categoriesForSelect"
                label="Categories"
                multiple
                chips
                variant="outlined"
                :loading="categoriesStore.loading"
                hint="Select the categories this tag belongs to"
                persistent-hint
              >
                <template v-slot:selection="{ item }">
                  <v-chip
                    :key="item.raw.value"
                    :color="getCategoryColor(item.raw.value)"
                    closable
                    @click:close="removeCategory(item.raw.value)"
                  >
                    {{ item.title }}
                  </v-chip>
                </template>
                <template v-slot:item="{ item, props }">
                  <v-list-item v-bind="props">
                    <template v-slot:prepend>
                      <v-checkbox-btn :model-value="isCategorySelected(item.raw.value)"></v-checkbox-btn>
                    </template>
                    <template v-slot:title>
                      <span :style="`color: ${getCategoryColor(item.raw.value)}`">{{ item.raw.title }}</span>
                    </template>
                  </v-list-item>
                </template>
              </v-select>
              <div class="text-caption text-grey">
                Assign this tag to one or more categories to help organize your practice
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
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useCategoriesStore } from '@/stores/categories.js'

// Props
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  tag: {
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

// Store
const categoriesStore = useCategoriesStore()

// Data
const form = ref(null)
const formValid = ref(false)
const saving = ref(false)
const formData = ref({
  name: '',
  category_ids: []
})
const originalTag = ref(null)

// Category colors
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

const categoriesForSelect = computed(() => {
  return categoriesStore.categories.map(category => ({
    title: category.name,
    value: category.id
  }))
})

// Validation rules
const nameRules = [
  v => !!v || 'Tag name is required',
  v => (v && v.length >= 2) || 'Tag name must be at least 2 characters',
  v => (v && v.length <= 30) || 'Tag name must be less than 30 characters'
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
    const tagData = {
      name: formData.value.name,
      category_ids: formData.value.category_ids.map(id => Number(id))
    }
    
    // Emit save event with the form data
    emit('save', tagData)
  } catch (error) {
    console.error('Error saving tag:', error)
  } finally {
    saving.value = false
  }
}

function removeCategory(categoryId) {
  formData.value.category_ids = formData.value.category_ids.filter(id => id !== categoryId)
}

function isCategorySelected(categoryId) {
  return formData.value.category_ids.includes(categoryId)
}

function getCategoryColor(categoryId) {
  return categoryColors[categoryId] || '#9e9e9e'
}

// Reset form when dialog opens/closes
watch(() => props.modelValue, (isOpen) => {
  if (isOpen && props.tag) {
    // Clone the tag to avoid modifying the original
    originalTag.value = { ...props.tag }
    
    // Set form values
    formData.value = {
      name: props.tag.name || '',
      category_ids: props.tag.category_ids ? [...props.tag.category_ids] : []
    }
  } else {
    // Reset form when dialog closes
    if (form.value) {
      form.value.reset()
    }
    formData.value = {
      name: '',
      category_ids: []
    }
    originalTag.value = null
  }
})

// Load categories if not already loaded
onMounted(async () => {
  if (categoriesStore.categories.length === 0) {
    await categoriesStore.fetchCategories()
  }
})
</script>
