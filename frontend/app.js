async function run() {
  const res = await fetch("http://localhost:8080/api/jobs", {
    method: "POST",
    headers: {"Content-Type":"application/json"},
    body: JSON.stringify({profile: document.getElementById("profile").value})
  })
  const {job_id} = await res.json()
  poll(job_id)
}

async function poll(id) {
  const p = await fetch(`http://localhost:8080/api/jobs/${id}/progress`)
  document.getElementById("out").textContent = await p.text()
  setTimeout(()=>poll(id), 1000)
}