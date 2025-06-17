document.addEventListener('DOMContentLoaded', function () {
  initTimers();
  htmx.onLoad(initTimers);
});

function initTimers() {
  document.querySelectorAll('[id^="timer-"]').forEach(el => {
    if (el._timerInitialized) return;
    el._timerInitialized = true;

    const startTime = el.dataset.startTime;
    if (!startTime || startTime.trim() === "") {
      el.textContent = "00:00:00";
      return;
    }

    let startDate;
    try {
      // Try parsing as ISO string first
      startDate = new Date(startTime);

      // Fallback for Unix timestamp (if number)
      if (isNaN(startDate.getTime()) && !isNaN(startTime)) {
        startDate = new Date(parseInt(startTime) * 1000);
      }
    } catch (e) {
      console.error("Date parsing error:", e);
      el.textContent = "00:00:00";
      return;
    }

    if (isNaN(startDate.getTime())) {
      console.error("Invalid date:", startTime);
      el.textContent = "00:00:00";
      return;
    }

    function updateTimer() {
      const now = new Date();
      const diff = now - startDate;

      // Safeguard against negative time
      if (diff < 0) {
        el.textContent = "00:00:00";
        return;
      }

      const totalSeconds = Math.floor(diff / 1000);
      const hours = Math.floor(totalSeconds / 3600);
      const minutes = Math.floor((totalSeconds % 3600) / 60);
      const seconds = totalSeconds % 60;

      el.textContent =
        `${hours.toString().padStart(2, '0')}:` +
        `${minutes.toString().padStart(2, '0')}:` +
        `${seconds.toString().padStart(2, '0')}`;
    }

    updateTimer();
    const intervalId = setInterval(updateTimer, 1000);

    // Store interval ID on element
    el._timerInterval = intervalId;
  });
}

// Enhanced cleanup
document.addEventListener('htmx:beforeSwap', function (evt) {
  const target = evt.detail.target instanceof Element ? evt.detail.target : document.body;
  target.querySelectorAll('[id^="timer-"]').forEach(el => {
    if (el._timerInterval) {
      clearInterval(el._timerInterval);
      delete el._timerInterval;
      delete el._timerInitialized;
    }
  });
});