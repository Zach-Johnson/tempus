import axios from "axios";

// Create axios instance with default config
const api = axios.create({
    baseURL: "/api/v1",
    headers: { "Content-Type": "application/json" },
    timeout: 10000,
});

// Request interceptor for adding auth token if needed
api.interceptors.request.use(
    (config) => { // You can add auth token logic here if needed
        return config;
    },
    (error) => {
        return Promise.reject(error);
    },
);

// Response interceptor for error handling
api.interceptors.response.use((response) => {
    return response;
}, (error) => {
    // Handle common errors
    const { response } = error;
    if (response) {
        // Log the error details
        console.error("API Error:", response.status, response.data);
    } else if (error.request) {
        // The request was made but no response was received
        console.error("Network Error:", error.request);
    } else {
        // Something happened in setting up the request
        console.error("Request Error:", error.message);
    }

    return Promise.reject(error);
});

// Categories API
const categoriesAPI = {
    getAll: (params = {}) => api.get("/categories", { params }),
    get: (id) => api.get(`/categories/${id}`),
    create: (data) => api.post("/categories", data),
    update: (id, data, updateMask) =>
        api.patch(
            `/categories/${id}`,
            { category: data, update_mask: updateMask },
        ),
    delete: (id) => api.delete(`/categories/${id}`),
};

// Tags API
const tagsAPI = {
    getAll: (params = {}) => api.get("/tags", { params }),
    get: (id) => api.get(`/tags/${id}`),
    create: (data) => api.post("/tags", data),
    update: (id, data, updateMask) =>
        api.patch(`/tags/${id}`, { tag: data, update_mask: updateMask }),
    delete: (id) => api.delete(`/tags/${id}`),
};

// Exercises API
const exercisesAPI = {
    getAll: (params = {}) => api.get("/exercises", { params }),
    get: (id) => api.get(`/exercises/${id}`),
    create: (data) => api.post("/exercises", data),
    update: (id, data, updateMask) =>
        api.patch(
            `/exercises/${id}`,
            { exercise: data, update_mask: updateMask },
        ),
    delete: (id) => api.delete(`/exercises/${id}`),
    addImage: (exerciseId, imageData) => {
        const formData = new FormData();
        formData.append("image_data", imageData.blob);
        formData.append("filename", imageData.filename);
        formData.append("mime_type", imageData.mimeType);
        formData.append("description", imageData.description || "");

        return api.post(`/exercises/${exerciseId}/images`, formData, {
            headers: {
                "Content-Type": "multipart/form-data",
            },
        });
    },
    deleteImage: (id) => api.delete(`/exercise-images/${id}`),
    addLink: (exerciseId, data) =>
        api.post(`/exercises/${exerciseId}/links`, data),
    deleteLink: (id) => api.delete(`/exercise-links/${id}`),
    getStats: (exerciseId, params = {}) =>
        api.get(`/exercises/${exerciseId}/stats`, { params }),
};

// Practice Sessions API
const sessionsAPI = {
    getAll: (params = {}) => api.get("/sessions", { params }),
    get: (id) => api.get(`/sessions/${id}`),
    create: (data) => api.post("/sessions", data),
    update: (id, data, updateMask) =>
        api.patch(`/sessions/${id}`, {
            session: data,
            update_mask: updateMask,
        }),
    delete: (id) => api.delete(`/sessions/${id}`),
    addExercise: (sessionId, data) =>
        api.post(`/sessions/${sessionId}/exercises`, data),
    updateExercise: (sessionId, exerciseId, data, updateMask) =>
        api.patch(
            `/sessions/${sessionId}/exercises/${exerciseId}`,
            { exercise: data, update_mask: updateMask },
        ),
    deleteExercise: (id) => api.delete(`/session-exercises/${id}`),
    getStats: (params = {}) => api.get("/sessions/stats", { params }),
};

// Exercise History API
const historyAPI = {
    getAll: (params = {}) => api.get("/history", { params }),
    get: (id) => api.get(`/history/${id}`),
    create: (data) => api.post("/history", data),
    update: (id, data, updateMask) =>
        api.patch(`/history/${id}`, { history: data, update_mask: updateMask }),
    delete: (id) => api.delete(`/history/${id}`),
};

export {
    api as default,
    categoriesAPI,
    exercisesAPI,
    historyAPI,
    sessionsAPI,
    tagsAPI,
};
