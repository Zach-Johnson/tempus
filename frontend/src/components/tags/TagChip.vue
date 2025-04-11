<template>
  <v-chip
    :color="color"
    :variant="variant"
    :size="size"
    :closable="closable"
    @click="handleClick"
    @click:close="handleClose"
    class="tag-chip"
  >
    <template v-if="showIcon">
      <v-icon :icon="icon" size="small" class="mr-1"></v-icon>
    </template>
    {{ tag.name }}
  </v-chip>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  tag: {
    type: Object,
    required: true
  },
  color: {
    type: String,
    default: 'primary'
  },
  variant: {
    type: String,
    default: 'tonal' // Options: 'flat', 'outlined', 'tonal', etc.
  },
  size: {
    type: String,
    default: 'small' // Options: 'x-small', 'small', 'default', 'large', 'x-large'
  },
  clickable: {
    type: Boolean,
    default: false
  },
  closable: {
    type: Boolean,
    default: false
  },
  showIcon: {
    type: Boolean,
    default: true
  },
  icon: {
    type: String,
    default: 'mdi-tag'
  },
  disableNavigation: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['click', 'close'])
const router = useRouter()

// Computed properties
const isClickable = computed(() => {
  return props.clickable || (!props.disableNavigation && !props.closable)
})

// Methods
function handleClick() {
  emit('click', props.tag)
  
  // Navigate to tag's exercises view unless disabled
  if (isClickable.value && !props.disableNavigation) {
    router.push({ 
      name: 'exercises',
      query: { tag: props.tag.id }
    })
  }
}

function handleClose(event) {
  // Prevent click event from propagating
  event.stopPropagation()
  emit('close', props.tag)
}
</script>

<style scoped>
.tag-chip {
  margin-right: 4px;
  margin-bottom: 4px;
}
</style>
