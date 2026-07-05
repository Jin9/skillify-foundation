# GitHub - openai/skills: Skills Catalog for Codex · GitHub

Source: https://github.com/openai/skills
Accessed: 2026-04-26
Category: openai / official skills catalog
Provenance: normalized 2026-07-05 from raw HTML capture; original Accessed date preserved

## Why This Source Matters

OpenAI's official skills repository: real folder structures plus the .system/.curated/.experimental governance tiers worth borrowing. Tier 3 official examples.


openai   / ** skills ** Public

-   Notifications  You must be signed in to change notification settings
-   Fork 1.1k
-

 Star  17.5k

Branches Tags

Open more actions menu

## Folders and files

| Name | Name | Last commit message | Last commit date |
|---|---|---|---|

| ## Latest commit ## History 99 Commits 99 Commits |
| skills | skills |  |  |
| .gitignore | .gitignore |  |  |
| README.md | README.md |  |  |
| contributing.md | contributing.md |  |  |
|  |

## Repository files navigation

# Agent Skills

Agent Skills are folders of instructions, scripts, and resources that AI agents can discover and use to perform at specific tasks. Write once, use everywhere.

Codex uses skills to help package capabilities that teams and individuals can use to complete specific tasks in a repeatable way. This repository catalogs skills for use and distribution with Codex.

Learn more:

- [Using skills in Codex](https://developers.openai.com/codex/skills)
- [Create custom skills in Codex](https://developers.openai.com/codex/skills/create-skill)
- [Agent Skills open standard](https://agentskills.io)

## Installing a skill

Skills in `.system` are automatically installed in the latest version of Codex.

To install curated or experimental skills, you can use the `$skill-installer` inside Codex.

Curated skills can be installed by name (defaults to `skills/.curated`):

```
$skill-installer gh-address-comments

```

For experimental skills, specify the skill folder. For example:

```
$skill-installer install the create-plan skill from the .experimental folder

```

Or provide the GitHub directory URL:

```
$skill-installer install https://github.com/openai/skills/tree/main/skills/.experimental/create-plan

```

After installing a skill, restart Codex to pick up new skills.

## License

The license of an individual skill can be found directly inside the skill's directory inside the `LICENSE.txt` file.

## About

 Skills Catalog for Codex

### Resources

 Readme

### Contributing

 Contributing

###  Uh oh!

There was an error while loading. Please reload this page.

Activity

Custom properties

### Stars

**17.5k** stars

### Watchers

**101** watching

### Forks

**1.1k** forks

 Report repository

##  Releases

No releases published

##  Packages 0

###  Uh oh!

There was an error while loading. Please reload this page.

##  Contributors

-

-

-

###  Uh oh!

There was an error while loading. Please reload this page.

## Languages

-   Python 77.4%
-   JavaScript 15.3%
-   Shell 3.1%
-   Jupyter Notebook 1.4%
-   Swift 1.4%
-   PowerShell 1.4%
