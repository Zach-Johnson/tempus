<template>
  <v-dialog v-model="dialogModel" max-width="600" persistent>
    <v-card>
      <v-card-title class="text-h5">
        {{ isEdit ? `Edit Category: ${originalCategory?.name}` : 'Create New Category' }}
      </v-card-title>
      
      <v-card-text>
        <v-form ref="form" @submit.prevent="save" v-model="formValid">
          <v-row>
            <v-col cols="12">
              <v-text-field
                v-model="formData.name"
                label="Category Name"
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
                placeholder="Provide a description for this category (optional)"
              ></v-textarea>
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
import { ref, computed, watch } from 'vue'

// Props
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  category: {
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

// Data
const form = ref(null)
const formValid = ref(false)
const saving = ref(false)
const formData = ref({
  name: '',
  description: ''
})
const originalCategory = ref(null)

// Computed
const dialogModel = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// Validation rules
const nameRules = [
  v => !!v || 'Category name is required',
  v => (v && v.length >= 2) || 'Category name must be at least 2 characters',
  v => (v && v.length <= 50) || 'Category name must be less than 50 characters'
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
    const categoryData = {
      name: formData.value.name,
      description: formData.value.description || ''
    }
    
    // Emit save event with the form data
    emit('save', categoryData)
  } catch (error) {
    console.error('Error saving category:', error)
  } finally {
    saving.value = false
  }
}

// Reset form when dialog opens/closes
watch(() => props.modelValue, (isOpen) => {
  if (isOpen && props.category) {
    // Clone the category to avoid modifying the original
    originalCategory.value = { ...props.category }
    
    // Set form values
    formData.value = {
      name: props.category.name || '',
      description: props.category.description || ''
    }
  } else {
    // Reset form when dialog closes
    if (form.value) {
      form.value.reset()
    }
    formData.value = {
      name: '',
      description: ''
    }
    originalCategory.value = null
  }
})
</script>
