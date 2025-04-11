import {
    addDays,
    endOfMonth,
    endOfWeek,
    format,
    parseISO,
    startOfMonth,
    startOfWeek,
    subDays,
} from "date-fns";

/**
 * Format a date for display
 * @param {string|Date} date - Date to format
 * @param {string} formatString - Format pattern (default: 'MMM d, yyyy')
 * @returns {string} Formatted date string
 */
export const formatDate = (date, formatString = "MMM d, yyyy") => {
    if (!date) return "";

    // If date is a string, parse it
    const dateObj = typeof date === "string" ? parseISO(date) : date;
    return format(dateObj, formatString);
};

/**
 * Format a time for display
 * @param {string|Date} date - Date to format
 * @param {string} formatString - Format pattern (default: 'h:mm a')
 * @returns {string} Formatted time string
 */
export const formatTime = (date, formatString = "h:mm a") => {
    if (!date) return "";

    // If date is a string, parse it
    const dateObj = typeof date === "string" ? parseISO(date) : date;
    return format(dateObj, formatString);
};

/**
 * Format a date and time for display
 * @param {string|Date} date - Date to format
 * @param {string} formatString - Format pattern (default: 'MMM d, yyyy h:mm a')
 * @returns {string} Formatted date and time string
 */
export const formatDateTime = (date, formatString = "MMM d, yyyy h:mm a") => {
    if (!date) return "";

    // If date is a string, parse it
    const dateObj = typeof date === "string" ? parseISO(date) : date;
    return format(dateObj, formatString);
};

/**
 * Get the current date range for this week
 * @returns {Object} Object with startDate and endDate
 */
export const getCurrentWeekRange = () => {
    const today = new Date();
    return {
        startDate: startOfWeek(today, { weekStartsOn: 1 }), // Start on Monday
        endDate: endOfWeek(today, { weekStartsOn: 1 }),
    };
};

/**
 * Get the current date range for this month
 * @returns {Object} Object with startDate and endDate
 */
export const getCurrentMonthRange = () => {
    const today = new Date();
    return { startDate: startOfMonth(today), endDate: endOfMonth(today) };
};

/**
 * Get the date range for the last N days
 * @param {number} days - Number of days to look back
 * @returns {Object} Object with startDate and endDate
 */
export const getLastNDaysRange = (days) => {
    const today = new Date();
    return { startDate: subDays(today, days), endDate: today };
};

/**
 * Format a timestamp from the API as a human-readable relative time
 * @param {string|Date} date - Date to format
 * @returns {string} Relative time string
 */
export const getRelativeTimeString = (date) => {
    if (!date) return "";

    // If date is a string, parse it
    const dateObj = typeof date === "string" ? parseISO(date) : date;

    const now = new Date();
    const diffMs = now - dateObj;
    const diffSecs = Math.floor(diffMs / 1000);
    const diffMins = Math.floor(diffSecs / 60);
    const diffHours = Math.floor(diffMins / 60);
    const diffDays = Math.floor(diffHours / 24);

    if (diffSecs < 60) {
        return "just now";
    } else if (diffMins < 60) {
        return `${diffMins} minute${diffMins !== 1 ? "s" : ""} ago`;
    } else if (diffHours < 24) {
        return `${diffHours} hour${diffHours !== 1 ? "s" : ""} ago`;
    } else if (diffDays < 7) {
        return `${diffDays} day${diffDays !== 1 ? "s" : ""} ago`;
    } else {
        return formatDate(dateObj);
    }
};
