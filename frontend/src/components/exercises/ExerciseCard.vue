<template>
  <v-card
    class="exercise-card"
    :to="computedTo"
    :flat="flat"
    :variant="variant"
    @click="$emit('click', exercise)"
  >
    <v-card-title class="text-subtitle-1">
      {{ exercise.name }}
    </v-card-title>
    
    <v-card-text>
      <p v-if="exercise.description" class="text-body-2 mb-2">
        {{ truncateText(exercise.description, 100) }}
      </p>
      
      <div v-if="showCategories && exercise.categoryIds && exercise.categoryIds.length > 0" class="mb-2">
        <v-chip-group>
          <category-chip
            v-for="categoryId in exercise.categoryIds.slice(0, 3)"
            :key="categoryId"
            :category="getCategoryById(categoryId)"
            size="x-small"
          ></category-chip>
          <v-chip
            v-if="exercise.categoryIds.length > 3"
            size="x-small"
            color="secondary"
            variant="outlined"
          >
            +{{ exercise.categoryIds.length - 3 }}
          </v-chip>
        </v-chip-group>
      </div>
      
      <div v-if="showTags && exercise.tagIds && exercise.tagIds.length > 0">
        <v-chip-group>
          <tag-chip
            v-for="tagId in exercise.tagIds.slice(0, 3)"
            :key="tagId"
            :tag="getTagById(tagId)"
            size="x-small"
          ></tag-chip>
          <v-chip
            v-if="exercise.tagIds.length > 3"
            size="x-small"
            color="primary"
            variant="outlined"
          >
            +{{ exercise.tagIds.length - 3 }}
          </v-chip>
        </v-chip-group>
      </div>
    </v-card-text>
    
    <slot name="actions"></slot>
  </v-card>
</template>

<script setup>
import { computed } from 'vue'

// Computed props
const computedTo = computed(() => {
  if (props.to) return props.to
  return props.exercise && props.exercise.id 
    ? { name: 'exercise-detail', params: { id: props.exercise.id } }
    : null
})
import { useCategoriesStore } from '@/stores/categories.js'
import { useTagsStore } from '@/stores/tags.js'
import CategoryChip from '@/components/categories/CategoryChip.vue'
import TagChip from '@/components/tags/TagChip.vue'

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
  }
})

defineEmits(['click'])

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
}

.exercise-card:hover {
  transform: translateY(-2px);
}
</style>
