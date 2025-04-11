<template>
  <v-app-bar>
    <v-app-bar-nav-icon @click="toggleDrawer"></v-app-bar-nav-icon>
    
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
  
  <v-navigation-drawer v-model="drawer" temporary>
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
const drawer = ref(false)

const isDarkTheme = computed(() => theme.global.name.value === 'dark')

const menuItems = [
  { title: 'Dashboard', icon: 'mdi-view-dashboard', to: '/' },
  { title: 'Practice Sessions', icon: 'mdi-calendar-clock', to: '/sessions' },
  { title: 'New Session', icon: 'mdi-calendar-plus', to: '/sessions/new' },
  { title: 'Exercises', icon: 'mdi-music-note', to: '/exercises' },
  { title: 'Categories', icon: 'mdi-folder', to: '/categories' },
  { title: 'Tags', icon: 'mdi-tag-multiple', to: '/tags' },
  { title: 'Statistics', icon: 'mdi-chart-bar', to: '/stats' }
]

function toggleDrawer() {
  drawer.value = !drawer.value
}

function toggleTheme() {
  appStore.toggleDarkMode()
  theme.global.name.value = isDarkTheme.value ? 'light' : 'dark'
}

// Initialize theme from app store
onMounted(() => {
  appStore.initDarkMode()
  theme.global.name.value = appStore.darkMode ? 'dark' : 'light'
})
</script>

<style scoped>
.v-app-bar-title a {
  color: inherit;
  font-weight: 500;
}
</style>
