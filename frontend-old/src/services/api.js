import axios from "axios";

// Create an axios instance for the API
const api = axios.create({
    baseURL: import.meta.env.VUE_APP_API_URL || "http://localhost:8081",
    headers: {
        "Content-Type": "application/json",
        "Accept": "application/json",
    },
});

// Request interceptor for API calls
api.interceptors.request.use(
    (config) => {
        // You can add auth headers here if needed
        return config;
    },
    (error) => {
        return Promise.reject(error);
    },
);

// Response interceptor for API calls
api.interceptors.response.use(
    (response) => {
        return response.data;
    },
    async (error) => {
        const { response } = error;

        // Handle different error status codes
        if (response) {
            switch (response.status) {
                case 401:
                    console.error("Unauthorized access");
                    break;
                case 404:
                    console.error("Resource not found");
                    break;
                case 500:
                    console.error("Server error");
                    break;
                default:
                    console.error(
                        "An error occurred:",
                        response.data.message || "Unknown error",
                    );
            }

            return Promise.reject(response.data);
        }

        return Promise.reject(error);
    },
);

export default api;
