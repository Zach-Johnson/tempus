<template>
  <div style="position: relative; height: 100%;">
    <Bar
      :data="chartData"
      :options="chartOptions"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { Bar } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale } from 'chart.js'

// Register ChartJS components
ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)

// Props
const props = defineProps({
  chartData: {
    type: Array,
    required: true
  }
})

// Computed properties for chart data
const chartData = computed(() => {
  return {
    labels: props.chartData.map(item => item.date),
    datasets: [
      {
        label: 'Practice Minutes',
        data: props.chartData.map(item => item.minutes),
        backgroundColor: '#1976D2',
        borderColor: '#1565C0',
        borderWidth: 1
      }
    ]
  }
})

// Chart options
const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: false
    },
    tooltip: {
      callbacks: {
        label: function(context) {
          const value = context.raw;
          if (value < 60) {
            return `${value} minutes`;
          } else {
            const hours = Math.floor(value / 60);
            const minutes = value % 60;
            if (minutes === 0) {
              return `${hours} hours`;
            } else {
              return `${hours}h ${minutes}min`;
            }
          }
        }
      }
    }
  },
  scales: {
    y: {
      beginAtZero: true,
      title: {
        display: true,
        text: 'Minutes'
      }
    },
    x: {
      ticks: {
        maxRotation: 45,
        minRotation: 45
      }
    }
  }
}
</script>
