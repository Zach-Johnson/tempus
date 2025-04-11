<template>
  <div class="exercise-list">
    <div v-if="loading" class="d-flex align-center">
      <v-progress-circular indeterminate size="20" width="2" color="primary" class="mr-2"></v-progress-circular>
      <span class="text-body-2">Loading exercises...</span>
    </div>
    
    <div v-else-if="!exercises || exercises.length === 0" class="text-body-2 text-grey">
      {{ noExercisesMessage }}
    </div>
    
    <div v-else>
      <div v-if="displayType === 'table'" class="exercise-table-container">
        <v-data-table
          :headers="tableHeaders"
          :items="exercises"
          :items-per-page="itemsPerPage"
          density="compact"
          class="elevation-1"
        >
          <template v-slot:item.description="{ item }">
            <span v-if="item.description">{{ truncateText(item.description, 80) }}</span>
            <span v-else class="text-grey">No description</span>
          </template>
          
          <template v-slot:item.categories="{ item }">
            <v-chip-group>
              <category-chip
                v-for="categoryId in item.categoryIds"
                :key="categoryId"
                :category="getCategoryById(categoryId)"
                size="x-small"
                class="mr-1"
              ></category-chip>
            </v-chip-group>
          </template>
          
          <template v-slot:item.tags="{ item }">
            <v-chip-group>
              <tag-chip
                v-for="tagId in item.tagIds"
                :key="tagId"
                :tag="getTagById(tagId)"
                size="x-small"
                class="mr-1"
              ></tag-chip>
            </v-chip-group>
          </template>
          
          <template v-slot:item.actions="{ item }">
            <v-btn
              icon
              variant="text"
              size="small"
              color="primary"
              :to="{ name: 'exercise-detail', params: { id: item.id }}"
              class="mr-1"
            >
              <v-icon>mdi-eye</v-icon>
            </v-btn>
            
            <v-btn
              v-if="allowSelect"
              icon
              variant="text"
              size="small"
              color="primary"
              @click="selectExercise(item)"
            >
              <v-icon>mdi-plus</v-icon>
            </v-btn>
          </template>
        </v-data-table>
      </div>
      
      <div v-else-if="displayType === 'list'" class="exercise-list-container">
        <v-list density="compact">
          <v-list-item
            v-for="exercise in exercises"
            :key="exercise.id"
            :to="selectable ? undefined : { name: 'exercise-detail', params: { id: exercise.id }}"
            @click="selectable ? selectExercise(exercise) : undefined"
          >
            <template v-slot:prepend>
              <v-avatar size="32" color="primary" class="text-white">
                {{ exercise.name.charAt(0).toUpperCase() }}
              </v-avatar>
            </template>
            
            <v-list-item-title>{{ exercise.name }}</v-list-item-title>
            <v-list-item-subtitle v-if="exercise.description">
              {{ truncateText(exercise.description, 60) }}
            </v-list-item-subtitle>
          </v-list-item>
        </v-list>
      </div>
      
      <div v-else class="exercise-grid-container">
        <v-row>
          <v-col 
            v-for="exercise in exercises" 
            :key="exercise.id"
            cols="12"
            sm="6"
            md="4"
            lg="3"
          >
            <v-card
              :to="selectable ? undefined : { name: 'exercise-detail', params: { id: exercise.id }}"
              @click="selectable ? selectExercise(exercise) : undefined"
              variant="outlined"
              class="h-100"
            >
              <v-card-title class="text-subtitle-1">{{ exercise.name }}</v-card-title>
              <v-card-text>
                <p v-if="exercise.description" class="text-body-2 mb-2">
                  {{ truncateText(exercise.description, 100) }}
                </p>
                
                <div v-if="exercise.categoryIds && exercise.categoryIds.length > 0" class="mb-2">
                  <v-chip-group>
                    <category-chip
                      v-for="categoryId in exercise.categoryIds.slice(0, 3)"
                      :key="categoryId"
                      :category="getCategoryById(categoryId)"
                      size="x-small"
                    ></category-chip>
                    <v-chip
                      v-if="exercise.categoryIds.length > 3"
                      size="x-small"
                      color="secondary"
                      variant="outlined"
                    >
                      +{{ exercise.categoryIds.length - 3 }}
                    </v-chip>
                  </v-chip-group>
                </div>
                
                <div v-if="exercise.tagIds && exercise.tagIds.length > 0">
                  <v-chip-group>
                    <tag-chip
                      v-for="tagId in exercise.tagIds.slice(0, 3)"
                      :key="tagId"
                      :tag="getTagById(tagId)"
                      size="x-small"
                    ></tag-chip>
                    <v-chip
                      v-if="exercise.tagIds.length > 3"
                      size="x-small"
                      color="primary"
                      variant="outlined"
                    >
                      +{{ exercise.tagIds.length - 3 }}
                    </v-chip>
                  </v-chip-group>
                </div>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </div>
      
      <div v-if="showPagination && totalCount > itemsPerPage" class="pagination-container mt-4">
        <v-pagination
          v-model="localPage"
          :length="Math.ceil(totalCount / itemsPerPage)"
          :total-visible="7"
        ></v-pagination>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useExercisesStore } from '@/stores/exercises.js'
import { useCategoriesStore } from '@/stores/categories.js'
import { useTagsStore } from '@/stores/tags.js'
import CategoryChip from '@/components/categories/CategoryChip.vue'
import TagChip from '@/components/tags/TagChip.vue'

const props = defineProps({
  exercises: {
    type: Array,
    default: () => []
  },
  categoryId: {
    type: [Number, String],
    default: null
  },
  tagId: {
    type: [Number, String],
    default: null
  },
  loading: {
    type: Boolean,
    default: false
  },
  displayType: {
    type: String,
    default: 'table',  // 'table', 'list', or 'grid'
    validator: (value) => ['table', 'list', 'grid'].includes(value)
  },
  noExercisesMessage: {
    type: String,
    default: 'No exercises available'
  },
  showPagination: {
    type: Boolean,
    default: false
  },
  totalCount: {
    type: Number,
    default: 0
  },
  itemsPerPage: {
    type: Number,
    default: 10
  },
  page: {
    type: Number,
    default: 1
  },
  selectable: {
    type: Boolean,
    default: false
  },
  allowSelect: {
    type: Boolean, 
    default: false
  },
  selectedExerciseIds: {
    type: Array,
    default: () => []
  }
