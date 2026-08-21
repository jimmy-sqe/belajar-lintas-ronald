import { Token } from '@squantumengine/horizon';
import type { Config } from 'tailwindcss';

// @todo use Token.COLOR_ROLES.COLOR_ROLES_TAILWIND from horizon, to override primary colors
const COLOR_ROLES = {
  COLOR_ROLES_TAILWIND: {
    primary: {
      DEFAULT: 'var(--hz-horizon-primary)',
      dark: 'var(--hz-horizon-primary-dark)',
      medium: 'var(--hz-horizon-primary-medium)',
      light: 'var(--hz-horizon-primary-light)',
      soft: 'var(--hz-horizon-primary-soft)',
      pale: 'var(--hz-horizon-primary-pale)'
    },
    success: {
      DEFAULT: 'var(--hz-horizon-success)',
      dark: 'var(--hz-horizon-success-dark)',
      medium: 'var(--hz-horizon-success-medium)',
      light: 'var(--hz-horizon-success-light)',
      soft: 'var(--hz-horizon-success-soft)',
      pale: 'var(--hz-horizon-success-pale)'
    },
    warning: {
      DEFAULT: 'var(--hz-horizon-warning)',
      dark: 'var(--hz-horizon-warning-dark)',
      medium: 'var(--hz-horizon-warning-medium)',
      light: 'var(--hz-horizon-warning-light)',
      soft: 'var(--hz-horizon-warning-soft)',
      pale: 'var(--hz-horizon-warning-pale)'
    },
    error: {
      DEFAULT: 'var(--hz-horizon-error)',
      dark: 'var(--hz-horizon-error-dark)',
      medium: 'var(--hz-horizon-error-medium)',
      light: 'var(--hz-horizon-error-light)',
      soft: 'var(--hz-horizon-error-soft)',
      pale: 'var(--hz-horizon-error-pale)'
    },
    info: {
      DEFAULT: 'var(--hz-horizon-info)',
      dark: 'var(--hz-horizon-info-dark)',
      medium: 'var(--hz-horizon-info-medium)',
      light: 'var(--hz-horizon-info-light)',
      soft: 'var(--hz-horizon-info-soft)',
      pale: 'var(--hz-horizon-info-pale)'
    }
  }
};

const { COLORS } = Token;

const config: Config = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}'
  ],
  theme: {
    extend: {
      colors: {
        ...COLOR_ROLES.COLOR_ROLES_TAILWIND,
        hz: {
          ...COLORS
        }
      }
    }
  },
  plugins: []
};
export default config;
