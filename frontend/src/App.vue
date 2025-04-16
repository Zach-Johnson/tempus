<template>
  <v-app>
    <app-header />
    
    <v-main :class="{ 'with-sidebar': isSidebarOpen }">
      <v-container fluid>
        <router-view />
      </v-container>
    </v-main>
    
    <app-footer />
    
    <!-- Global snackbar -->
    <v-snackbar
      v-model="appStore.snackbar.show"
      :timeout="appStore.snackbar.timeout"
      :color="appStore.snackbar.color"
      location="top"
    >
      {{ appStore.snackbar.text }}
      <template v-slot:actions>
        <v-btn variant="text" @click="appStore.hideSnackbar">
          Close
        </v-btn>
      </template>
    </v-snackbar>
  </v-app>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useAppStore } from '@/stores/app.js';
import AppHeader from '@/components/layout/AppHeader.vue';
import AppFooter from '@/components/layout/AppFooter.vue';

const appStore = useAppStore();
const isSidebarOpen = ref(true); // Default to open

onMounted(() => {
  // Initialize app settings
  appStore.initDarkMode();
  
  // Initialize sidebar state from localStorage
  const storedOpen = localStorage.getItem("sidebarOpen");
  if (storedOpen !== null) {
    isSidebarOpen.value = storedOpen === "true";
  }
  
  // Listen for storage events to keep sidebar state in sync across components
  window.addEventListener('storage', (event) => {
    if (event.key === 'sidebarOpen') {
      isSidebarOpen.value = event.newValue === 'true';
    }
  });
});
</script>

<style>
/* Add padding when sidebar is open */
.with-sidebar {
  padding-left: 256px; /* Width of drawer */
}

@media (max-width: 960px) {
  .with-sidebar {
    padding-left: 0; /* Remove padding on small screens */
  }
}
</style>
