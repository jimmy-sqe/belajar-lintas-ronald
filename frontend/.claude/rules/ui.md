# UI & Component Rules

Detail for the UI directives summarized in `CLAUDE.md`. Covers the
Horizon component library, theming/providers, and the two sanctioned
non-Horizon visualization libraries (charts & maps).

## Horizon UI is the sole component library

- **NEVER** create a custom `Button`, `TextField`, `Select`, `Dialog`,
  `Table`, or any other component that already exists in Horizon.
- **NEVER** use raw `<button>`, `<input>`, `<select>`, `<table>`, or
  `<dialog>` HTML elements when a Horizon component covers the use case.
- **NEVER** install third-party UI libraries (Material UI, Chakra,
  Ant Design, shadcn/ui) — Horizon is the single source of truth.
  The ONLY sanctioned non-Horizon UI libraries are ECharts (charts)
  and Leaflet (maps) — see "Charts & maps" below.
- **ALWAYS** import from `@squantumengine/horizon` (named exports) or
  `@squantumengine/horizon/lib/*` for direct imports.
- **NEVER** use the `hz-` Tailwind prefix in OUR JSX — that prefix is
  internal to Horizon. Use project-owned classes for custom
  layout/spacing.
- Targeting Horizon-emitted classes in CSS selectors (e.g.,
  `.filters .hz-chip { cursor: pointer; }`) IS allowed — just not in
  `className=` props.

## Before building UI, read Horizon docs

For any new component or icon work:

- Read `.claude/docs/horizon-components.md` first for the full props
  reference. Do NOT guess prop names.
- For icons, verify the name against the Horizon icon catalog at
  https://sqehorizon.squantumengine.com/?path=/docs/foundation-icon--overview

## Project setup providers

Provider stack is wired in `src/app/layout.tsx` (PRESERVE on edits).
Horizon theming is applied via a `theme-hz-blue` CSS class on
`<body>`, NOT a JSX `<ThemeProvider>` wrapper. The Horizon stylesheet
is imported globally:

```tsx
import '@squantumengine/horizon/lib/index.css';
import { ToasterProvider } from '@squantumengine/horizon';

// in layout.tsx body className includes "theme-hz-blue":
<body className="... theme-hz-blue ...">
  <ToasterProvider>{children}</ToasterProvider>
</body>
```

Available theme class suffixes (replace `blue` in `theme-hz-<theme>`):
`blue`, `green`, `magenta`, `neutral`, `orange`, `purple`, `red`,
`teal`, `yellow`.

→ Full provider order: `docs/ARCHITECTURE.md#provider-stack`

## Charts & maps (the only sanctioned non-Horizon UI libs)

Horizon covers neither data visualization nor maps. These two
libraries are the ONLY approved exceptions to the third-party UI ban:

- **Charts / analytics → Apache ECharts.** Install `echarts` +
  `echarts-for-react`. NEVER use Recharts, Chart.js, Nivo, Victory,
  or D3-direct for standard charts.
- **Maps / geospatial → Leaflet.** Install `leaflet` +
  `react-leaflet` (+ `@types/leaflet`). NEVER use Google Maps JS,
  Mapbox GL, or OpenLayers.

Both touch `window`/DOM, so they are client-only in the App Router:

- Mark the wrapping component `'use client'` AND load it via
  `next/dynamic` with `{ ssr: false }` to avoid hydration and
  `window is not defined` errors.
- Import Leaflet's CSS once in the map component:
  `import 'leaflet/dist/leaflet.css';`.
- Keep each chart/map component self-contained inside its feature
  folder (`_components/`), preserving the prune-by-folder discipline
  (see `architecture.md` → route groups).

Theming: feed Horizon design tokens (colors from the active
`theme-hz-*` theme) into the ECharts `option` object instead of
hardcoding hex values, so charts track the selected theme.

```tsx
// FleetMap is client-only — render it without SSR.
const FleetMap = dynamic(() => import('./_components/FleetMap'), {
  ssr: false,
});
```

## Quick-lookup Horizon components

| If you need... | Use this import from `@squantumengine/horizon` |
|---|---|
| Any clickable action | `Button` (variant: primary / secondary / text) |
| Text input or textarea | `TextField` (set `multiline` for textarea) |
| Form-bound text input | `FormTextField` (React Hook Form integration) |
| Dropdown / select | `Select` or `SelectMenu` |
| Checkbox / multi-check | `Checkbox` / `Checkbox.Group` |
| Radio options | `Radio` / `Radio.Group` |
| On/off toggle | `Switch` |
| Search input | `SearchBar` |
| Date selection | `DatePicker` or `RangeDatePicker` |
| Time selection | `TimePicker` |
| Content container | `Card` (has `Card.Meta` sub-component) |
| Modal / dialog | `Dialog` + `DialogHeader` + `DialogBody` + `DialogFooter` |
| Data table | `Table` + `useTable` hook |
| Page navigation | `Pagination` |
| Tab panels | `Tabs` |
| Side navigation | `Sidebar` |
| Page header with back | `Header` |
| Alert / banner | `Info` (type: info / success / warning / error) |
| Status badge / tag | `Label` (type: success / danger / info / warning / default) |
| Filter chip | `Chip` |
| Tooltip / popup | `Popover` |
| Stepper / wizard | `Steps` |
| Selection list | `Listing` |
| Accordion | `Collapse` |
| Horizontal rule | `Divider` |
| Loading spinner | `Spinner` |
| Loading placeholder | `Skeleton` |
| Heading (h1–h6) | `Title` |
| Body text | `Paragraph` |
| Icon | `Icon` |
| Zoomable image | `ImageZoom` |
| Toast notification | `useToaster()` hook (needs `ToasterProvider`) |

(Full props reference: `.claude/docs/horizon-components.md`.)

## Figma-to-code workflow

1. Read the Figma design (via Figma MCP if available) to get
   component names and properties.
2. Read `.claude/docs/horizon-components.md` for the full props
   reference.
3. Map each Figma component to its Horizon equivalent using the
   quick-lookup table above.
4. Map Figma properties to Horizon props — variant names → `variant`,
   size tokens → `size`.
5. For layout/spacing, use the project's own Tailwind or scoped CSS
   — Horizon handles component internals only.
6. If no Horizon component matches, build a custom one BUT use
   Horizon primitives inside it.
