import {
    differenceInMinutes,
    format,
    formatDistanceToNow,
    parseISO,
} from "date-fns";
import { defineStore } from "pinia";
import { computed, ref } from "vue";

export const useAppStore = defineStore("app", () => {
    // State
    const isLoading = ref(false);
    const darkMode = ref(false);
    const snackbar = ref({
        show: false,
        text: "",
        color: "success", // 'success', 'info', 'warning', 'error'
        timeout: 3000,
    });

    // Utility functions for formatting (to be used across components)
    function formatDate(dateString) {
        if (!dateString) return "";

        try {
            const date = typeof dateString === "string"
                ? parseISO(dateString)
                : dateString;
            return format(date, "MMM d, yyyy");
        } catch (err) {
            console.error("Error formatting date:", err);
            return "";
        }
    }

    function formatDateTime(dateString) {
        if (!dateString) return "";

        try {
            const date = typeof dateString === "string"
                ? parseISO(dateString)
                : dateString;
            return format(date, "MMM d, yyyy h:mm a");
        } catch (err) {
            console.error("Error formatting date time:", err);
            return "";
        }
    }

    function formatTimeAgo(dateString) {
        if (!dateString) return "";

        try {
            const date = typeof dateString === "string"
                ? parseISO(dateString)
                : dateString;
            return formatDistanceToNow(date, { addSuffix: true });
        } catch (err) {
            console.error("Error formatting time ago:", err);
            return "";
        }
    }

    function formatDuration(startTime, endTime) {
        if (!startTime || !endTime) return "";

        try {
            const start = typeof startTime === "string"
                ? parseISO(startTime)
                : startTime;
            const end = typeof endTime === "string"
                ? parseISO(endTime)
                : endTime;

            const minutes = differenceInMinutes(end, start);
            return formatMinutes(minutes);
        } catch (err) {
            console.error("Error calculating duration:", err);
            return "";
        }
    }

    function formatMinutes(minutes) {
        if (minutes === undefined || minutes === null) return "0min";

        if (minutes < 60) return `${minutes}min`;

        const hours = Math.floor(minutes / 60);
        const remainingMinutes = minutes % 60;

        if (remainingMinutes === 0) return `${hours}h`;
        return `${hours}h ${remainingMinutes}min`;
    }

    // Snackbar functions
    function showSnackbar(text, color = "success", timeout = 3000) {
        snackbar.value = { show: true, text, color, timeout };
    }

    function hideSnackbar() {
        snackbar.value.show = false;
    }

    function showSuccessMessage(text) {
        showSnackbar(text, "success");
    }

    function showErrorMessage(text) {
        showSnackbar(text, "error", 5000);
    }

    function showWarningMessage(text) {
        showSnackbar(text, "warning", 4000);
    }

    function showInfoMessage(text) {
        showSnackbar(text, "info");
    }

    // Toggle dark mode
    function toggleDarkMode() {
        darkMode.value = !darkMode.value;
        localStorage.setItem("darkMode", darkMode.value ? "true" : "false");
    }

    function setDarkMode(value) {
        darkMode.value = value;
        localStorage.setItem("darkMode", darkMode.value ? "true" : "false");
    }

    // Initialize dark mode from localStorage on app startup
    function initDarkMode() {
        const storedDarkMode = localStorage.getItem("darkMode");
        if (storedDarkMode !== null) {
            darkMode.value = storedDarkMode === "true";
        } else {
            // Check if user prefers dark mode via OS settings
            darkMode.value = window.matchMedia &&
                window.matchMedia("(prefers-color-scheme: dark)").matches;
        }
    }

    return {
        // State
        isLoading,
        darkMode,
        snackbar,

        // Formatting functions
        formatDate,
        formatDateTime,
        formatTimeAgo,
        formatDuration,
        formatMinutes,

        // Snackbar functions
        showSnackbar,
        hideSnackbar,
        showSuccessMessage,
        showErrorMessage,
        showWarningMessage,
        showInfoMessage,

        // Dark mode functions
        toggleDarkMode,
        setDarkMode,
        initDarkMode,
    };
});
