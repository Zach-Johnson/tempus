<style scoped>
.sidebar {
  width: 280px;
  height: 100%;
  background-color: #f8f9fa;
  transition: width 0.3s;
  border-right: 1px solid #e9ecef;
}

.sidebar-collapsed {
  width: 60px;
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border-bottom: 1px solid #e9ecef;
}

.sidebar-header h5 {
  margin: 0;
  font-weight: 600;
}

.toggle-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 0.25rem;
  border-radius: 0.25rem;
}

.toggle-btn:hover {
  background-color: rgba(0, 0, 0, 0.05);
}

.sidebar-content {
  padding: 1rem;
}

.section-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: #6c757d;
  margin-bottom: 0.5rem;
  text-transform: uppercase;
}

.list-group-item {
  cursor: pointer;
  transition: background-color 0.2s;
  padding: 0.5rem 0.75rem;
}

.list-group-item:hover {
  background-color: rgba(0, 0, 0, 0.03);
}

.recent-item {
  display: flex;
  flex-direction: column;
  font-size: 0.875rem;
}

.recent-date {
  font-weight: 500;
}

.recent-duration {
  color: #6c757d;
  font-size: 0.75rem;
}

.icon-nav {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.icon-btn {
  font-size: 1.25rem;
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 0.25rem;
  color: #495057;
}

.icon-btn:hover {
  background-color: rgba(0, 0, 0, 0.05);
}
</style>
<template>
  <div class="sidebar" :class="{ 'sidebar-collapsed': isCollapsed }">
    <div class="sidebar-header">
      <h5 v-if="!isCollapsed">Quick Links</h5>
      <button class="toggle-btn" @click="toggleSidebar">
        <i :class="isCollapsed ? 'bi bi-chevron-right' : 'bi bi-chevron-left'"></i>
      </button>
    </div>
    
    <div class="sidebar-content">
      <!-- Quick actions when expanded -->
      <div v-if="!isCollapsed" class="quick-actions">
        <div class="section-title">Actions</div>
        <div class="action-buttons">
          <button 
            class="btn btn-sm btn-primary mb-2 w-100" 
            @click="navigateTo('/practice/new')"
          >
            <i class="bi bi-plus-circle me-2"></i> New Session
          </button>
          <button 
            class="btn btn-sm btn-outline-primary mb-2 w-100" 
            @click="navigateTo('/exercises/new')"
          >
            <i class="bi bi-music-note me-2"></i> New Exercise
          </button>
        </div>
      </div>
      
      <!-- Recent items -->
      <div v-if="!isCollapsed" class="recent-items mt-4">
        <div class="section-title">Recent Sessions</div>
        <ul class="list-group list-group-flush">
          <li 
            v-for="session in recentSessions" 
            :key="session.id" 
            class="list-group-item"
            @click="navigateTo(`/practice/${session.id}`)"
          >
            <div class="recent-item">
              <div class="recent-date">
                {{ formatDate(session.date) }}
              </div>
              <div class="recent-duration">
                {{ formatDuration(session.duration_minutes) }}
              </div>
            </div>
          </li>
          <li 
            v-if="recentSessions.length === 0" 
            class="list-group-item text-muted"
          >
            No recent sessions
          </li>
        </ul>
      </div>
      
      <!-- Collapsed quick icons -->
      <div v-if="isCollapsed" class="icon-nav">
        <button class="icon-btn" @click="navigateTo('/practice/new')">
          <i class="bi bi-plus-circle"></i>
        </button>
        <button class="icon-btn" @click="navigateTo('/exercises/new')">
          <i class="bi bi-music-note"></i>
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useSessionsStore } from '@/store/modules/sessions';
import { formatDate } from '@/utils/dateUtils';
import { formatDuration } from '@/utils/formatters';

export default {
  name: 'Sidebar',
  
  setup() {
    const router = useRouter();
    const sessionsStore = useSessionsStore();
    
    const isCollapsed = ref(localStorage.getItem('sidebar-collapsed') === 'true');
    
    // Load recent sessions
    onMounted(() => {
      sessionsStore.fetchSessions({ page: 1, page_size: 5 });
    });
    
    // Get recent sessions from store
    const recentSessions = computed(() => {
      return sessionsStore.sessions.slice(0, 5);
    });
    
    // Toggle sidebar collapsed state
    const toggleSidebar = () => {
      isCollapsed.value = !isCollapsed.value;
      localStorage.setItem('sidebar-collapsed', isCollapsed.value);
    };
    
    // Navigation helper
    const navigateTo = (route) => {
      router.push(route);
    };
    
    return {
      isCollapsed,
      recentSessions,
      toggleSidebar,
      navigateTo,
      formatDate,
      formatDuration
    };
  }
};
</script>
