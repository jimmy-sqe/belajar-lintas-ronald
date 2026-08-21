# Horizon Component API — Full Reference

Complete props reference for every component in `@squantumengine/horizon`.
All components are named exports: `import { Button, TextField, ... } from '@squantumengine/horizon';`

---

## Button

```tsx
<Button variant="primary" size="md">Click me</Button>
<Button variant="secondary" size="sm" full>Full width</Button>
<Button variant="text" loading>Loading...</Button>
<Button variant="primary" inverted>Inverted</Button>
```

| Prop        | Type                                 | Default     | Description                   |
| ----------- | ------------------------------------ | ----------- | ----------------------------- |
| `variant`   | `'primary' \| 'secondary' \| 'text'` | `'primary'` | Visual style                  |
| `size`      | `'sm' \| 'md' \| 'lg'`               | `'md'`      | Button size (h-8, h-10, h-12) |
| `full`      | `boolean`                            | `false`     | Full width                    |
| `disabled`  | `boolean`                            | `false`     | Disabled state                |
| `loading`   | `boolean`                            | `false`     | Shows spinner, hides label    |
| `inverted`  | `boolean`                            | `false`     | For dark backgrounds          |
| `className` | `string`                             | —           | Additional CSS classes        |

Extends all native `ButtonHTMLAttributes`.

---

## TextField

```tsx
<TextField label="Email" placeholder="you@example.com" />
<TextField label="Bio" multiline rows={4} />
<TextField label="Name" clearable errorMessage="Required" />
<TextField label="Amount" prefix="Rp" suffix=".00" />
```

| Prop           | Type        | Default | Description                       |
| -------------- | ----------- | ------- | --------------------------------- |
| `label`        | `string`    | —       | Floating label text               |
| `animateLabel` | `boolean`   | —       | Enable label animation            |
| `multiline`    | `boolean`   | `false` | Renders as `<textarea>`           |
| `rows`         | `number`    | —       | Textarea rows (when multiline)    |
| `clearable`    | `boolean`   | —       | Show clear button                 |
| `full`         | `boolean`   | —       | Full width                        |
| `prefix`       | `ReactNode` | —       | Content before input              |
| `suffix`       | `ReactNode` | —       | Content after input               |
| `iconPrefix`   | `ReactNode` | —       | Icon before input                 |
| `description`  | `string`    | —       | Helper text below                 |
| `errorMessage` | `ReactNode` | —       | Error text (replaces description) |
| `showCounter`  | `boolean`   | —       | Character counter                 |
| `inputRef`     | `RefObject` | —       | Ref to the input/textarea         |

### FormTextField (React Hook Form)

```tsx
import { FormTextField } from '@squantumengine/horizon';

<FormTextField
  controller={{ control, name: 'email' }}
  label="Email"
  errorMessage={errors.email?.message}
  valueTransformer={(val) => val.trim()}
/>;
```

Additional props: `controller` (UseControllerProps), `valueTransformer`.

---

## Select

```tsx
const options = [
  { value: 'id', label: 'Indonesia' },
  { value: 'my', label: 'Malaysia' },
];

<Select options={options} label="Country" placeholder="Choose..." onChange={(val) => {}} />
<Select options={options} isMultiple isSearchable label="Tags" />
```

