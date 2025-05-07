<template>
  <div style="position: relative; height: 100%;">
    <div v-if="loading" class="d-flex justify-center align-center" style="height: 100%">
      <v-progress-circular indeterminate color="primary"></v-progress-circular>
    </div>
    <div v-else-if="chartData.datasets.length > 0" style="height: 100%">
      <Bar :data="chartData" :options="chartOptions" />
    </div>
    <div v-else class="d-flex justify-center align-center" style="height: 100%">
      <p class="text-body-1">No practice data available for the last 30 days</p>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
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

// For debugging
const debug = ref(false)

// Generate colors for categories
const categoryColors = computed(() => {
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

  const colors = {
    0: '#9e9e9e' // Fallback for uncategorized
  }

  categoriesStore.categories.forEach((category, index) => {
    colors[category.id] = colorPalette[index % colorPalette.length]
  })

  return colors
})

// Computed properties for chart data
const chartData = computed(() => {
  if (debug.value) {
    console.log('practiceFrequency:', statsStore.practiceFrequency);
    console.log('categoryDistribution:', statsStore.categoryDistribution);
  }

  if (!statsStore.practiceFrequency || statsStore.practiceFrequency.length === 0) {
    return { labels: [], datasets: [] };
  }

  // Get unique dates from practice frequency data
  const dates = [...new Set(statsStore.practiceFrequency.map(item => item.date))];

  // Sort the dates chronologically
  const sortedDates = dates.sort((a, b) => new Date(a) - new Date(b));

  const monthNames = [
    'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
    'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'
  ];

  // Format the dates for display as labels
  const labels = sortedDates.map(dateStr => {
    const date = new Date(dateStr);
    return `${monthNames[date.getUTCMonth()]} ${date.getUTCDate()}`;
  });

  // Create a map of category ID to dataset
  const datasetsMap = {};

  // Initialize datasets for each category from categoryDistribution
  if (statsStore.categoryDistribution && statsStore.categoryDistribution.length > 0) {
    statsStore.categoryDistribution.forEach(category => {
      datasetsMap[category.categoryId] = {
        label: category.categoryName,
        backgroundColor: categoryColors.value[category.categoryId] || categoryColors.value[0],
        data: Array(labels.length).fill(0) // Initialize with zeros
      };
    });
  } else {
    // Fallback to a single dataset if no category data
    datasetsMap['total'] = {
      label: 'Practice Time',
      backgroundColor: '#1976D2',
      data: Array(labels.length).fill(0)
    };
  }

  // Add data points from practiceFrequency
  statsStore.practiceFrequency.forEach(point => {
    const dateStr = point.date;
    const dateIndex = sortedDates.indexOf(dateStr);

    if (dateIndex !== -1) {
      const minutes = point.minutes;

      if (point.categoryId !== null && datasetsMap[point.categoryId]) {
        // Add category-specific data
        datasetsMap[point.categoryId].data[dateIndex] += minutes;
      } else if (datasetsMap['total']) {
        // Add to total if no category is specified and we're using the fallback
        datasetsMap['total'].data[dateIndex] += minutes;
      }
    }
  });

  // Convert map to array, filtering out empty datasets
  const datasets = Object.values(datasetsMap)
    .filter(dataset => dataset.data.some(value => value > 0));

  return {
    labels,
    datasets
  };
});

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
        label: function (context) {
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
</script>
