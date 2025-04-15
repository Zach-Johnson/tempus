<template>
  <div style="position: relative; height: 100%;">
    <div v-if="loading" class="d-flex justify-center align-center" style="height: 100%">
      <v-progress-circular indeterminate color="primary"></v-progress-circular>
    </div>
    <Bar
      v-else-if="chartData.datasets.length > 0"
      :data="chartData"
      :options="chartOptions"
    />
    <div v-else class="d-flex justify-center align-center" style="height: 100%">
      <p class="text-body-1">No practice data available for the last 30 days</p>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import { Bar } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale } from 'chart.js'
import { useCategoriesStore } from '@/stores/categories.js'
import { useStatsStore } from '@/stores/stats.js'

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

// Generate colors for categories
const categoryColors = {
  0: '#9e9e9e', // Gray for uncategorized items (category_id 0)
}

// Prepare category colors when store loads
onMounted(() => {
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
})

// Computed properties for chart data
const chartData = computed(() => {
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

  // Get unique category IDs
  const uniqueCategoryIds = new Set()
  statsStore.categoryDistribution.forEach(cat => uniqueCategoryIds.add(cat.category_id))
  
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
  if (!categoriesMap[0]) {
    categoriesMap[0] = {
      label: 'Uncategorized',
      backgroundColor: categoryColors[0],
      data: labels.map(() => 0)
    }
  }
  
  // Fill in the data (assuming we have daily data with category breakdown)
  // Note: If practiceFrequency doesn't have category breakdowns, 
  // we'll need a different approach or backend API update
  sortedFrequency.forEach((item, index) => {
    const dateStr = new Date(item.date).toLocaleString('default', { month: 'short' }) + 
                   ' ' + new Date(item.date).getDate()
    const dateIndex = labels.indexOf(dateStr)
    
    if (dateIndex === -1) return
    
    // If we have category info
    if (item.categoryId) {
      if (categoriesMap[item.categoryId]) {
        categoriesMap[item.categoryId].data[dateIndex] += item.minutes
      }
    } else {
      // If we don't have category breakdown, proportionally distribute
      // based on overall category distribution percentages
      const minutes = item.minutes
      statsStore.categoryDistribution.forEach(category => {
        if (categoriesMap[category.category_id]) {
          const share = minutes * (category.percentage / 100)
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
const chartOptions = {
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
}
</script>