| Prop             | Type                                       | Default      | Description                                              |
| ---------------- | ------------------------------------------ | ------------ | -------------------------------------------------------- |
| `options`        | `SelectOption[]`                           | **required** | `{ value: string \| number, label: string, data?: any }` |
| `value`          | `string \| number \| (string \| number)[]` | —            | Controlled value                                         |
| `label`          | `string`                                   | —            | Floating label                                           |
| `placeholder`    | `string`                                   | —            | Placeholder text                                         |
| `isMultiple`     | `boolean`                                  | —            | Multi-select mode                                        |
| `isSearchable`   | `boolean`                                  | —            | Enable search/filter                                     |
| `isDisabled`     | `boolean`                                  | —            | Disabled state                                           |
| `isError`        | `boolean`                                  | —            | Error visual state                                       |
| `errorMessage`   | `string`                                   | —            | Error text                                               |
| `footerMessage`  | `string`                                   | —            | Helper text below                                        |
| `prefixIcon`     | `ReactNode`                                | —            | Icon inside trigger                                      |
| `full`           | `boolean`                                  | —            | Full width                                               |
| `isLoading`      | `boolean`                                  | —            | Loading state                                            |
| `filterOptions`  | `boolean`                                  | —            | Client-side filter                                       |
| `onChange`       | `(value, data?) => void`                   | —            | Change handler                                           |
| `renderOptions`  | `(params) => ReactNode`                    | —            | Custom option renderer                                   |
| `onChangeSearch` | `(value: string) => void`                  | —            | Server-side search handler                               |
| `valueSearch`    | `string`                                   | —            | Controlled search value                                  |
| `className`      | `string`                                   | —            | Additional CSS classes                                   |
| `zIndex`         | `number`                                   | —            | Popup z-index                                            |

---

## SelectMenu

Simpler select alternative. Same `SelectOption` type.

```tsx
<SelectMenu options={options} label="Role" placeholder="Select role" />
<SelectMenu options={options} multiple showSearch />
```

| Prop                  | Type                                              | Default      | Description           |
| --------------------- | ------------------------------------------------- | ------------ | --------------------- |
| `options`             | `SelectOption[]`                                  | **required** | Options list          |
| `value`               | `SelectOption \| SelectOption[]`                  | —            | Controlled value      |
| `label`               | `string`                                          | —            | Label text            |
| `placeholder`         | `string`                                          | —            | Placeholder           |
| `multiple`            | `boolean`                                         | —            | Multi-select          |
| `showSearch`          | `boolean`                                         | —            | Search within options |
| `disabled`            | `boolean`                                         | —            | Disabled state        |
| `errorMessage`        | `string`                                          | —            | Error text            |
| `footerText`          | `string`                                          | —            | Helper text           |
| `onChange`            | `(value: SelectOption \| SelectOption[]) => void` | —            | Change handler        |
| `customRenderOptions` | `(options?) => ReactNode`                         | —            | Custom renderer       |

---

## Checkbox

```tsx
<Checkbox>Remember me</Checkbox>
<Checkbox checked indeterminate>Select all</Checkbox>

<Checkbox.Group
  options={['Option A', 'Option B', 'Option C']}
  defaultValue={['Option A']}
  onChange={(checked) => {}}
/>
```

| Prop            | Type               | Default | Description              |
| --------------- | ------------------ | ------- | ------------------------ |
| `checked`       | `boolean`          | —       | Controlled checked state |
| `indeterminate` | `boolean`          | —       | Indeterminate state      |
| `value`         | `string \| number` | —       | Value for group usage    |

**Checkbox.Group**: `options`, `value`, `defaultValue`, `disabled`, `onChange`.

---

## Radio

```tsx
<Radio.Group options={['Monthly', 'Yearly']} defaultValue="Monthly" onChange={(e) => {}} />
```

**Radio.Group**: `options`, `value`, `defaultValue`, `disabled`, `onChange`.

---

## Switch

```tsx
<Switch>Enable notifications</Switch>
```

Extends native input attributes. Toggle on/off control.

---

## SearchBar

```tsx
<SearchBar placeholder="Search..." size="default" onChange={(value) => {}} />
<SearchBar size="compact" searchIconPosition="end" />
```

| Prop                 | Type                      | Default     | Description                               |
| -------------------- | ------------------------- | ----------- | ----------------------------------------- |
| `size`               | `'default' \| 'compact'`  | `'default'` | Search bar height                         |
| `placeholder`        | `string`                  | —           | Placeholder text                          |
| `searchIconPosition` | `'start' \| 'end'`        | `'start'`   | Icon placement                            |
| `onChange`           | `(value: string) => void` | —           | Change handler                            |
| `debounceDuration`   | `number`                  | —           | Debounce delay in ms (DebouncedSearchBar) |
| `value`              | `string`                  | —           | Controlled value                          |
| `defaultValue`       | `string`                  | —           | Initial value                             |

---

## DatePicker

