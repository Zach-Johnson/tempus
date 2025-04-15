<template>
  <div style="position: relative; height: 100%;">
    <div v-if="loading" class="d-flex justify-center align-center" style="height: 100%">
      <v-progress-circular indeterminate color="primary"></v-progress-circular>
    </div>
    <div v-else-if="chartData.datasets.length > 0" style="height: 100%">
      <Bar
        :data="chartData"
        :options="chartOptions"
      />
    </div>
    <div v-else class="d-flex justify-center align-center" style="height: 100%">
      <p class="text-body-1">No practice data available for the last 30 days</p>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, onMounted, watch } from 'vue'
import { Bar } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale } from 'chart.js'
import { useCategoriesStore } from '@/stores/categories.js'
import { useStatsStore } from '@/stores/stats.js'
import { useAppStore } from '@/stores/app.js'

// Register ChartJS components
ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)

// Props
const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  }
})

const categoriesStore = useCategoriesStore()
const statsStore = useStatsStore()
const appStore = useAppStore()

// For debugging
const debug = ref(true)

// Generate colors for categories
const categoryColors = {
  0: '#9e9e9e', // Gray for uncategorized items (category_id 0)
}

// Prepare category colors
const setupCategoryColors = () => {
  // Generate colors for each category
  const colorPalette = [
    '#1976D2', // Primary blue
    '#E53935', // Red
    '#43A047', // Green
    '#FB8C00', // Orange
    '#8E24AA', // Purple
    '#00ACC1', // Cyan
    '#FFB300', // Amber
    '#5E35B1', // Deep Purple
    '#1E88E5', // Blue
    '#00897B', // Teal
  ]

  categoriesStore.categories.forEach((category, index) => {
    categoryColors[category.id] = colorPalette[index % colorPalette.length]
  })
}

// Computed properties for chart data
const chartData = computed(() => {
  if (debug.value) {
    console.log('practiceFrequency:', statsStore.practiceFrequency)
    console.log('categoryDistribution:', statsStore.categoryDistribution)
  }

  if (!statsStore.practiceFrequency || statsStore.practiceFrequency.length === 0) {
    return { labels: [], datasets: [] }
  }

  // Get unique dates from practice frequency data
  const sortedFrequency = [...statsStore.practiceFrequency].sort((a, b) => 
    new Date(a.date) - new Date(b.date)
  )
  
  // Get the unique dates as labels
  const labels = [...new Set(sortedFrequency.map(item => {
    // Format the date as 'MMM D' (e.g., 'Jan 1')
    const date = new Date(item.date)
    return `${date.toLocaleString('default', { month: 'short' })} ${date.getDate()}`
  }))]

  // Check if we have category distribution data
  if (!statsStore.categoryDistribution || statsStore.categoryDistribution.length === 0) {
    // Simplified approach if no category data: just show total practice time
    const dataset = {
      label: 'Practice Time',
      backgroundColor: '#1976D2',
      data: labels.map(() => 0)
    }

    // Fill in the data
    sortedFrequency.forEach(item => {
      const dateStr = new Date(item.date).toLocaleString('default', { month: 'short' }) + 
                     ' ' + new Date(item.date).getDate()
      const dateIndex = labels.indexOf(dateStr)
      
      if (dateIndex !== -1) {
        dataset.data[dateIndex] += item.minutes
      }
    })

    return {
      labels,
      datasets: [dataset]
    }
  }

  // If we have category distribution data, proceed with full implementation
  
  // Create datasets for each category
  const categoriesMap = {}
  
  // Initialize with zeros for all dates
  statsStore.categoryDistribution.forEach(category => {
    categoriesMap[category.category_id] = {
      label: category.category_name,
      backgroundColor: categoryColors[category.category_id] || categoryColors[0],
      data: labels.map(() => 0) // Initialize with zeros for all dates
    }
  })
  
  // Add uncategorized category if needed
  if (!categoriesMap[0] && categoryColors[0]) {
    categoriesMap[0] = {
      label: 'Uncategorized',
      backgroundColor: categoryColors[0],
      data: labels.map(() => 0)
    }
  }
  
  // Fill in the data
  sortedFrequency.forEach((item) => {
    const dateStr = new Date(item.date).toLocaleString('default', { month: 'short' }) + 
                   ' ' + new Date(item.date).getDate()
    const dateIndex = labels.indexOf(dateStr)
    
    if (dateIndex === -1) return
    
    // If we don't have category breakdown in practiceFrequency,
    // proportionally distribute based on category distribution percentages
    const minutes = item.minutes
    
    // Only distribute if we have non-zero minutes
    if (minutes > 0) {
      statsStore.categoryDistribution.forEach(category => {
        if (categoriesMap[category.category_id]) {
          // Use percentage (if available) or distribute evenly
          const percentage = category.percentage || 
                           (100 / statsStore.categoryDistribution.length)
          const share = minutes * (percentage / 100)
          categoriesMap[category.category_id].data[dateIndex] += share
        }
      })
    }
  })
  
  // Convert to datasets array, filtering out empty categories
  const datasets = Object.values(categoriesMap)
    .filter(category => category.data.some(value => value > 0))
  
  return {
    labels,
    datasets
  }
})

// Chart options
const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: true,
      position: 'top'
    },
    tooltip: {
      callbacks: {
        label: function(context) {
          const value = context.raw;
          if (value < 60) {
            return `${context.dataset.label}: ${Math.round(value)} minutes`;
          } else {
            const hours = Math.floor(value / 60);
            const minutes = Math.round(value % 60);
            if (minutes === 0) {
              return `${context.dataset.label}: ${hours} hours`;
            } else {
              return `${context.dataset.label}: ${hours}h ${minutes}min`;
            }
          }
        }
      }
    }
  },
  scales: {
    x: {
      stacked: true,
      ticks: {
        maxRotation: 45,
        minRotation: 45
      }
    },
    y: {
      stacked: true,
      beginAtZero: true,
      title: {
        display: true,
        text: 'Minutes'
      }
    }
  }
}))

// Watch for changes in stats data for debugging
watch(() => statsStore.practiceStats, (newStats) => {
  if (debug.value && newStats) {
    console.log('New practice stats received:', newStats)
  }
}, { deep: true })

// Setup and initialize
onMounted(() => {
  setupCategoryColors()
  
  // Enable this line to debug
  // debug.value = true
  
  if (debug.value) {
    // Log initial data
    console.log('Initial store data:')
    console.log('Categories:', categoriesStore.categories)
    console.log('Practice stats:', statsStore.practiceStats)
    console.log('Practice frequency:', statsStore.practiceFrequency)
    console.log('Category distribution:', statsStore.categoryDistribution)
  }
})
</script>
