const startBtn = document.getElementById("startBtn");
const profileInput = document.getElementById("profileInput");
const statusEl = document.getElementById("status");
const loaderEl = document.getElementById("loader");
const liveOutputEl = document.getElementById("liveOutput");
const resultOutputEl = document.getElementById("resultOutput");

let eventSource = null;

function setStatus(text) {
  statusEl.textContent = text;
}

function showLoader(show) {
  loaderEl.style.display = show ? "block" : "none";
}

function resetUI() {
  liveOutputEl.textContent = "";
  resultOutputEl.textContent = "";
  setStatus("");
  showLoader(false);

  if (eventSource) {
    eventSource.close();
    eventSource = null;
  }
}

async function startJob() {
  resetUI();

  const profile = profileInput.value.trim();
  if (!profile) {
    alert("Bitte SimC Profil eingeben");
    return;
  }

  setStatus("Job wird erstellt …");
  showLoader(true);

  const res = await fetch("/api/jobs", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ profile })
  });

  if (!res.ok) {
    setStatus("Fehler beim Erstellen des Jobs");
    showLoader(false);
    return;
  }

  const { job_id } = await res.json();
  setStatus("Simulation läuft …");

  startStream(job_id);
}

function startStream(jobId) {
  eventSource = new EventSource(`/api/jobs/${jobId}/stream`);

  eventSource.addEventListener("output", (event) => {
    liveOutputEl.textContent += event.data + "\n";
    liveOutputEl.scrollTop = liveOutputEl.scrollHeight;
  });

  eventSource.addEventListener("status", async (event) => {
    if (event.data === "done") {
      setStatus("Simulation abgeschlossen ✅");
      showLoader(false);
      eventSource.close();
      await fetchAndRenderResult(jobId);
    }

    if (event.data === "failed") {
      setStatus("Simulation fehlgeschlagen ❌");
      showLoader(false);
      eventSource.close();
    }
  });

  eventSource.onerror = () => {
    setStatus("Stream unterbrochen");
    showLoader(false);
    eventSource.close();
  };
}

async function fetchAndRenderResult(jobId) {
  setStatus("Ergebnis wird verarbeitet …");

  const res = await fetch(`/api/jobs/${jobId}/result`);
  if (!res.ok) {
    setStatus("Kein Ergebnis gefunden");
    return;
  }

  const { output } = await res.json();

  // Fallback: roher Output
  resultOutputEl.textContent = output;

  // Parser aufrufen
  const parsedRes = await fetch("/api/parse-simc", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ output })
  });

  if (!parsedRes.ok) {
    console.warn("Parser fehlgeschlagen");
    return;
  }

  const parsed = await parsedRes.json();
  renderParsedResult(parsed);
}

function renderParsedResult(data) {
  resultOutputEl.innerHTML = `
    <h2>Simulation Ergebnis</h2>

    <div style="display:grid;grid-template-columns:repeat(3,1fr);gap:12px;">
      <div><strong>Spieler</strong><br>${data.summary.player || "-"}</div>
      <div><strong>Spec</strong><br>${data.summary.spec || "-"}</div>
      <div><strong>DPS</strong><br>${Math.round(data.summary.dps || 0).toLocaleString()}</div>
    </div>
  `;
}

startBtn.addEventListener("click", startJob);