```tsx
<DatePicker placeholder="Select date" format="DD/MM/YYYY" onChange={(dateStr) => {}} />
<DatePicker showTime format="DD/MM/YYYY HH:mm" />
<DatePicker picker="month" />
```

| Prop                | Type                                   | Default | Description            |
| ------------------- | -------------------------------------- | ------- | ---------------------- |
| `value`             | `string \| Date \| Dayjs \| null`      | —       | Controlled value       |
| `format`            | `string`                               | —       | Display format (dayjs) |
| `picker`            | `'year' \| 'month' \| 'day'`           | `'day'` | Picker granularity     |
| `showTime`          | `boolean`                              | —       | Include time picker    |
| `locale`            | `'en' \| 'id'`                         | —       | Localization           |
| `disabled`          | `boolean`                              | —       | Disabled state         |
| `disabledDate`      | `(date: Dayjs) => boolean`             | —       | Disable specific dates |
| `onChange`          | `(dateString: string \| null) => void` | —       | Change handler         |
| `placeholder`       | `string`                               | —       | Placeholder text       |
| `showIcon`          | `boolean`                              | —       | Show calendar icon     |
| `isShowTodayButton` | `boolean`                              | —       | Show today shortcut    |

---

## RangeDatePicker

```tsx
<RangeDatePicker placeholder="Select range" showPreset onChange={(dates, dateStrings) => {}} />
```

| Prop           | Type                             | Default | Description            |
| -------------- | -------------------------------- | ------- | ---------------------- |
| `value`        | `[Dayjs \| null, Dayjs \| null]` | —       | Controlled range       |
| `locale`       | `'en' \| 'id'`                   | —       | Localization           |
| `format`       | `string`                         | —       | Display format         |
| `showPreset`   | `boolean`                        | —       | Show preset ranges     |
| `disabledDate` | `(date: Dayjs) => boolean`       | —       | Disable specific dates |
| `onChange`     | `(dates, dateStrings) => void`   | —       | Change handler         |
| `onOk`         | `(data?) => void`                | —       | Confirm handler        |
| `disabled`     | `boolean`                        | —       | Disabled state         |

---

## TimePicker

```tsx
<TimePicker format="HH:mm" placeholder="Select time" onChange={(dayjs) => {}} />
```

| Prop           | Type                        | Default      | Description        |
| -------------- | --------------------------- | ------------ | ------------------ |
| `value`        | `Dayjs \| null`             | —            | Controlled value   |
| `format`       | `string`                    | **required** | Time format        |
| `disabled`     | `boolean`                   | —            | Disabled state     |
| `full`         | `boolean`                   | —            | Full width         |
| `onChange`     | `(date: Dayjs) => void`     | —            | Change handler     |
| `onOpenChange` | `(isOpen: boolean) => void` | —            | Open state handler |

---

## Card

```tsx
<Card size="md" title="Card Title" extra={<Button variant="text">Action</Button>}>
  Card content here
</Card>
<Card size="sm" borderless cover={<img src="..." />} footer={<span>Footer</span>} />
```

| Prop             | Type                   | Default      | Description           |
| ---------------- | ---------------------- | ------------ | --------------------- |
| `size`           | `'sm' \| 'md' \| 'lg'` | **required** | Card padding size     |
| `title`          | `string \| ReactNode`  | —            | Card header title     |
| `extra`          | `ReactNode`            | —            | Top-right action area |
| `cover`          | `ReactNode`            | —            | Cover image/content   |
| `footer`         | `ReactNode`            | —            | Footer content        |
| `borderless`     | `boolean`              | —            | Remove border         |
| `isHoverOutline` | `boolean`              | —            | Outline on hover      |
| `onClick`        | `() => void`           | —            | Clickable card        |

**Card.Meta**: `<Card.Meta title="..." description="..." />`

---

## Dialog

Compound component. Always compose with sub-components.

