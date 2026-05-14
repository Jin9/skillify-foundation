# SWE-agent: Agent-Computer Interfaces Enable Automated Software Engineering

Source: https://arxiv.org/abs/2405.15793  
Authors: John Yang, Carlos E. Jimenez, Alexander Wettig, Kilian Lieret, Shunyu Yao, Karthik Narasimhan, Ofir Press  
Submitted: 2024-05-06; last revised 2024-11-11  
Accessed: 2026-05-14  
Category: Research paper / coding agents and tool interfaces

## Why This Source Matters

SWE-agent argues that language-model agents are end users of software tools and benefit from interfaces designed for their capabilities. This is directly relevant to skills because a good skill often defines an agent-facing interface: commands, scripts, file layouts, test loops, and observation formats.

## Core Data Points

- The paper studies automated software engineering with language-model agents.
- It introduces an agent-computer interface designed for repository navigation, editing, and program execution.
- The authors argue that interface design affects agent behavior and performance.
- SWE-agent is evaluated on SWE-bench and HumanEvalFix.
- The paper reports improved results over non-interactive language-model baselines.
- Important interface features include tools for editing files, navigating code, and running tests.

## Design Implications For Skillify

- Scripts inside skills should be agent-friendly, not just human-friendly.
- A skill should define what observations the agent should expect from commands and how to act on failures.
- Tool wrappers should produce concise, structured output when possible.
- Skills for software engineering should include repository navigation strategy and validation loops.
- A good `SKILL.md` can function as an agent-computer interface specification for a narrow workflow.

## Skill Pattern Extracted

```text
1. Define the task interface: files, commands, inputs, outputs.
2. Provide navigation rules for finding relevant code.
3. Provide edit rules for making bounded changes.
4. Provide validation commands.
5. Provide failure interpretation rules.
6. Return a concise summary with changed files and remaining risks.
```

