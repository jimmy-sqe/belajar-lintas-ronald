import { appName } from '@/config/environment';

export const prefixKey = appName.replace(/\s+/g, '_');

export const sessionKeyEnum = {
  SESSION: `${prefixKey}_s`,
  SELECTED_COMPANY: `${prefixKey}_s_c`
} as const;
