---
name: thai-translator
description: >
  Translates English documents into structurally identical bilingual (Thai/English) documents. Use when the user asks to 'translate to thai', 'mirror this document in thai', or 'translate to bilingual thai'. Do NOT use for translating to languages other than Thai.
---

# Thai Translator

## Purpose

Translate English documents into structurally identical bilingual (Thai and English) documents. It maintains all formatting while keeping headings, technical terms, and key concepts in both languages.

## When to use this skill

- Use when: The user asks to "translate to thai" or "translate to bilingual thai" for an English document.
- Use when: The user asks to "mirror this document in thai" or "convert to thai language".
- Do NOT use when: The user asks to translate to any language other than Thai, or explicitly asks for an English-only or Thai-only version without bilingual terms.

## Workflow

1. Read the user's request and identify the target English document.
2. Read the source document to understand its structure, formatting (markdown, HTML, etc.), and tone.
3. Translate the document section by section into Thai, using a bilingual approach:
   - **Bilingual Headings:** Keep headings in both languages (e.g., `# [Thai Translation] ([English Original])`).
   - **Bilingual Terms:** Domain-specific terminology, proper nouns, and key concepts must be translated into Thai followed by the original English term in parentheses (e.g., "การเรียนรู้ของเครื่อง (Machine Learning)").
   - **Formatting Parity:** All markdown headings, lists, tables, and bold/italic formatting must be perfectly preserved.
   - **Untranslated Elements:** Any code blocks, filenames, URLs, or variable names are left strictly in English.
4. Save the translated document to the requested output path.
5. Validate the structural parity between the source and the translated document.

## Output Format

The output must be a single translated file with the same format as the source (e.g., Markdown, text).
- Default naming convention: Append `-th` to the original filename before the extension (e.g., `document-th.md`), unless the user specifies a different path.

## Constraints

- DO NOT change the structure of the document (e.g., do not combine paragraphs or alter list hierarchies).
- DO NOT fully translate technical terms without keeping the original English term alongside them in parentheses.
- DO NOT translate code blocks, URLs, or system-specific identifiers.
- DO NOT add commentary, notes, or preambles outside of the translated text itself.

## Examples

### Example 1: Translate a technical markdown document

**User says**: "Translate README.md to thai"

**Action**:
1. Read `README.md`.
2. Translate the contents maintaining markdown features and applying bilingual rules to headings and terms.
3. Save the result as `README-th.md`.

**Result**: A fully translated `README-th.md` with identical structure, containing bilingual terms like "ฐานข้อมูล (Database)".

## Validation

- [ ] Frontmatter name matches the folder.
- [ ] Description includes concrete trigger phrases.
- [ ] Workflow and output contract are explicit.
- [ ] The translated document has the exact same structural elements as the original document.
- [ ] The translated document correctly utilizes the bilingual format for headings and technical terms.
