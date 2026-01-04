new Chart(ctx, {
  type: 'bar',
  data: {
    labels: ['DPS'],
    datasets: [{ data: [window.dps] }]
  }
});