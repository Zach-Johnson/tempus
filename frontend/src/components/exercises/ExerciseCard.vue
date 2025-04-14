<template>
  <v-card
    class="exercise-card"
    :to="selectable ? undefined : computedTo"
    :flat="flat"
    :variant="variant"
    @click="handleClick"
    :class="{ 'selected-exercise': selected }"
  >
    <v-card-title class="text-subtitle-1">
      {{ exercise.name }}
    </v-card-title>
    
    <v-card-text>
      <p v-if="exercise.description" class="text-body-2 mb-2">
        {{ truncateText(exercise.description, 100) }}
      </p>
      
      <div v-if="showCategories && exercise.categoryIds && exercise.categoryIds.length > 0" class="mb-2">
        <span class="text-caption text-grey">
          {{ exercise.categoryIds.length }} {{ exercise.categoryIds.length === 1 ? 'category' : 'categories' }}
        </span>
      </div>
      
      <div v-if="showTags && exercise.tagIds && exercise.tagIds.length > 0">
        <span class="text-caption text-grey">
          {{ exercise.tagIds.length }} {{ exercise.tagIds.length === 1 ? 'tag' : 'tags' }}
        </span>
      </div>
    </v-card-text>
    
    <slot name="actions">
      <v-card-actions v-if="selectable">
        <v-spacer></v-spacer>
        <v-btn 
          variant="text"
          :color="selected ? 'error' : 'primary'"
          size="small"
        >
          {{ selected ? 'Remove' : 'Add' }}
        </v-btn>
      </v-card-actions>
    </slot>
  </v-card>
</template>

<script setup>
import { computed } from 'vue'
import { useCategoriesStore } from '@/stores/categories.js'
import { useTagsStore } from '@/stores/tags.js'

const props = defineProps({
  exercise: {
    type: Object,
    required: true
  },
  to: {
    type: [String, Object],
    default: null
  },
  flat: {
    type: Boolean,
    default: false
  },
  variant: {
    type: String,
    default: 'elevated' // elevated, flat, tonal, outlined, etc.
  },
  showCategories: {
    type: Boolean,
    default: true
  },
  showTags: {
    type: Boolean,
    default: true
  },
  selected: {
    type: Boolean,
    default: false
  },
  selectable: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['click'])

// Computed props
const computedTo = computed(() => {
  if (props.to) return props.to
  return props.exercise && props.exercise.id 
    ? { name: 'exercise-detail', params: { id: props.exercise.id } }
    : null
})

// Stores
const categoriesStore = useCategoriesStore()
const tagsStore = useTagsStore()

// Methods
function truncateText(text, maxLength) {
  if (!text) return ''
  return text.length > maxLength 
    ? text.substring(0, maxLength) + '...' 
    : text
}

function handleClick() {
  emit('click', props.exercise)
}

// These functions are no longer used in the simplified template
function getCategoryById(categoryId) {
  return categoriesStore.categoryById(categoryId)
}

function getTagById(tagId) {
  return tagsStore.tagById(tagId)
}
</script>

<style scoped>
.exercise-card {
  height: 100%;
  transition: transform 0.2s;
  cursor: pointer;
}

.exercise-card:hover {
  transform: translateY(-2px);
}

.selected-exercise {
  border: 2px solid rgb(var(--v-theme-primary)) !important;
  background-color: rgba(var(--v-theme-primary), 0.05);
}
</style>