```tsx
import { Dialog, DialogHeader, DialogBody, DialogFooter, Button } from '@squantumengine/horizon';

<Dialog open={isOpen} onClose={() => setIsOpen(false)}>
  <DialogHeader>Confirm Action</DialogHeader>
  <DialogBody>Are you sure?</DialogBody>
  <DialogFooter>
    <Button variant="secondary" onClick={() => setIsOpen(false)}>
      Cancel
    </Button>
    <Button onClick={handleConfirm}>Confirm</Button>
  </DialogFooter>
</Dialog>;
```

| Prop                  | Type         | Default      | Description             |
| --------------------- | ------------ | ------------ | ----------------------- |
| `open`                | `boolean`    | **required** | Controls visibility     |
| `onClose`             | `() => void` | **required** | Close handler           |
| `closeOnClickOutside` | `boolean`    | —            | Close on backdrop click |
| `hideCloseBtn`        | `boolean`    | —            | Hide the X button       |
| `className`           | `string`     | —            | Dialog wrapper class    |

`DialogHeader`, `DialogBody`, `DialogFooter` each accept `children` and `className`.

---

## Table

Uses TanStack React Table v8. Create instance with `useTable` hook.

```tsx
import { Table, useTable, Pagination } from '@squantumengine/horizon';
import type { TableColumnTypeDef } from '@squantumengine/horizon';

const columns: TableColumnTypeDef<DataType> = [
  { accessorKey: 'name', header: 'Name', width: 200 },
  { accessorKey: 'email', header: 'Email' },
  {
    accessorKey: 'status',
    header: 'Status',
    horizontalAlign: 'center',
    cell: ({ getValue }) => <Label type="success" label={getValue()} />,
  },
];

const { table } = useTable({ data, columns });

<Table table={table} data={data} isStickyHeader />
<Pagination total={100} current={page} onChange={setPage} />
```

| Prop                  | Type            | Default      | Description                             |
| --------------------- | --------------- | ------------ | --------------------------------------- |
| `table`               | `Table<T>`      | **required** | TanStack table instance from `useTable` |
| `data`                | `T[]`           | **required** | Data array                              |
| `isStickyHeader`      | `boolean`       | —            | Sticky header                           |
| `isLoading`           | `boolean`       | —            | Loading skeleton                        |
| `freezeLeftColCount`  | `number`        | —            | Frozen left columns                     |
| `freezeRightColCount` | `number`        | —            | Frozen right columns                    |
| `onBodyRowClick`      | `(row) => void` | —            | Row click handler                       |

Column extensions: `width`, `horizontalAlign` (`'left' | 'center' | 'right'`), `headerVerticalAlign` (`'top' | 'middle' | 'bottom'`).

---

## Pagination

```tsx
<Pagination
  total={200}
  current={1}
  pageSize={10}
  showSizeChanger
  onChange={setPage}
  onPageSizeChange={setPageSize}
/>
```

| Prop               | Type                     | Default             | Description                |
| ------------------ | ------------------------ | ------------------- | -------------------------- |
| `total`            | `number`                 | **required**        | Total items                |
| `current`          | `number`                 | —                   | Current page               |
| `pageSize`         | `number`                 | `10`                | Items per page             |
| `pageSizeOptions`  | `number[]`               | `[10, 20, 50, 100]` | Page size dropdown options |
| `showSizeChanger`  | `boolean`                | —                   | Show page size selector    |
| `onChange`         | `(page: number) => void` | —                   | Page change handler        |
| `onPageSizeChange` | `(size: number) => void` | —                   | Size change handler        |

---

## Tabs

```tsx
<Tabs
  size="md"
  items={[
    { key: 'tab1', label: 'Overview', children: <Overview /> },
    { key: 'tab2', label: 'Details', children: <Details /> },
    { key: 'tab3', label: 'Disabled', children: null, disabled: true }
  ]}
  onChange={(key) => {}}
/>
```

| Prop               | Type                                     | Default      | Description                           |
| ------------------ | ---------------------------------------- | ------------ | ------------------------------------- |
| `items`            | `TabItem[]`                              | **required** | `{ key, label, children, disabled? }` |
| `size`             | `'sm' \| 'md' \| 'lg'`                   | —            | Tab size                              |
| `defaultActiveKey` | `string`                                 | —            | Initially active tab                  |
| `tabPosition`      | `'top' \| 'bottom' \| 'left' \| 'right'` | `'top'`      | Tab bar position                      |
| `centered`         | `boolean`                                | —            | Center tab labels                     |
| `onChange`         | `(activeKey: string) => void`            | —            | Tab change handler                    |

