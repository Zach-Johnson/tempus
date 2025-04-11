<template>
  <div class="app-container">
    <Navbar />
    <div class="main-content">
      <!-- <Sidebar /> -->
      <div class="content-area">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script>
import Navbar from '@/components/common/Navbar.vue';
// import Sidebar from '@/components/common/Sidebar.vue';

export default {
  name: 'App',
  components: {
    Navbar,
    // Sidebar
  }
};
</script>

<style>
.app-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.main-content {
  display: flex;
  flex: 1;
}

.content-area {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

/* Global styles */
body {
  font-family: 'Roboto', sans-serif;
  background-color: #f8f9fa;
  margin: 0;
  padding: 0;
}

h1, h2, h3, h4, h5, h6 {
  font-weight: 500;
  margin-bottom: 1rem;
}

.card {
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05);
  border: none;
  margin-bottom: 1.5rem;
}

.card-header {
  background-color: #ffffff;
  border-bottom: 1px solid #f0f0f0;
  padding: 1rem 1.25rem;
}

.form-label {
  font-weight: 500;
}

.btn-primary {
  background-color: #007bff;
  border-color: #007bff;
}

.btn-primary:hover {
  background-color: #0069d9;
  border-color: #0062cc;
}
</style>
