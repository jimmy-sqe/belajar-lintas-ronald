# Specs

Feature specifications live here, organized by quarter and feature slug:

```
specs/
└── <quarter>/             # e.g., 2026-Q2
    └── <feature-slug>/    # e.g., voucher-validation
        ├── PRD.md          # Product Requirements (from /write-prd)
        ├── tech-spec.md    # Technical Design (from /design-tech-spec)
        ├── api-contract.yaml   # API contract (from /design-tech-plan)
        ├── fe-blueprint.md     # FE engineering implementation guide (from /design-tech-plan)
        ├── qa-handles.yaml     # FE↔QA data-testid contract (from /design-tech-plan)
        ├── test-plan.md    # Test plan (from /generate-test-plan)
        ├── qa-report-*.md  # QA reports (from /qa-validate)
        └── .qa-state.yaml  # internal state
```

Run `/write-prd <feature-name>` from the lintas marketplace to create the first PRD.
