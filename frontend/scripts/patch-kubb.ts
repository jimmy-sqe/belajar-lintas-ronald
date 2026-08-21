/**
 * Post-kubb-generate patch script.
 *
 * 1. Replaces the generated .kubb/config.ts with a re-export from src/utils/form-data.ts
 *    so that the custom buildFormData / restoreFormData implementations are preserved.
 * 2. Fixes `formData as FormData` type errors in generated services by changing the
 *    cast to `formData as unknown as <RequestType>`.
 *
 * Usage: bun scripts/patch-kubb.ts
 */

import { writeFileSync, readFileSync, readdirSync, statSync } from 'node:fs';
import { resolve, join } from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = resolve(__filename, '..');

// Patch 1: Replace .kubb/config.ts
const target = resolve(__dirname, '../src/openapi/.kubb/config.ts');

const content = `export { buildFormData, restoreFormData } from "@/utils/form-data";
`;

writeFileSync(target, content, 'utf-8');
console.log(`✔ Patched ${target}`);

// Patch 2: Fix `formData as FormData` type assertion in generated services
const servicesDir = resolve(__dirname, '../src/openapi/services');

function walkDir(dir: string): string[] {
  const files: string[] = [];
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry);
    if (statSync(full).isDirectory()) {
      files.push(...walkDir(full));
    } else if (full.endsWith('.ts')) {
      files.push(full);
    }
  }
  return files;
}

let patchCount = 0;
for (const file of walkDir(servicesDir)) {
  const src = readFileSync(file, 'utf-8');
  // Match: data: formData as FormData
  // Replace with: data: formData as unknown as <RequestType>
  // The request type is the TVariables generic in the preceding request<...> call
  const patched = src.replace(
    /data: formData as FormData/g,
    'data: formData as unknown as typeof requestData'
  );
  if (patched !== src) {
    writeFileSync(file, patched, 'utf-8');
    patchCount++;
    console.log(`✔ Patched FormData cast in ${file}`);
  }
}
console.log(`✔ Patched ${patchCount} file(s) with FormData fix`);

// Patch 3: Fix recursive Zod schemas that reference themselves via z.lazy()
// TypeScript can't infer the type of a const that references itself.
const zodDir = resolve(__dirname, '../src/openapi/zod');

let zodPatchCount = 0;
for (const file of walkDir(zodDir)) {
  const src = readFileSync(file, 'utf-8');
  // Find schemas that use z.lazy(() => <sameSchemaName>)
  const patched = src.replace(
    /export const (\w+)(:\s*z\.ZodType)?\s*=\s*(z\.object\(\{[\s\S]*?z\.lazy\(\(\) => \1\)[\s\S]*?\}\))/g,
    (match, name, existingType) => {
      if (existingType) return match; // Already typed
      return match.replace(`export const ${name} =`, `export const ${name}: z.ZodType =`);
    }
  );
  if (patched !== src) {
    writeFileSync(file, patched, 'utf-8');
    zodPatchCount++;
    console.log(`✔ Patched recursive Zod schema in ${file}`);
  }
}
console.log(`✔ Patched ${zodPatchCount} file(s) with recursive Zod fix`);
