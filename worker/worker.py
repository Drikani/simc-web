import redis, subprocess, uuid, time, os

r = redis.Redis(host="redis", decode_responses=True)
DATA = "/data"

while True:
    job = r.blpop("jobs", timeout=0)
    job_id, simc_text = job[1].split("::", 1)

    simc_file = f"{DATA}/{job_id}.simc"
    xml_file = f"{DATA}/{job_id}.xml"

    with open(simc_file, "w") as f:
        f.write(simc_text)

    subprocess.run([
        "simc",
        simc_file,
        f"output={xml_file}",
        "xml=1"
    ], timeout=300)

    r.set(f"done:{job_id}", "1")