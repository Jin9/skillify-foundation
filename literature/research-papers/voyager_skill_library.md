# Voyager: Open-Ended Embodied Agent With A Skill Library

Source: https://arxiv.org/abs/2305.16291  
Authors: Guanzhi Wang, Yuqi Xie, Yunfan Jiang, Ajay Mandlekar, Chaowei Xiao, Yuke Zhu, Linxi Fan, Anima Anandkumar  
Submitted: 2023-05-25; last revised 2023-10-19  
Accessed: 2026-05-14  
Category: Research paper / lifelong learning and skill libraries

## Why This Source Matters

Voyager is one of the clearest research examples of a skill library compounding an agent's capabilities over time. Although it is an embodied Minecraft agent rather than a coding assistant, its architecture maps well to reusable `SKILL.md` folders: store successful procedures, retrieve them later, and compose them for harder tasks.

## Core Data Points

- Voyager is an LLM-powered lifelong learning agent in Minecraft.
- It combines three main components:
  - an automatic curriculum for exploration,
  - an expanding library of executable skills,
  - iterative prompting with environment feedback, execution errors, and self-verification.
- Its skills are executable code that represent temporally extended behaviors.
- The skill library supports retrieval and reuse in new situations.
- The paper argues that learned skills are interpretable, compositional, and help reduce catastrophic forgetting.
- Reported results include more unique items, longer travel distances, and faster milestone unlocking than prior baselines.
- The agent can use learned skills in a new world to solve novel tasks.

## Design Implications For Skillify

- A skill library should optimize for retrieval, reuse, and composition, not just documentation.
- Skills should encode stable procedures that have succeeded in real traces.
- Executable scripts inside skills are analogous to Voyager's executable skill programs.
- Skill descriptions should include enough semantic signal for retrieval.
- Skill maintenance should include pruning or refactoring old skills that overlap or become stale.
- Compositionality requires narrow skills with clean boundaries.

## Skill Pattern Extracted

```text
1. Identify a repeated task.
2. Capture the successful procedure as a reusable skill.
3. Store the skill with a clear retrieval description.
4. Reuse the skill in new tasks.
5. Update the skill only when new evidence improves generality.
6. Compose multiple narrow skills for larger workflows.
```

