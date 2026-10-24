# STAR Stories

[← Section index](../README.md) · Prev: [Resume Walkthrough and Project Deep Dive](../01-resume-walkthrough-and-project-deep-dive/README.md) · Next: [Interview Rounds and Questions to Ask](../03-interview-rounds-and-questions-to-ask/README.md)

**Interview weight: 1.0%**  ·  Revise in: 15 min

## In one line
Prepare three true stories in Situation, Task, Action, Result form and reuse them for most behavioral questions.

## Must know
- **S**ituation (1 sentence), **T**ask (your goal), **A**ction (what *you* did, 60% of the time), **R**esult (outcome plus what you learned).
- Use only real facts. If you did not measure a result, say what you observed instead.
- One story can answer many questions: pick the angle the question asks for.
- Keep each story to 2 minutes. Practise out loud.
- Keep it generic: no customer or company names, nothing confidential.

## Story 1: Performance or scale problem
**S:** A [scan / batch job / API] was too slow for [reason: more data, more accounts, more traffic].
**T:** Bring it down to [target] without losing correctness.
**A:** I measured first ([profiling / metrics / logs]), found [the bottleneck], then [your change, e.g. bounded concurrency, batching, caching, skipping unchanged writes]. I tested with [how].
**R:** [your number] before and after, or "noticeably faster, measured by [what]". Lesson: [what you would measure earlier next time].

## Story 2: Production incident or bug
**S:** In production, [symptom: memory growth, restarts, latency spike, stuck consumers].
**T:** Find the cause and fix it with minimal user impact.
**A:** I [checked metrics and logs, took a pprof snapshot or traces], found [root cause, e.g. a goroutine never cancelled, a missing timeout], fixed it by [change], and added [alert / test / runbook].
**R:** [what stopped happening and for how long]. Lesson: [what you changed in the process].

## Story 3: Disagreement, ownership, or mentoring
**S:** [A teammate or lead] preferred [option A]; I thought [option B] was better because [reason].
**T:** Reach a decision that was good for the product and the team.
**A:** I [wrote down trade-offs / ran a small test / asked for data], listened to their concerns, and [we chose X]. If it went against me, I [supported it and checked the result].
**R:** [outcome]. Lesson: [what you learned about communication or decisions].

## Common traps
- Saying "we" for everything, so the interviewer cannot tell what you did.
- Skipping the result or the lesson.
- Blaming a teammate. Describe the problem, not the person.
- Choosing a story where you had no real decision to make.

## Interview Q&A

**Q1. Can you tell me about a time you debugged a hard production issue?**

**A:** Use Story 2. Stress how you found the cause, not just the fix.

**Q2. Can you tell me about a time you disagreed with someone?**

**A:** Use Story 3. Show you used data, stayed respectful and committed to the final call.

**Q3. Why are you leaving your current company?**

**A:** Say what you want next ([growth, scale, domain, ownership]), not what is wrong there. Stay positive about your current team.

**Q4. What is a mistake you made, or your weakness? And your strength?**

**A:** Real and small, plus the fix and what you do differently now. Strength: one skill, proved by Story 1 or 2.

**Q5. What is the biggest technical challenge you have solved?**

**A:** Use Story 1 or 2. Pick one with a real trade-off, say what you measured, and end with the lesson.
