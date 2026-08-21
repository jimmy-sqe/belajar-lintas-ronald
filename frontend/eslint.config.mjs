import { FlatCompat } from '@eslint/eslintrc';
import js from '@eslint/js';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import tsParser from '@typescript-eslint/parser';
import typescriptEslint from '@typescript-eslint/eslint-plugin';
import react from 'eslint-plugin-react';
import importPlugin from 'eslint-plugin-import';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const compat = new FlatCompat({
  baseDirectory: __dirname,
  recommendedConfig: js.configs.recommended,
  allConfig: js.configs.all
});

import nextPlugin from '@next/eslint-plugin-next';

export default [
  {
    ignores: [
      '.next/**',
      'next-env.d.ts',
      'public/mockServiceWorker.js',
      'next.config.mjs',
      'postcss.config.mjs',
      'tailwind.config.ts',
      'tailwind.config.js',
      '**/*.config.*',
      'src/openapi/.kubb/**'
    ]
  },
  js.configs.recommended,
  ...compat.extends('next/core-web-vitals'),
  ...compat.extends('plugin:react/recommended'),
  ...compat.extends('plugin:@typescript-eslint/recommended'),
  {
    files: ['**/*.ts', '**/*.tsx'],

    plugins: {
      react,
      '@typescript-eslint': typescriptEslint,
      '@next/next': nextPlugin,
      import: importPlugin
    },

    languageOptions: {
      parser: tsParser,
      ecmaVersion: 5,
      sourceType: 'script',

      parserOptions: {
        project: './tsconfig.json'
      }
    },

    rules: {
      '@typescript-eslint/no-unused-vars': [
        'error',
        { varsIgnorePattern: '^_', argsIgnorePattern: '^_', caughtErrorsIgnorePattern: '^_' }
      ],
      '@typescript-eslint/no-var-requires': 'off',
      'react/react-in-jsx-scope': 'off',
      '@typescript-eslint/consistent-type-exports': 'error',
      '@typescript-eslint/consistent-type-imports': 'error',
      'no-mixed-spaces-and-tabs': 0,

      'import/order': [
        'error',
        {
          groups: ['external', 'builtin', 'internal', 'sibling', 'parent', 'index'],

          pathGroups: [
            {
              pattern: 'react',
              group: 'external',
              position: 'before'
            },
            {
              pattern: 'react*',
              group: 'external',
              position: 'before'
            },
            {
              pattern: 'react*/**',
              group: 'external',
              position: 'before'
            },
            {
              pattern: '@/features/**',
              group: 'internal',
              position: 'before'
            },
            {
              pattern: '@/common/stores/**',
              group: 'internal',
              position: 'before'
            },
            {
              pattern: '@/common/components/**',
              group: 'internal',
              position: 'before'
            },
            {
              pattern: '@/common/utils/**',
              group: 'internal',
              position: 'before'
            },
            {
              pattern: '@/common/hooks/**',
              group: 'internal',
              position: 'before'
            },
            {
              pattern: '@/common/types/**',
              group: 'internal',
              position: 'before'
            },
            {
              pattern: '@/common/constants/**',
              group: 'internal',
              position: 'before'
            },
            {
              pattern: '@/services/**',
              group: 'internal',
              position: 'before'
            }
          ],

          pathGroupsExcludedImportTypes: ['react', 'react-router'],

          alphabetize: {
            order: 'asc',
            caseInsensitive: true
          }
        }
      ]
    }
  },
  {
    files: ['src/openapi/**/*.ts', 'src/openapi/**/*.tsx'],
    rules: {
      '@typescript-eslint/no-explicit-any': 'off'
    }
  }
];
