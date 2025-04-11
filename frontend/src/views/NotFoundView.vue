<template>
  <div class="not-found text-center">
    <v-container>
      <v-row justify="center">
        <v-col cols="12" sm="10" md="8" lg="6">
          <v-card class="pa-6" variant="flat">
            <div class="mb-5">
              <v-icon
                icon="mdi-drum"
                size="x-large"
                color="primary"
                class="mb-4"
              ></v-icon>
              <h1 class="text-h2 mb-2">404</h1>
              <h2 class="text-h4 mb-6">Page Not Found</h2>
              <p class="text-body-1 mb-8">
                The page you are looking for doesn't exist or has been moved.
              </p>
            </div>
            
            <div class="d-flex justify-center flex-wrap">
              <v-btn
                color="primary"
                size="large"
                prepend-icon="mdi-home"
                class="ma-2"
                :to="{ name: 'home' }"
              >
                Return to Dashboard
              </v-btn>
              <v-btn
                variant="outlined"
                size="large"
                prepend-icon="mdi-arrow-left"
                class="ma-2"
                @click="goBack"
              >
                Go Back
              </v-btn>
            </div>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

function goBack() {
  router.go(-1)
}
</script>

<style scoped>
.not-found {
  padding: 5% 0;
}

@media (min-width: 600px) {
  .not-found {
    padding: 10% 0;
  }
}
</style>
