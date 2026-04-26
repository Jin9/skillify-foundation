# Pseudocode Generation Rules

## Core Principles
*   **"Antigravity" Philosophy**: Keep the logic lightweight, effortless to read, and free of language-specific boilerplate.
*   **Language-Agnostic**: Write logic that can be equally understood by a Python, Go, or TypeScript developer.
*   **Progressive Refinement**: Start with high-level intent, then drill down into detailed logical steps.

## Pseudocode Formatting Rules
Use clear, capitalized keywords for control flow. Keep variable assignments and function calls readable.

**Example Format:**
```text
FUNCTION ProcessData(records):
    IF records IS EMPTY THEN
        RETURN []
    END IF
    
    SET processed_results = []
    
    FOR EACH record IN records:
        TRY:
            SET cleaned_data = CALL Sanitize(record)
            IF cleaned_data.isValid THEN
                APPEND cleaned_data TO processed_results
            END IF
        CATCH Error:
            LOG "Failed to process record"
        END TRY
    END FOR
    
    RETURN processed_results
END FUNCTION
```
