document.addEventListener('DOMContentLoaded', function () {
  formatTime();
  htmx.onLoad(formatTime);
});

function formatTime() {
  document
    .querySelectorAll('[data-timestamp]')
    .forEach(timestamp => {
      const timeVal = timestamp.dataset.timestamp;
      const dateText = new Date(timeVal);
      timestamp.innerHTML = formatToCustomLocal(dateText);
    });
}

function formatToCustomLocal(date) {
  const day = String(date.getDate()).padStart(2, '0');
  const month = date.toLocaleString('en-US', { month: 'short' });
  const year = String(date.getFullYear()).slice(2);
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');

  // Extract timezone abbreviation (e.g., WIB)
  const tzAbbr = date
    .toLocaleTimeString('en-us', { timeZoneName: 'short' })
    .split(' ')
    .pop();

  return `${day} ${month} ${year} ${hours}:${minutes} ${tzAbbr}`;
}
