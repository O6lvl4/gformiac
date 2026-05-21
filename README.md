# gformiac

[日本語](README.ja.md)

Google Forms Infrastructure as Code — Declaratively manage Google Forms from YAML definitions.

## Install

```bash
go install github.com/O6lvl4/gformiac@latest
```

## Quick Start

```yaml
# form.yaml
title: "Customer Satisfaction Survey"
description: "Please share your feedback."

items:
  - title: "Your name"
    type: short_answer
    required: true

  - title: "How satisfied are you?"
    type: choice
    required: true
    choice:
      type: radio
      options:
        - "Very satisfied"
        - "Satisfied"
        - "Neutral"
        - "Dissatisfied"
        - "Very dissatisfied"

  - title: "Any comments?"
    type: paragraph
```

```bash
gformiac plan                  # Preview changes (dry-run)
gformiac apply                 # Create or update the form
gformiac import <formID>       # Import an existing form to YAML
```

## Setup

1. Create a project in [Google Cloud Console](https://console.cloud.google.com/)
2. Enable the Google Forms API
3. Create an OAuth 2.0 Client ID and save it as `credentials.json`
4. On first `apply`, authenticate via browser — `token.json` is saved automatically

For service accounts, just set the `GOOGLE_APPLICATION_CREDENTIALS` environment variable.

## Commands

| Command | Description |
|---|---|
| `plan` | Preview diff between local definition and remote form |
| `apply` | Create or update the form |
| `import <formID>` | Import an existing form into a YAML definition |

### Flags

```
-f, --file string          Form definition file (default "form.yaml")
    --credentials string   OAuth2 credentials file (default "credentials.json")
    --token string         OAuth token file (default "token.json")
    --state string         State file (default "gformiac.state.json")
    --auto-approve         Skip confirmation prompt (for apply)
```

## Item Types

| type | Description | Extra fields |
|---|---|---|
| `short_answer` | Short text | — |
| `paragraph` | Long text | — |
| `choice` | Selection | `choice.type`: `radio` / `checkbox` / `dropdown` |
| `scale` | Linear scale | `scale.low`, `scale.high`, `scale.low_label`, `scale.high_label` |
| `date` | Date picker | — |
| `time` | Time picker | — |
| `page_break` | Section break | — |

## YAML Schema

```yaml
title: "Form title"                # required
description: "Form description"    # optional

items:
  - title: "Question text"         # required
    type: short_answer             # required
    description: "Help text"       # optional
    required: true                 # optional (default: false)

  - title: "Pick one"
    type: choice
    choice:
      type: radio                  # radio / checkbox / dropdown
      options:
        - "Option A"
        - "Option B"

  - title: "Rating"
    type: scale
    scale:
      low: 1
      high: 10
      low_label: "Not at all"
      high_label: "Absolutely"
```

## State Management

Running `apply` generates `gformiac.state.json`, which tracks the mapping between the local definition and the remote form/item IDs. This enables diff detection on subsequent `plan` / `apply` runs.

```
credentials.json        ← add to .gitignore
token.json              ← add to .gitignore
gformiac.state.json     ← add to .gitignore
```

## License

MIT
