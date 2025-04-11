<template>
  <v-app>
    <app-header />
    
    <v-main>
      <v-container fluid>
        <router-view />
      </v-container>
    </v-main>
  </v-app>
</template>

<script setup>
import AppHeader from '@/components/layout/AppHeader.vue'
</script>

<style>
.v-main {
  background-color: var(--v-background-base);
  min-height: 100vh;
}
</style>
