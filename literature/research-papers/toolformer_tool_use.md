# Toolformer: Language Models Can Teach Themselves To Use Tools

Source: https://arxiv.org/abs/2302.04761  
Authors: Timo Schick, Jane Dwivedi-Yu, Roberto Dessi, Roberta Raileanu, Maria Lomeli, Luke Zettlemoyer, Nicola Cancedda, Thomas Scialom  
Submitted: 2023-02-09  
Accessed: 2026-05-14  
Category: Research paper / tool use

## Why This Source Matters

Toolformer helps explain why tool-use guidance belongs in skills. A model may know a lot, but external tools handle arithmetic, lookup, execution, search, translation, calendars, and other functions more reliably when the agent knows when and how to call them.

## Core Data Points

- The paper trains language models to decide which APIs to call, when to call them, what arguments to pass, and how to incorporate results.
- Training uses self-supervision from a small number of demonstrations per API.
- Tool categories include calculator, question-answering, search, translation, and calendar.
- Tool-augmented models improve zero-shot performance on several downstream tasks.
- The paper's main lesson for agent design is that tool use is not merely access to tools; it is selection, argument construction, and result integration.

## Design Implications For Skillify

- Skills should say when a tool is appropriate, what inputs it expects, and how to interpret its outputs.
- Tool guidance should include refusal or fallback conditions when a tool is unavailable or risky.
- Scripts are best when task reliability depends on exact computation, parsing, validation, or environment checks.
- A skill should not tell the agent to use tools generically; it should define a decision rule.
- Output handling is part of tool use. A script returning JSON should be paired with instructions for interpreting the fields.

## Skill Pattern Extracted

```text
For each tool-supported step:
- trigger condition,
- required inputs,
- command or API shape,
- expected output,
- validation rule,
- fallback behavior.
```

