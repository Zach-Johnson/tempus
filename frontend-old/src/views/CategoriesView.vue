<template>
  <div class="categories-view">
    <div class="container-fluid">
      <div class="page-header mb-4">
        <h1>Categories</h1>
        <p class="text-muted">Organize your drum exercises into categories for easy access.</p>
      </div>
      
      <!-- Categories list component -->
      <CategoryList />
    </div>
  </div>
</template>

<script>
import CategoryList from '@/components/categories/CategoryList.vue';

export default {
  name: 'CategoriesView',
  
  components: {
    CategoryList
  }
};
</script>

<style scoped>
.page-header h1 {
  font-weight: 600;
  color: #212529;
}
</style>
