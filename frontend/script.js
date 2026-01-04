const form = document.getElementById("job-form");
const textarea = document.getElementById("job-input");
const outputContainer = document.getElementById("job-output");
const jobList = document.getElementById("job-list");

let jobs = [];

form.addEventListener("submit", async (e) => {
    e.preventDefault();

    const res = await fetch("/api/jobs", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ profile: textarea.value })
    });
    const data = await res.json();
    const jobId = data.job_id;

    jobs.push({ id: jobId, status: "queued" });
    renderJobs();

    outputContainer.textContent = `Job ${jobId} started...\n`;

    const evtSource = new EventSource(`/api/jobs/${jobId}/stream`);
    evtSource.onmessage = function(e) {
        outputContainer.textContent += e.data + "\n";
    };
    evtSource.onerror = function() {
        outputContainer.textContent += "\n--- Stream closed ---\n";
        evtSource.close();
    };

    // Status Polling
    const pollStatus = setInterval(async () => {
        const statusRes = await fetch(`/api/jobs/${jobId}/progress`);
        const statusData = await statusRes.json();
        updateJobStatus(jobId, statusData.progress);
        renderJobs();
        if (["done", "failed"].includes(statusData.progress)) {
            clearInterval(pollStatus);
        }
    }, 200);
});

function updateJobStatus(jobId, status) {
    const job = jobs.find(j => j.id === jobId);
    if (job) job.status = status;
}

function renderJobs() {
    jobList.innerHTML = "";
    jobs.forEach(job => {
        const div = document.createElement("div");
        div.className = `job ${job.status}`;
        div.textContent = `Job ${job.id} - ${job.status}`;
        jobList.appendChild(div);
    });
}
