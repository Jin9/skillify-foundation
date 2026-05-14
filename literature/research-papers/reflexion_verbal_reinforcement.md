# Reflexion: Language Agents With Verbal Reinforcement Learning

Source: https://arxiv.org/abs/2303.11366  
Authors: Noah Shinn, Federico Cassano, Edward Berman, Ashwin Gopinath, Karthik Narasimhan, Shunyu Yao  
Submitted: 2023-03-20  
Accessed: 2026-05-14  
Category: Research paper / agent reflection and memory

## Why This Source Matters

Reflexion provides a research basis for post-task reflection, evaluation loops, and memory buffers. These ideas map directly to skill lifecycle practices: run the skill, inspect failures, capture lessons, and improve future executions without changing model weights.

## Core Data Points

- The paper proposes reinforcing language agents through verbal feedback rather than model fine-tuning.
- Agents reflect on feedback signals and store reflective text in episodic memory.
- Feedback can be scalar, natural language, external, or internally simulated.
- The approach is evaluated across sequential decision-making, coding, and reasoning tasks.
- The paper reports large improvements on HumanEval coding when reflection is used.
- The key mechanism is iterative: try, receive feedback, reflect, store the reflection, and use it on the next attempt.

## Design Implications For Skillify

- Skills should include explicit failure-capture points for tasks that require iteration.
- Post-run notes should be compact and actionable, not verbose transcripts.
- A skill's references can include known failure modes and repairs gathered from previous traces.
- Evaluation artifacts should be designed so the agent can turn failures into future improvements.
- Reflection should not replace deterministic tests; it should interpret failed tests and guide the next attempt.

## Skill Pattern Extracted

```text
1. Attempt the task using the current procedure.
2. Run an external or internal check.
3. Capture the failure signal.
4. Write a short reflection: root cause, missing context, next correction.
5. Retry with the reflection in context.
6. Promote stable lessons into the skill only after repeated evidence.
```

## Caution

Reflection can amplify a wrong lesson if the feedback signal is weak. Skillify should require evidence before converting a one-off reflection into permanent instructions.