---

## Sidebar

```tsx
<Sidebar
  items={[
    { id: 'dashboard', text: 'Dashboard', icon: <Icon name="home" /> },
    {
      id: 'users',
      text: 'Users',
      icon: <Icon name="users" />,
      children: [
        { id: 'list', text: 'All Users' },
        { id: 'add', text: 'Add User' }
      ]
    }
  ]}
  selectedMenu="dashboard"
  collapsible
  onClick={(item) => navigate(item.url)}
/>
```

| Prop               | Type                          | Default      | Description                           |
| ------------------ | ----------------------------- | ------------ | ------------------------------------- |
| `items`            | `SidebarItemProps[]`          | **required** | `{ id, text, icon, children?, url? }` |
| `selectedMenu`     | `string`                      | —            | Active item ID                        |
| `collapsible`      | `boolean`                     | —            | Allow collapse/expand                 |
| `expanded`         | `boolean`                     | —            | Controlled expanded state             |
| `footer`           | `ReactNode`                   | —            | Footer content                        |
| `width`            | `number`                      | —            | Sidebar width                         |
| `defaultOpenMenus` | `string[]`                    | —            | Initially open submenus               |
| `onClick`          | `(item) => void`              | —            | Item click handler                    |
| `onExpanded`       | `(expanded: boolean) => void` | —            | Expand/collapse handler               |

---

## Header

```tsx
<Header title="Page Title" onBack={() => navigate(-1)} position="sticky">
  <Button variant="text">Action</Button>
</Header>
```

| Prop             | Type                                 | Default | Description        |
| ---------------- | ------------------------------------ | ------- | ------------------ |
| `title`          | `ReactNode \| string`                | —       | Header title       |
| `onBack`         | `() => void`                         | —       | Shows back arrow   |
| `position`       | `'sticky' \| 'relative'`             | —       | CSS position       |
| `background`     | `'none' \| 'solid' \| 'transparent'` | —       | Background style   |
| `layout`         | `'center' \| 'stretch'`              | —       | Content layout     |
| `children`       | `ReactNode`                          | —       | Right-side content |
| `secondaryChild` | `ReactNode`                          | —       | Secondary row      |

---

## Icon

```tsx
<Icon name="chevron-right" size="md" />
<Icon name="check-circle" variant="success" />
<Icon name="close" size="sm" color="#FF0000" />
```

| Prop       | Type                                 | Default      | Description      |
| ---------- | ------------------------------------ | ------------ | ---------------- |
| `name`     | `string`                             | **required** | Icon identifier  |
| `size`     | `'sm' \| 'md' \| 'lg' \| 'xl'`       | —            | Icon size        |
| `variant`  | `'primary' \| 'success' \| 'danger'` | —            | Predefined color |
| `color`    | `string`                             | —            | Custom color     |
| `disabled` | `boolean`                            | —            | Dimmed state     |

---

## Info

```tsx
<Info type="success" title="Success" description="Operation completed." />
<Info type="error" description="Something went wrong." closeable />
<Info type="warning" variant="antique" description="Deprecated feature." />
```

| Prop          | Type                                          | Default      | Description       |
| ------------- | --------------------------------------------- | ------------ | ----------------- |
| `type`        | `'info' \| 'success' \| 'warning' \| 'error'` | **required** | Alert type        |
| `variant`     | `'simple' \| 'antique'`                       | `'simple'`   | Visual variant    |
| `title`       | `ReactNode`                                   | —            | Alert title       |
| `description` | `ReactNode`                                   | **required** | Alert message     |
| `closeable`   | `boolean`                                     | —            | Show close button |

---

## Label

```tsx
<Label type="success" label="Active" />
<Label type="danger" label="Failed" border />
```

