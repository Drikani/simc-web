document.addEventListener("DOMContentLoaded", () => {
  const form = document.getElementById("job-form");
  const profileInput = document.getElementById("profileInput");
  const statusDiv = document.getElementById("status");
  const liveOutput = document.getElementById("liveOutput");

  form.addEventListener("submit", async (e) => {
    e.preventDefault();

    const profile = profileInput.value.trim();
    if (!profile) {
      alert("Please paste a SimC profile.");
      return;
    }

    // Status zurücksetzen
    statusDiv.textContent = "Submitting job...";
    liveOutput.textContent = "";

    try {
      // Job anlegen
      const res = await fetch("/api/jobs", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ profile }),
      });
      const data = await res.json();
      if (!data.job_id) {
        statusDiv.textContent = "Failed to create job";
        return;
      }

      const jobID = data.job_id;
      statusDiv.textContent = `Job submitted (ID: ${jobID}). Waiting for result...`;

      // SSE starten
      const evtSource = new EventSource(`/api/jobs/${jobID}/stream`);
      evtSource.onmessage = (event) => {
        try {
          const json = JSON.parse(event.data);
          // JSON hübsch formatieren
          liveOutput.textContent = JSON.stringify(json, null, 2);
        } catch (err) {
          liveOutput.textContent = event.data; // Falls keine JSON, z.B. Error
        }
      };

      evtSource.onerror = () => {
        statusDiv.textContent = "Job stream closed";
        evtSource.close();
      };
    } catch (err) {
      console.error(err);
      statusDiv.textContent = "Error submitting job";
    }
  });
});