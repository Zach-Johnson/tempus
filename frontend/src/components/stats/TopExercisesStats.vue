<template>
  <div class="top-exercises">
    <div v-if="loading" class="d-flex justify-center align-center py-4">
      <v-progress-circular indeterminate color="primary"></v-progress-circular>
    </div>
    
    <div v-else-if="topExercises.length === 0" class="text-center py-4">
      <p class="text-body-2 text-grey">No exercise data available</p>
    </div>
    
    <div v-else>
      <v-list density="compact">
        <v-list-subheader>Most Practiced Exercises</v-list-subheader>
        
        <v-list-item
          v-for="exercise in topExercises"
          :key="exercise.id"
          :to="{ name: 'exercise-detail', params: { id: exercise.id }}"
        >
          <template v-slot:prepend>
            <v-avatar size="32" color="primary" class="text-white">
              {{ exercise.name.charAt(0).toUpperCase() }}
            </v-avatar>
          </template>
          
          <v-list-item-title>{{ exercise.name }}</v-list-item-title>
          <v-list-item-subtitle>
            <v-progress-linear
              :model-value="(exercise.totalMinutes / maxExerciseTime) * 100"
              color="primary"
              height="8"
              rounded
              class="my-1"
            ></v-progress-linear>
            
            <div class="d-flex mt-1">
              <span v-if="exercise.practiceCount" class="text-caption text-grey mr-2">
                <v-icon size="x-small" color="grey">mdi-replay</v-icon>
                {{ exercise.practiceCount }} times
              </span>
              <span v-if="exercise.maxBpm" class="text-caption text-grey">
                <v-icon size="x-small" color="grey">mdi-speedometer</v-icon>
                {{ exercise.maxBpm }} BPM max
              </span>
            </div>
          </v-list-item-subtitle>
          
          <template v-slot:append>
            <strong>{{ formatTime(exercise.totalMinutes) }}</strong>
          </template>
        </v-list-item>
      </v-list>
      
      <div class="d-flex justify-end mt-2">
        <v-btn
          variant="text"
          size="small"
          to="/exercises"
        >
          View All Exercises
        </v-btn>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useAppStore } from '@/stores/app'

const props = defineProps({
  stats: {
    type: Object,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  },
  limit: {
    type: Number,
    default: 5
  }
})

const appStore = useAppStore()

// Extract and sort top exercises
const topExercises = computed(() => {
  if (!props.stats || !props.stats.exerciseStats) return []
  
  return Object.values(props.stats.exerciseStats)
    .sort((a, b) => b.totalTime - a.totalTime)
    .slice(0, props.limit)
    .map(exercise => ({
      ...exercise,
      totalMinutes: exercise.totalTime || 0
    }))
})

// Calculate the maximum exercise time for scaling the progress bars
const maxExerciseTime = computed(() => {
  if (topExercises.value.length === 0) return 0
  return Math.max(...topExercises.value.map(e => e.totalMinutes))
})

// Format time helper function
function formatTime(minutes) {
  return appStore.formatMinutes(minutes)
}
</script>
