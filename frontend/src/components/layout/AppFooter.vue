<template>
  <v-footer app padless class="text-center">
    <v-card flat tile width="100%" class="transparent">
      <v-card-text>
        <v-btn
          v-for="icon in socialIcons"
          :key="icon.icon"
          :icon="icon.icon"
          class="mx-2"
          variant="text"
          target="_blank"
          :href="icon.url"
        ></v-btn>
      </v-card-text>

      <v-divider></v-divider>

      <v-card-text class="text-body-2">
        {{ new Date().getFullYear() }} — <strong>Tempus Drum Practice Tracker</strong>
      </v-card-text>
    </v-card>
  </v-footer>
</template>

<script setup>
const socialIcons = [
  { icon: 'mdi-github', url: 'https://github.com/Zach-Johnson/tempus' },
]
</script>
