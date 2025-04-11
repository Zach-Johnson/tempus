<template>
  <div class="exercises-view">
    <div class="container-fluid">
      <div class="page-header mb-4">
        <h1>Drum Exercises</h1>
        <p class="text-muted">Browse and manage your collection of drum exercises.</p>
      </div>
      
      <!-- Show the list by default or when no route params -->
      <ExerciseList v-if="!$route.params.id" />
      
      <!-- Router outlet for nested routes like exercise details -->
      <router-view v-else />
    </div>
  </div>
</template>

<script>
import ExerciseList from '@/components/exercises/ExerciseList.vue';

export default {
  name: 'ExercisesView',
  
  components: {
    ExerciseList
  }
};
</script>

<style scoped>
.page-header h1 {
  font-weight: 600;
  color: #212529;
}
</style>
