<template>
  <div class="tag-list">
    <div v-if="loading" class="d-flex align-center">
      <v-progress-circular indeterminate size="20" width="2" color="primary" class="mr-2"></v-progress-circular>
      <span class="text-body-2">Loading tags...</span>
    </div>
    
    <div v-else-if="!tags || tags.length === 0" class="text-body-2 text-grey">
      {{ noTagsMessage }}
    </div>
    
    <div v-else class="tag-list-container">
      <tag-chip
        v-for="tag in tags"
        :key="tag.id"
        :tag="tag"
        :color="color"
        :variant="variant"
        :size="size"
        :closable="selectable && modelValue && modelValue.includes(tag.id)"
        :clickable="selectable"
        :disable-navigation="disableNavigation"
        @click="toggleTag(tag)"
        @close="removeTag(tag)"
      ></tag-chip>
      
      <div v-if="showAddButton && selectable" class="d-inline-block">
        <v-btn
          variant="outlined"
          :size="size"
          icon
          class="tag-add-btn"
          @click="openTagSelect"
        >
          <v-icon>mdi-plus</v-icon>
        </v-btn>
      </div>
    </div>
    
    <!-- Tag Selection Dialog -->
    <v-dialog v-model="tagSelectDialog" max-width="500">
      <v-card>
        <v-card-title>Add Tags</v-card-title>
        <v-card-text>
          <v-text-field
            v-model="tagSearch"
            label="Search tags"
            variant="outlined"
            density="compact"
            prepend-inner-icon="mdi-magnify"
            clearable
            class="mb-4"
          ></v-text-field>
          
          <v-list v-if="availableTags.length > 0" density="compact">
            <v-list-item
              v-for="tag in filteredAvailableTags"
              :key="tag.id"
              :value="tag.id"
              @click="addTag(tag)"
            >
              <template v-slot:prepend>
                <v-checkbox-btn
                  :model-value="isTagSelected(tag)"
                  color="primary"
                ></v-checkbox-btn>
              </template>
              <v-list-item-title>{{ tag.name }}</v-list-item-title>
            </v-list-item>
          </v-list>
          
          <div v-else class="text-center py-4">
            <p class="text-body-2 text-grey">No tags available</p>
          </div>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn 
            color="primary" 
            variant="text" 
            @click="tagSelectDialog = false"
          >
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useTagsStore } from '@/stores/tags.js'
import TagChip from './TagChip.vue'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  },
  tags: {
    type: Array,
    default: () => []
  },
  categoryId: {
    type: [Number, String],
    default: null
  },
  color: {
    type: String,
    default: 'primary'
  },
  variant: {
    type: String,
    default: 'tonal'
  },
  size: {
    type: String,
    default: 'small'
  },
  selectable: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  noTagsMessage: {
    type: String,
    default: 'No tags'
  },
  showAddButton: {
    type: Boolean,
    default: true
  },
  disableNavigation: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'add', 'remove'])
const tagsStore = useTagsStore()

// Data
const tagSelectDialog = ref(false)
const tagSearch = ref('')

// Computed
const availableTags = computed(() => {
  // If we have a category filter, get tags for that category
  if (props.categoryId) {
    return tagsStore.tagsByCategory(parseInt(props.categoryId))
  }
  
  // Otherwise, get all tags
  return tagsStore.tagsSorted
})

const filteredAvailableTags = computed(() => {
  if (!tagSearch.value) return availableTags.value
  
  const search = tagSearch.value.toLowerCase()
  return availableTags.value.filter(tag => 
    tag.name.toLowerCase().includes(search)
  )
})

const selectedTags = computed(() => {
  if (!props.modelValue) return []
  
  // Convert tag IDs to tag objects
  return props.modelValue
    .map(tagId => {
      const tag = tagsStore.tagById(tagId)
      return tag || null
    })
    .filter(tag => tag !== null)
})

// Methods
function isTagSelected(tag) {
  return props.modelValue && props.modelValue.includes(tag.id)
}

function openTagSelect() {
  tagSearch.value = ''
  tagSelectDialog.value = true
}

function toggleTag(tag) {
  if (!props.selectable) return
  
  if (isTagSelected(tag)) {
    removeTag(tag)
  } else {
    addTag(tag)
  }
}

function addTag(tag) {
  if (!props.selectable) return
  
  // Skip if already selected
  if (isTagSelected(tag)) return
  
  const newValue = [...(props.modelValue || []), tag.id]
  emit('update:modelValue', newValue)
  emit('add', tag)
}

function removeTag(tag) {
  if (!props.selectable) return
  
  const newValue = (props.modelValue || []).filter(id => id !== tag.id)
  emit('update:modelValue', newValue)
  emit('remove', tag)
}

// Fetch tags if needed
onMounted(async () => {
  if (tagsStore.tags.length === 0) {
    await tagsStore.fetchTags()
  }
})
</script>

<style scoped>
.tag-list-container {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
}

.tag-add-btn {
  margin-right: 4px;
  margin-bottom: 4px;
}
</style>
