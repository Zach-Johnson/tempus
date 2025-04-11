import {sessionsAPI} from '@/services/api.js'
import {defineStore} from 'pinia'
import {computed, ref} from 'vue'

export const useSessionsStore = defineStore('sessions', () => {
    // State
    const sessions = ref([])
    const loading = ref(false)
    const error = ref(null)
    const currentSession = ref(null)
    const currentSessionLoading = ref(false)
    const practiceStats = ref(null)
    const statsLoading = ref(false)

    // Getters
    const sessionsSortedByDate = computed(
        () => {return [...sessions.value].sort(
            (a,
             b) => {// Sort by start_time in descending order (newest first)
                    return new Date(b.start_time) - new Date(a.start_time)})})

    const sessionById = computed(
        () => {return (id) =>
                          sessions.value.find(session => session.id === id)})

    const sessionsByExercise = computed(
        () => {return (exerciseId) => sessions.value.filter(
                          session => session.exercises &&
                              session.exercises.some(
                                  ex => ex.exercise_id === exerciseId))})

    const totalPracticeTime =
        computed(() => {return sessions.value.reduce((total, session) => {
                     if (session.start_time && session.end_time) {
                         const startTime = new Date(session.start_time)
                         const endTime = new Date(session.end_time)
                         const durationMinutes =
                             (endTime - startTime) / (1000 * 60)
                         return total + durationMinutes
                     }
                     return total
                 }, 0)})

    // Actions
    async function fetchSessions(params = {}) {
        loading.value = true
        error.value = null

        try {
            const response = await sessionsAPI.getAll(params)
            sessions.value = response.data.sessions || []
        } catch (err) {
            error.value = err.message || 'Failed to fetch practice sessions'
            console.error('Error fetching practice sessions:', err)
        } finally {
            loading.value = false
        }
    }

    async function fetchSession(id) {
        currentSessionLoading.value = true
        error.value = null

        try {
            const response = await sessionsAPI.get(id)
            currentSession.value = response.data

            // Also update the session in the sessions array if it exists
            const index = sessions.value.findIndex(s => s.id === id)
            if (index !== -1) {
                sessions.value[index] = response.data
            }
            else {
                sessions.value.push(response.data)
            }
        } catch (err) {
            error.value =
                err.message || `Failed to fetch practice session with ID ${id}`
            console.error(`Error fetching practice session ${id}:`, err)
        } finally {
            currentSessionLoading.value = false
        }
    }

    async function createSession(sessionData) {
