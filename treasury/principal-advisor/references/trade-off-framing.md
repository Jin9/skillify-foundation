# Trade-off framing — options, doors, and over-engineering

How to structure an `advise` answer so it lands as a decision, not a menu.

## Options discipline

- Always at least two options. An option = name + pros + cons + cost and maintenance posture. Two strong options beat four thin ones.
- Design order: business need and volume first → non-functional requirements → boundaries and ownership → only then stack, service count, and standards. A stack argument that starts at the stack is upside down.
- Include the do-nothing or boring-baseline option when it is defensible; it often wins.

## Evaluation axes

State the position on at least two of these per decision — pick the axes the decision actually moves:

- complexity ↔ maintainability
- cost ↔ performance
- speed ↔ safety
- coupling ↔ flexibility

## Door tagging

- One-way door (stack choice, data-migration strategy, structural split, vendor commitment): raise the rigor — decide-then-explain, name the rollback or revisit trigger, do not rush.
- Two-way door: decide fast and adjust freely; say explicitly that it is reversible so the user stops over-deliberating.
- Compliance and regulatory constraints are one-way doors, not options.

## Over-engineering signals

Say it plainly when any of these appear:

- The solution process is more complex than the problem size warrants.
- The design takes five minutes to explain — that is the smell itself.
- Microservices without answers to the three questions: which business scope does each serve? is separation a priority now? can this team maintain the split?
- A resume-driven or new-because-new technology choice.

## Transition cost

An option's price includes the journey, not just the destination. For any change to a running system, weigh: the migration path, the coexistence period while both systems are live, rollback at each stage, data movement, and the exit cost if the new choice fails. An option that wins on the target state and loses on the transition often loses overall.

## Criteria, not crowns

- Never a bare "X is best." Emit the conditions under which each option wins ("A while the team is under five and latency is soft; B once you need independent deploys").
- Then still recommend exactly one option for this user's stated constraints — an advisor recommends; a menu dodges.
- Name the flip criteria: the observable condition that should change the recommendation.
- When comparing products, models, or vendors that move fast, cite specifics as dated examples and anchor the advice in tiers and criteria that outlive the examples.

## Rendering in chat

Compact table for three or more options; prose for two. One recommendation line at the end, bolded, carrying its strongest reason and its flip criterion.
