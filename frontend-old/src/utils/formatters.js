/**
 * Format a number of minutes as hours and minutes
 * @param {number} minutes - Total minutes
 * @returns {string} Formatted duration string
 */
export const formatDuration = (minutes) => {
    if (!minutes && minutes !== 0) return "";

    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;

    if (hours === 0) {
        return `${mins} min${mins !== 1 ? "s" : ""}`;
    } else if (mins === 0) {
        return `${hours} hour${hours !== 1 ? "s" : ""}`;
    } else {
        return `${hours} hour${hours !== 1 ? "s" : ""} ${mins} min${
            mins !== 1 ? "s" : ""
        }`;
    }
};

/**
 * Convert BPM to a readable speed description
 * @param {number} bpm - Beats per minute
 * @returns {string} Speed description
 */
export const getBpmSpeedDescription = (bpm) => {
    if (!bpm) return "N/A";

    if (bpm < 60) {
        return "Slow";
    } else if (bpm < 90) {
        return "Moderate";
    } else if (bpm < 120) {
        return "Medium";
    } else if (bpm < 160) {
        return "Fast";
    } else {
        return "Very Fast";
    }
};

/**
 * Convert a rating to a text description
 * @param {number} rating - Rating (1-5)
 * @returns {string} Rating description
 */
export const getRatingDescription = (rating) => {
    if (!rating) return "Not Rated";

    const descriptions = [
        "Poor", // 1
        "Fair", // 2
        "Good", // 3
        "Very Good", // 4
        "Excellent", // 5
    ];

    return descriptions[Math.min(Math.max(Math.floor(rating) - 1, 0), 4)];
};

/**
 * Truncate a string to a given length and add ellipsis if needed
 * @param {string} text - Text to truncate
 * @param {number} length - Maximum length
 * @returns {string} Truncated text
 */
export const truncateText = (text, length = 100) => {
    if (!text) return "";
    if (text.length <= length) return text;

    return text.substring(0, length).trim() + "...";
};

/**
 * Format a time signature
 * @param {string} timeSignature - Time signature (e.g., '4/4', '3/4')
 * @returns {string} Formatted time signature
 */
export const formatTimeSignature = (timeSignature) => {
    if (!timeSignature) return "Common Time (4/4)";

    const commonTimeSignatures = {
        "4/4": "Common Time (4/4)",
        "3/4": "Waltz Time (3/4)",
        "2/4": "March Time (2/4)",
        "6/8": "Compound Duple (6/8)",
        "9/8": "Compound Triple (9/8)",
        "12/8": "Compound Quadruple (12/8)",
        "5/4": "Quintuple (5/4)",
        "7/8": "Septuple (7/8)",
    };

    return commonTimeSignatures[timeSignature] || timeSignature;
};
