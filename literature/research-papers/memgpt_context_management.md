# MemGPT: Towards LLMs As Operating Systems

Source: https://arxiv.org/abs/2310.08560  
Authors: Charles Packer, Sarah Wooders, Kevin Lin, Vivian Fang, Shishir G. Patil, Ion Stoica, Joseph E. Gonzalez  
Submitted: 2023-10-12  
Accessed: 2026-05-14  
Category: Research paper / context management and memory

## Why This Source Matters

MemGPT gives a conceptual foundation for progressive disclosure. Skills, references, scripts, and memories are practical ways to manage limited context by keeping only the most relevant instructions in the active window.

## Core Data Points

- The paper frames limited context windows as a memory-management problem.
- It proposes virtual context management inspired by hierarchical memory in operating systems.
- The system manages tiers of information so an LLM can work with information larger than the active context window.
- It also uses interrupts to manage control flow between the system and user.
- Evaluated domains include long-document analysis and multi-session chat.
- The work supports the idea that agents need explicit mechanisms for loading, storing, and evicting context.

## Design Implications For Skillify

- Progressive disclosure is not just a style preference; it is a context-management strategy.
- `SKILL.md` should contain only what must be active after trigger.
- Large references should be loaded only when task conditions require them.
- Memories and postmortems should be compact and indexed by task, not dumped wholesale into always-on instructions.
- Skill bodies should include "load this reference only when..." guidance.
- Validation and scripts reduce context pressure by moving deterministic work out of natural-language instructions.

## Skill Pattern Extracted

```text
Hot context:
- name
- description
- minimal workflow

Warm context:
- focused reference files
- examples
- templates

Cold context:
- historical traces
- large docs
- generated artifacts
- raw source captures
```