| Prop     | Type                                                        | Default      | Description   |
| -------- | ----------------------------------------------------------- | ------------ | ------------- |
| `type`   | `'success' \| 'danger' \| 'info' \| 'warning' \| 'default'` | **required** | Color variant |
| `label`  | `string`                                                    | **required** | Display text  |
| `border` | `boolean`                                                   | —            | Show border   |

---

## Chip

```tsx
<Chip label="Category" size="sm" isActive onClick={() => {}} />
<Chip label="Removable" suffix="close" onSuffixClick={handleRemove} />
```

| Prop            | Type                                     | Default      | Description           |
| --------------- | ---------------------------------------- | ------------ | --------------------- |
| `label`         | `string`                                 | **required** | Chip text             |
| `size`          | `'sm' \| 'lg'`                           | —            | Chip size             |
| `isActive`      | `boolean`                                | —            | Active/selected state |
| `isDisable`     | `boolean`                                | —            | Disabled state        |
| `suffix`        | `'close' \| 'chevron-down' \| ReactNode` | —            | Trailing content      |
| `prefix`        | `ReactNode`                              | —            | Leading content       |
| `onClick`       | `() => void`                             | —            | Click handler         |
| `onSuffixClick` | `() => void`                             | —            | Suffix click handler  |

---

## Popover

```tsx
<Popover content={<div>Popover body</div>} title="Info" trigger="click" placement="bottom">
  <Button variant="text">Open</Button>
</Popover>
```

| Prop           | Type                                                                                                                                                             | Default | Description           |
| -------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | --------------------- |
| `content`      | `ReactNode`                                                                                                                                                      | —       | Popover body          |
| `title`        | `ReactNode`                                                                                                                                                      | —       | Popover header        |
| `trigger`      | `'hover' \| 'click' \| 'focus'`                                                                                                                                  | —       | Trigger type          |
| `placement`    | `'top' \| 'bottom' \| 'left' \| 'right' \| 'topLeft' \| 'topRight' \| 'bottomLeft' \| 'bottomRight' \| 'leftTop' \| 'leftBottom' \| 'rightTop' \| 'rightBottom'` | —       | Position              |
| `arrow`        | `boolean`                                                                                                                                                        | —       | Show arrow            |
| `open`         | `boolean`                                                                                                                                                        | —       | Controlled open state |
| `offset`       | `number`                                                                                                                                                         | `8`     | Distance from trigger |
| `portal`       | `boolean`                                                                                                                                                        | —       | Render in portal      |
| `onOpenChange` | `(isOpen: boolean) => void`                                                                                                                                      | —       | Open state handler    |
| `disabled`     | `boolean`                                                                                                                                                        | —       | Disabled state        |

---

## Steps

```tsx
<Steps
  variant="progress"
  currentStep={1}
  orientation="horizontal"
  steps={[{ title: 'Account' }, { title: 'Profile' }, { title: 'Review' }]}
/>
```

| Prop             | Type                                  | Default      | Description                                        |
| ---------------- | ------------------------------------- | ------------ | -------------------------------------------------- |
| `steps`          | `Step[]`                              | **required** | `{ title?, content?, icon?, deactive?, checked? }` |
| `currentStep`    | `number`                              | **required** | Active step index                                  |
| `variant`        | `'progress' \| 'line' \| 'indicator'` | —            | Visual style                                       |
| `orientation`    | `'vertical' \| 'horizontal'`          | —            | Layout direction                                   |
| `labelPlacement` | `'vertical' \| 'horizontal'`          | —            | Label position                                     |
| `onStepClick`    | `(index: number) => void`             | —            | Step click handler                                 |

---

## Listing

```tsx
<Listing
  listItems={[
    { id: '1', content: 'Item One' },
    { id: '2', content: 'Item Two' }
  ]}
  selectionType="checkbox"
  selectedIds={['1']}
  onSelected={(ids) => {}}
/>
```

| Prop            | Type                                           | Default      | Description                          |
| --------------- | ---------------------------------------------- | ------------ | ------------------------------------ |
| `listItems`     | `ListItem[]`                                   | **required** | `{ id: string, content: ReactNode }` |
| `selectionType` | `'checkbox' \| 'radio' \| 'toggle' \| 'arrow'` | —            | Selection UI                         |
| `selectedIds`   | `string[]`                                     | —            | Selected item IDs                    |
| `onSelected`    | `(ids: string[]) => void`                      | —            | Selection change                     |
| `onClickItem`   | `(id, content) => void`                        | —            | Item click                           |

