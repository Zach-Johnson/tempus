<template>
  <div>
    <!-- Loading state -->
    <div v-if="loading" class="d-flex justify-center my-8">
      <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
    </div>

    <!-- Error state -->
    <v-alert v-else-if="error" type="error" class="my-4">
      {{ error }}
    </v-alert>

    <!-- Exercise details -->
    <template v-else>
      <div v-if="!exercisesStore.currentExercise">
        <v-alert type="info" class="my-4">
          Exercise not found. It may have been deleted.
          <div class="mt-4">
            <v-btn color="primary" to="/exercises">
              Back to Exercises
            </v-btn>
          </div>
        </v-alert>
      </div>
      <template v-else>
        <!-- Header section with actions -->
        <v-row class="mb-4">
          <v-col cols="12" sm="8">
            <v-btn
              variant="text"
              to="/exercises"
              prepend-icon="mdi-arrow-left"
              class="mb-2"
            >
              Back to Exercises
            </v-btn>
            <h1 class="text-h3">{{ exercisesStore.currentExercise.name }}</h1>
          </v-col>
          <v-col cols="12" sm="4" class="d-flex justify-end align-center">
            <v-btn
              color="primary"
              variant="outlined"
              class="mr-2"
              prepend-icon="mdi-pencil"
              @click="openEditDialog"
            >
              Edit
            </v-btn>
            <v-btn
              color="primary"
              variant="flat"
              prepend-icon="mdi-play"
              @click="startPractice"
            >
              Practice Now
            </v-btn>
          </v-col>
        </v-row>

        <!-- Exercise details card -->
        <v-card class="mb-6">
          <v-card-text>
            <v-row>
              <v-col cols="12" md="6">
                <div class="text-body-1 font-weight-medium mb-1">Description</div>
                <div class="text-body-2 mb-4">{{ exercisesStore.currentExercise.description || 'No description provided' }}</div>
                
                <div class="text-body-1 font-weight-medium mb-1">Categories</div>
                <div v-if="exerciseCategories.length === 0" class="text-body-2 mb-4 text-grey">
                  No categories assigned
                </div>
                <div v-else class="mb-4">
                  <category-list :categories="exerciseCategories" />
                </div>
                
                <div class="text-body-1 font-weight-medium mb-1">Tags</div>
                <div v-if="exerciseTags.length === 0" class="text-body-2 mb-4 text-grey">
                  No tags assigned
                </div>
                <div v-else class="mb-4">
                  <tag-list :tags="exerciseTags" />
                </div>
                
                <div class="text-body-1 font-weight-medium mb-1">Created</div>
                <div class="text-body-2 mb-4">{{ formatDate(exercisesStore.currentExercise.createdAt) }}</div>
                
                <div class="text-body-1 font-weight-medium mb-1">Last Updated</div>
                <div class="text-body-2">{{ formatDate(exercisesStore.currentExercise.updatedAt) }}</div>
              </v-col>
              
              <v-col cols="12" md="6">
                <div class="text-body-1 font-weight-medium mb-1">External Resources</div>
                
                <div v-if="!exercisesStore.currentExercise.links || exercisesStore.currentExercise.links.length === 0" class="text-body-2 mb-4 text-grey">
                  No external resources available
                </div>
                
                <v-list v-else density="compact" class="mb-4">
                  <v-list-item 
                    v-for="(link, index) in exercisesStore.currentExercise.links" 
                    :key="index"
                  >
                    <template v-slot:prepend>
                      <v-icon icon="mdi-link" size="small"></v-icon>
                    </template>
                    
                    <v-list-item-title>
                      <a :href="link.url" target="_blank" rel="noopener noreferrer">
                        {{ link.description || link.url }}
                      </a>
                    </v-list-item-title>
