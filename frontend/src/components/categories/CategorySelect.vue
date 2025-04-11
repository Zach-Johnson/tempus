<template>
  <div>
    <v-select
      :model-value="modelValue"
      @update:model-value="$emit('update:modelValue', $event)"
      :items="categoriesForSelect"
      :label="label"
      :placeholder="placeholder"
      :hint="hint"
      :rules="rules"
      :density="density"
      :variant="variant"
      :multiple="multiple"
      :chips="chips"
      :clearable="clearable"
      :disabled="disabled || categoriesStore.loading"
      :loading="categoriesStore.loading"
      hide-details
      :item-value="itemValue"
      :item-title="itemTitle"
      :required="required"
      :eager="true"
    >
      <template v-if="chips && multiple" v-slot:selection="{ item, index }">
        <v-chip
          v-if="index < 3"
          size="small"
          :color="chipColor"
          :variant="chipVariant"
        >
          {{ item.raw.title }}
        </v-chip>
        <span
          v-if="index === 3"
          class="text-caption text-grey ms-2"
        >
          (+{{ modelValue.length - 3 }} more)
        </span>
      </template>
      
      <template v-if="categoriesStore.loading && categoriesForSelect.length === 0" v-slot:no-data>
        <div class="pa-2">Loading categories...</div>
      </template>
      
      <template v-else-if="categoriesForSelect.length === 0" v-slot:no-data>
        <div class="pa-2">
          No categories available
          <v-btn
            variant="text"
            size="small"
            color="primary"
            class="ms-2"
            @click="openCreateDialog"
          >
            Create Category
          </v-btn>
        </div>
      </template>
    </v-select>
    
    <!-- Create Category Dialog -->
    <category-form-dialog
      v-model="dialogVisible"
      :category="{ name: '', description: '' }"
      :is-edit="false"
      @save="createCategory"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useCategoriesStore } from '@/stores/categories.js'
import { useAppStore } from '@/stores/app.js'
import CategoryFormDialog from './CategoryFormDialog.vue'

const props = defineProps({
  modelValue: {
    type: [Array, Number, String],
    default: () => []
  },
  label: {
    type: String,
    default: 'Categories'
  },
  placeholder: {
    type: String,
    default: 'Select categories'
  },
  hint: {
    type: String,
    default: ''
  },
  rules: {
    type: Array,
    default: () => []
  },
  density: {
    type: String,
    default: 'default'
  },
  variant: {
    type: String,
    default: 'outlined'
  },
  multiple: {
    type: Boolean,
    default: false
  },
  chips: {
    type: Boolean,
    default: false
  },
  clearable: {
    type: Boolean,
    default: true
  },
  disabled: {
    type: Boolean,
    default: false
  },
  itemValue: {
    type: String,
    default: 'value'
  },
  itemTitle: {
    type: String,
    default: 'title'
  },
  required: {
    type: Boolean,
    default: false
  },
  chipColor: {
    type: String,
    default: 'primary'
  },
  chipVariant: {
    type: String,
    default: 'tonal'
  }
})

const emit = defineEmits(['update:modelValue', 'category-created'])

const categoriesStore = useCategoriesStore()
const appStore = useAppStore()
const dialogVisible = ref(false)

// Computed
const categoriesForSelect = computed(() => {
  return categoriesStore.categories.map(category => ({
    title: category.name,
    value: category.id,
    description: category.description
  }))
})

// Methods
function openCreateDialog() {
  dialogVisible.value = true
}

async function createCategory(categoryData) {
  try {
    const newCategory = await categoriesStore.createCategory(categoryData)
    appStore.showSuccessMessage(`Category "${categoryData.name}" created successfully`)
    dialogVisible.value = false
    
    // Update selected value if multiple
    if (props.multiple) {
      const currentValue = Array.isArray(props.modelValue) ? [...props.modelValue] : []
      currentValue.push(newCategory.id)
      emit('update:modelValue', currentValue)
    } else {
      emit('update:modelValue', newCategory.id)
    }
    
    emit('category-created', newCategory)
  } catch (error) {
    appStore.showErrorMessage(`Error creating category: ${error.message}`)
  }
}

// Load categories if needed
onMounted(async () => {
  if (categoriesStore.categories.length === 0) {
    await categoriesStore.fetchCategories()
  }
})
</script>
