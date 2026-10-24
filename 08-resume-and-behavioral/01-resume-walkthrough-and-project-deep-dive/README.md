# Resume Walkthrough and Project Deep Dive

[← Section index](../README.md) · Next: [STAR Stories](../02-star-stories/README.md)

**Interview weight: 1.4%**  ·  Revise in: 20 min

## In one line
"I am a Go backend engineer with four years of experience building a multi-cloud discovery platform and distributed services for 5G core networks."

## Must know
- Every resume line is a question. Be ready to go two levels deeper on each.
- Never quote a number you cannot defend. Use your real figures or say "roughly".
- Say "I" for what you did and "we" for what the team did. Interviewers listen for the difference.
- Know why you chose each technology and what you rejected.
- Know one failure or incident per project and what you changed after it.

## 60-second self introduction (fill the brackets)
> I am a backend engineer with four years of Go. Recently I worked on a multi-cloud infrastructure discovery platform that maps cloud resources into a graph. My part was [your part, e.g. ingestion, graph model, APIs]. Before that I built distributed backend services for 5G core networks, where [your focus, e.g. reliability, latency, messaging]. My stack is Go, REST/gRPC, Kafka/NATS/Pub/Sub, PostgreSQL, MongoDB, Neo4j, Redis and Kubernetes. I am looking for [what you want next] and [why this company].

Order: present, past, why here. Stop at 60 seconds and let them pick a thread.

## Project walkthrough skeleton
Use this for each project. One or two sentences per line.
1. **Problem:** what the system does and who uses it. Scale: [your number] accounts / requests / events.
2. **My role:** what I owned. Team size [your number].
3. **Design:** services, data stores, messaging, and why each. Draw it if asked.
4. **Hard part:** one concrete technical problem, the options I weighed, what I chose.
5. **Result:** what improved: [your number] latency / throughput / cost / incidents. If unmeasured, say so.

## Likely deep-dive questions
**Discovery platform**
- How do you collect resources from several clouds, and how do you handle API rate limits and retries?
- Why a graph database? What are the nodes and edges? Show a query.
- How do you keep the graph correct when resources change or are deleted?
- How do you make ingestion idempotent and handle duplicate messages?
- How do you scale collectors and bound concurrency? (see 02-concurrency/09, 11)

**5G core services**
- How did you get reliability: retries, timeouts, circuit breaking, health checks?
- How did you debug latency or memory problems? (pprof, see 03-runtime-internals/06)
- How is service-to-service traffic secured? (mTLS, JWT, Istio, see 07-backend-ecosystem/10)
- Why NATS or Kafka for this, and what delivery guarantee did you rely on?

**General**
- Why Redis, PostgreSQL, MongoDB: what goes where? (07-backend-ecosystem)
- How did you monitor it? (Prometheus, Grafana, ELK)
- What would you change if you rebuilt it today?

## Common traps
- Reciting technology names without a reason for each choice, or claiming team results as personal ones.
- Inventing a metric under pressure. Say "I do not have the exact number, roughly [x]".
- Keeping confidential details out: describe the design generically, no customer or company names.

## Interview Q&A

**Q1. Can you tell me about your current project?**

**A:** Use the skeleton: problem, my role, design, the hard part, result. Keep it under 2 minutes.

**Q2. What was the hardest bug or incident?**

**A:** Symptom, how you found the cause, the fix, what you changed to prevent it. Details in 02-star-stories.

**Q3. What would you do differently?**

**A:** Name one real design trade-off and what you would pick now and why.

**Q4. Can you tell me about yourself?**

**A:** Use the 60-second introduction: present, past, why here. Do not recite the whole resume.