---

## Collapse

```tsx
<Collapse header="Section Title" expandIconPosition="end">
  Collapsible content here
</Collapse>
```

| Prop                 | Type                        | Default      | Description         |
| -------------------- | --------------------------- | ------------ | ------------------- |
| `header`             | `ReactNode`                 | **required** | Collapse trigger    |
| `children`           | `ReactNode`                 | **required** | Collapsible content |
| `isOpen`             | `boolean`                   | —            | Controlled state    |
| `expandIconPosition` | `'start' \| 'end'`          | —            | Chevron position    |
| `onChange`           | `(isOpen: boolean) => void` | —            | Toggle handler      |

---

## Divider

```tsx
<Divider />
<Divider type="dashed" />
```

| Prop   | Type                  | Default   | Description |
| ------ | --------------------- | --------- | ----------- |
| `type` | `'solid' \| 'dashed'` | `'solid'` | Line style  |

---

## Spinner

```tsx
<Spinner size="md" />
<Spinner size="lg" inverted />
```

| Prop       | Type                   | Default | Description          |
| ---------- | ---------------------- | ------- | -------------------- |
| `size`     | `'sm' \| 'md' \| 'lg'` | —       | Spinner size         |
| `inverted` | `boolean`              | —       | For dark backgrounds |

---

## Skeleton

```tsx
<Skeleton size="md" />
```

| Prop       | Type           | Default | Description          |
| ---------- | -------------- | ------- | -------------------- |
| `size`     | `'sm' \| 'md'` | —       | Skeleton height      |
| `inverted` | `boolean`      | —       | For dark backgrounds |

---

## Title

```tsx
<Title level={1}>Page Heading</Title>
<Title level={3} italic>Subheading</Title>
```

| Prop        | Type                         | Default      | Description           |
| ----------- | ---------------------------- | ------------ | --------------------- |
| `level`     | `1 \| 2 \| 3 \| 4 \| 5 \| 6` | **required** | Heading level (h1–h6) |
| `italic`    | `boolean`                    | —            | Italic style          |
| `underline` | `boolean`                    | —            | Underline style       |

---

## Paragraph

```tsx
<Paragraph size="r" weight="regular">Body text</Paragraph>
<Paragraph size="l" weight="bold">Emphasized text</Paragraph>
```

| Prop        | Type                                | Default | Description     |
| ----------- | ----------------------------------- | ------- | --------------- |
| `size`      | `'s' \| 'r' \| 'l' \| 'xl'`         | —       | Text size       |
| `weight`    | `'regular' \| 'semibold' \| 'bold'` | —       | Font weight     |
| `italic`    | `boolean`                           | —       | Italic style    |
| `underline` | `boolean`                           | —       | Underline style |

---

## Toaster

Requires `ToasterProvider` wrapping the app.

```tsx
import { useToaster } from '@squantumengine/horizon';

const { notify } = useToaster();
notify({ type: 'success', message: 'Saved successfully!' });
```

---

## ImageZoom

```tsx
<ImageZoom src="/image.png" alt="Zoomable image" />
```

| Prop  | Type     | Default      | Description      |
| ----- | -------- | ------------ | ---------------- |
| `src` | `string` | **required** | Image source URL |
| `alt` | `string` | —            | Alt text         |

---

## Design Tokens

```tsx
import { Token } from '@squantumengine/horizon';

Token.COLORS; // Color palette
Token.SIZES; // sm, md, lg
Token.TYPOGRAPHY; // Font definitions
Token.THEMES; // Theme definitions
Token.BUTTON_VARIANTS; // primary, secondary, text
Token.LABEL_TYPE; // success, danger, info, warning, default
Token.INFO_TYPE; // info, success, warning, error
Token.TAB_POSITION; // top, bottom, left, right
```

## Storybook

Full interactive documentation: https://sqehorizon.squantumengine.com
