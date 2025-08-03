document.addEventListener('DOMContentLoaded', function () {
  bootstrapListener();
  htmx.onLoad(bootstrapListener);
});

function bootstrapListener() {
  const filterForm = document.querySelector('form#datetime-filter')
  if (filterForm) {
    filterForm.addEventListener('submit', setDatetimeFormValues);
  }
}

function setDatetimeFormValues() {
  const startDate = document.querySelector('#start').value;
  const endDate = document.querySelector('#end').value;

  // Construct Date objects in local time
  const startDateTime = startDate ? new Date(startDate + 'T00:00:00') : null;
  const endDateTime = endDate ? new Date(endDate + 'T23:59:59') : null;

  // Format as ISO string with timezone offset (e.g., 2025-08-03T23:59:59+07:00)
  const toIsoWithOffset = date =>
    date.toLocaleString('sv-SE', { hour12: false }).replace(' ', 'T') +
    getOffsetString(date);

  function getOffsetString(date) {
    const offset = -date.getTimezoneOffset(); // minutes
    const sign = offset >= 0 ? '+' : '-';
    const pad = n => String(Math.floor(Math.abs(n))).padStart(2, '0');
    return sign + pad(offset / 60) + ':' + pad(offset % 60);
  }

  if (startDateTime) {
    document.querySelector('#startDateTime').value = toIsoWithOffset(startDateTime);
  }
  if (endDateTime) {
    document.querySelector('#endDateTime').value = toIsoWithOffset(endDateTime);
  }
}