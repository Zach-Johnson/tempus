<template>
  <v-app-bar>
    <v-btn icon @click="toggleSidebar">
      <v-icon>{{ isSidebarOpen ? 'mdi-menu-open' : 'mdi-menu' }}</v-icon>
    </v-btn>
    
    <v-app-bar-title>
      <router-link to="/" class="text-decoration-none">
        Tempus
      </router-link>
    </v-app-bar-title>
    
    <v-spacer></v-spacer>
    
    <v-btn icon @click="toggleTheme">
      <v-icon>{{ isDarkTheme ? 'mdi-weather-sunny' : 'mdi-weather-night' }}</v-icon>
    </v-btn>
  </v-app-bar>
  
  <v-navigation-drawer :model-value="isSidebarOpen" @update:model-value="updateSidebarState">
    <v-list-item class="pa-4">
      <v-list-item-title class="text-h6">
        Tempus
      </v-list-item-title>
      <v-list-item-subtitle>
        Drum Practice Tracker
      </v-list-item-subtitle>
    </v-list-item>
    
    <v-divider></v-divider>
    
    <v-list density="compact" nav>
      <v-list-item v-for="item in menuItems" :key="item.title" :to="item.to" :prepend-icon="item.icon">
        {{ item.title }}
      </v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useTheme } from 'vuetify'
import { useAppStore } from '@/stores/app.js'

const theme = useTheme()
const appStore = useAppStore()
const isSidebarOpen = ref(true) // Default to open

const isDarkTheme = computed(() => theme.global.name.value === 'dark')

const menuItems = [
  { title: 'Dashboard', icon: 'mdi-view-dashboard', to: '/' },
  { title: 'Practice Sessions', icon: 'mdi-calendar-clock', to: '/sessions' },
  { title: 'New Session', icon: 'mdi-calendar-plus', to: '/sessions/new' },
  { title: 'Exercises', icon: 'mdi-music-note', to: '/exercises' },
  { title: 'Categories', icon: 'mdi-folder', to: '/categories' },
  { title: 'Tags', icon: 'mdi-tag-multiple', to: '/tags' },
]

function toggleSidebar() {
  isSidebarOpen.value = !isSidebarOpen.value;
  localStorage.setItem("sidebarOpen", isSidebarOpen.value ? "true" : "false");
}

function updateSidebarState(newValue) {
  isSidebarOpen.value = newValue;
  localStorage.setItem("sidebarOpen", newValue ? "true" : "false");
}

function toggleTheme() {
  appStore.toggleDarkMode()
  theme.global.name.value = isDarkTheme.value ? 'light' : 'dark'
}

// Initialize from local storage
onMounted(() => {
  // Initialize theme
  appStore.initDarkMode()
  theme.global.name.value = appStore.darkMode ? 'dark' : 'light'
  
  // Initialize sidebar state from localStorage
  const storedOpen = localStorage.getItem("sidebarOpen");
  if (storedOpen !== null) {
    isSidebarOpen.value = storedOpen === "true";
  }
})
</script>

<style scoped>
.v-app-bar-title a {
  color: inherit;
  font-weight: 500;
}
</style>
