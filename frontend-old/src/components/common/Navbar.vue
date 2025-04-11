<template>
  <nav class="navbar navbar-expand-lg navbar-dark bg-primary">
    <div class="container-fluid">
      <router-link class="navbar-brand" to="/">
        <i class="bi bi-music-note-beamed me-2"></i>
        Drum Practice
      </router-link>
      
      <button 
        class="navbar-toggler" 
        type="button" 
        data-bs-toggle="collapse" 
        data-bs-target="#navbarMain" 
        aria-controls="navbarMain" 
        aria-expanded="false" 
        aria-label="Toggle navigation"
      >
        <span class="navbar-toggler-icon"></span>
      </button>
      
      <div class="collapse navbar-collapse" id="navbarMain">
        <ul class="navbar-nav me-auto mb-2 mb-lg-0">
          <li class="nav-item">
            <router-link class="nav-link" to="/" active-class="active" exact>
              Dashboard
            </router-link>
          </li>
          <li class="nav-item">
            <router-link class="nav-link" to="/practice" active-class="active">
              Practice Sessions
            </router-link>
          </li>
          <li class="nav-item">
            <router-link class="nav-link" to="/exercises" active-class="active">
              Exercises
            </router-link>
          </li>
          <li class="nav-item">
            <router-link class="nav-link" to="/categories" active-class="active">
              Categories
            </router-link>
          </li>
          <li class="nav-item">
            <router-link class="nav-link" to="/tags" active-class="active">
              Tags
            </router-link>
          </li>
          <!-- <li class="nav-item"> -->
          <!--   <router-link class="nav-link" to="/history" active-class="active"> -->
          <!--     Progress Tracking -->
          <!--   </router-link> -->
          <!-- </li> -->
        </ul>
        
        <ul class="navbar-nav ms-auto">
          <li class="nav-item">
            <button 
              class="btn btn-success" 
              @click="openNewSessionModal"
            >
              <i class="bi bi-plus-circle me-1"></i> New Session
            </button>
          </li>
        </ul>
      </div>
    </div>
  </nav>
  
  <!-- New Session Modal -->
  <div 
    class="modal fade" 
    id="newSessionModal" 
    tabindex="-1" 
    aria-labelledby="newSessionModalLabel" 
    aria-hidden="true"
  >
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="newSessionModalLabel">Start New Practice Session</h5>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="createSession">
            <div class="mb-3">
              <label for="sessionDate" class="form-label">Date</label>
              <input 
                type="datetime-local" 
                class="form-control" 
                id="sessionDate" 
                v-model="newSession.date" 
                required
              >
            </div>
            <div class="mb-3">
              <label for="sessionDuration" class="form-label">Duration (minutes)</label>
              <input 
                type="number" 
                class="form-control" 
                id="sessionDuration" 
                v-model="newSession.duration_minutes" 
                min="1" 
                required
              >
            </div>
            <div class="mb-3">
              <label for="sessionNotes" class="form-label">Notes</label>
              <textarea 
                class="form-control" 
                id="sessionNotes" 
                v-model="newSession.notes" 
                rows="3"
              ></textarea>
            </div>
            <div class="text-end">
              <button type="button" class="btn btn-secondary me-2" data-bs-dismiss="modal">Cancel</button>
              <button type="submit" class="btn btn-primary">Start Session</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useSessionsStore } from '@/store/modules/sessions';
import { Modal } from 'bootstrap'

export default {
  name: 'Navbar',
  
  setup() {
    const router = useRouter();
    const sessionsStore = useSessionsStore();
    
    // Initialize new session data with current date
    const newSession = ref({
      date: new Date().toISOString().slice(0, -8), // Format as YYYY-MM-DDTHH:MM
      duration_minutes: 30,
      notes: ''
    });
    
    let sessionModal = null;
    
      // Initialize Bootstrap modal
    onMounted(() => {
        const modalEl = document.getElementById('newSessionModal')
        if (modalEl) {
            sessionModal = new Modal(modalEl)
        }
    });
    
    const openNewSessionModal = () => {
      // Reset the form data with current date
      newSession.value = {
        date: new Date().toISOString().slice(0, -8),
        duration_minutes: 30,
        notes: ''
      };
      
      // Show the modal
      sessionModal.show();
    };
    
    const createSession = async () => {
      try {
        // Create the session
        const session = await sessionsStore.createSession(newSession.value);
        
        // Hide the modal
        sessionModal.hide();
        
        if (session && session.id) {
          // Navigate to the new session
          router.push(`/practice/${session.id}`);
        }
      } catch (error) {
        console.error('Failed to create session:', error);
        // Add error handling as needed
      }
    };
    
    return {
      newSession,
      openNewSessionModal,
      createSession
    };
  }
};
</script>

<style scoped>
.navbar {
  padding: 0.75rem 1rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.navbar-brand {
  font-weight: 700;
  font-size: 1.5rem;
}

.navbar-nav .nav-link {
  font-weight: 500;
  padding: 0.5rem 1rem;
}

.navbar-nav .nav-link.active {
  color: #fff;
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 0.25rem;
}
</style>
