# ReAct: Synergizing Reasoning And Acting In Language Models

Source: https://arxiv.org/abs/2210.03629  
Authors: Shunyu Yao, Jeffrey Zhao, Dian Yu, Nan Du, Izhak Shafran, Karthik Narasimhan, Yuan Cao  
Submitted: 2022-10-06; last revised 2023-03-10  
Accessed: 2026-05-14  
Category: Research paper / agent reasoning and tool use

## Why This Source Matters

ReAct is foundational for agent workflow design because it frames agent performance as an interleaving of reasoning and environment actions. Skill files often encode this same pattern: reason about the task, call tools or inspect resources, update the plan, and continue.

## Core Data Points

- The paper studies language models that generate reasoning traces and task-specific actions in an interleaved trajectory.
- Reasoning helps the model maintain and revise plans, track exceptions, and interpret observations.
- Actions let the model gather external information or affect an environment.
- In knowledge tasks, tool interaction can reduce hallucination and error propagation compared with reasoning-only approaches.
- In interactive tasks such as ALFWorld and WebShop, ReAct outperformed imitation and reinforcement-learning baselines using only a small number of examples.
- The authors emphasize interpretability: a trajectory with thoughts, actions, and observations is easier to inspect than opaque final answers.

## Design Implications For Skillify

- Strong skills should specify observation/action loops where the task requires interaction:
  - inspect state,
  - decide the next action,
  - run a tool or read a file,
  - evaluate the result,
  - adjust plan.
- Avoid skills that jump directly from input to final output when external verification is required.
- Tool use instructions should say what evidence the agent should collect, not only which tool to call.
- Debugging skills should include exception handling and plan revision rules.
- Review or research skills should preserve traceability from claims to observations.

## Skill Pattern Extracted

```text
1. Form an initial plan from the user goal.
2. Take a concrete action that changes or observes the environment.
3. Read the observation.
4. Update the plan based on the observation.
5. Repeat until evidence is sufficient.
6. Return a final answer grounded in the observed results.
```

