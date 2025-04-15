<template>
  <div class="category-summary">
    <div v-if="loading" class="d-flex justify-center align-center py-4">
      <v-progress-circular indeterminate color="primary"></v-progress-circular>
    </div>
    
    <div v-else-if="categories.length === 0" class="text-center py-4">
      <p class="text-body-2 text-grey">No categories found</p>
    </div>
    
    <div v-else>
      <v-list density="compact">
        <v-list-subheader>Practice Time by Category</v-list-subheader>
        
        <v-list-item
          v-for="category in sortedCategories"
          :key="category.id"
          :to="{ name: 'category-detail', params: { id: category.id }}"
        >
          <template v-slot:prepend>
            <v-avatar size="32" :color="getCategoryColor(category)" class="text-white">
              {{ category.name.charAt(0).toUpperCase() }}
            </v-avatar>
          </template>
          
          <v-list-item-title>{{ category.name }}</v-list-item-title>
          <v-list-item-subtitle>
            <v-progress-linear
              :model-value="(category.totalMinutes / maxCategoryTime) * 100"
              :color="getCategoryColor(category)"
              height="8"
              rounded
              class="my-1"
            ></v-progress-linear>
          </v-list-item-subtitle>
          
          <template v-slot:append>
            <strong>{{ formatTime(category.totalMinutes) }}</strong>
          </template>
        </v-list-item>
        
        <!-- Uncategorized practice time if any -->
        <v-list-item v-if="uncategorizedTime > 0">
          <template v-slot:prepend>
            <v-avatar size="32" color="grey" class="text-white">
              ?
            </v-avatar>
          </template>
          
          <v-list-item-title>Uncategorized</v-list-item-title>
          <v-list-item-subtitle>
            <v-progress-linear
              :model-value="(uncategorizedTime / maxCategoryTime) * 100"
              color="grey"
              height="8"
              rounded
              class="my-1"
            ></v-progress-linear>
          </v-list-item-subtitle>
          
          <template v-slot:append>
            <strong>{{ formatTime(uncategorizedTime) }}</strong>
          </template>
        </v-list-item>
      </v-list>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useCategoriesStore } from '@/stores/categories'
import { useAppStore } from '@/stores/app'

const props = defineProps({
  stats: {
    type: Object,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const categoriesStore = useCategoriesStore()
const appStore = useAppStore()

// Generate colors for categories
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

// Calculate category statistics from the provided stats
const categories = computed(() => {
  if (!props.stats || !props.stats.categoryStats) return []
  
  return categoriesStore.categories.map(category => {
    const categoryData = props.stats.categoryStats[category.id] || { totalTime: 0 }
    return {
      ...category,
      totalMinutes: categoryData.totalTime || 0
    }
  }).filter(category => category.totalMinutes > 0)
})

// Sort categories by total time (descending)
const sortedCategories = computed(() => {
  return [...categories.value].sort((a, b) => b.totalMinutes - a.totalMinutes)
})

// Calculate maximum category time for scaling progress bars
const maxCategoryTime = computed(() => {
  if (categories.value.length === 0) return 0
  return Math.max(...categories.value.map(c => c.totalMinutes), uncategorizedTime.value)
})

// Calculate uncategorized time if any
const uncategorizedTime = computed(() => {
  if (!props.stats || !props.stats.categoryStats) return 0
  return props.stats.categoryStats.uncategorized?.totalTime || 0
})

// Get color for a category
function getCategoryColor(category) {
  const colorIndex = (category.id % Object.keys(categoryColors).length) + 1
  return categoryColors[colorIndex] || '#9e9e9e'
}

// Format time helper function
function formatTime(minutes) {
  return appStore.formatMinutes(minutes)
}
</script>
