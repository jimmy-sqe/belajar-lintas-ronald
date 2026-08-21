import type { IPasswordRule } from '@/static/types/auth';

export const ruleRegex = {
  capital: /[A-Z]/,
  number: /[0-9]/,
  lower: /[a-z]/,
  symbol: /[^A-Za-z0-9]/
};

export const validateRules = (passwordRules: IPasswordRule[], passwordWatch: string) => {
  return passwordRules.map((rule) => {
    let passed = false;
    switch (rule.text) {
      case 'Huruf kapital':
        passed = ruleRegex.capital.test(passwordWatch);
        break;
      case 'Min. 8 karakter':
        passed = passwordWatch.length >= 8;
        break;
      case 'Angka':
        passed = ruleRegex.number.test(passwordWatch);
        break;
      case 'Huruf kecil':
        passed = ruleRegex.lower.test(passwordWatch);
        break;
      case 'Simbol':
        passed = ruleRegex.symbol.test(passwordWatch);
        break;
    }
    return { ...rule, passed };
  });
};
// Lazy-loaded password rules to avoid server-side Token.COLORS access
export const getInitialPassRules = () => {
  // Import Token only when called (client-side)
  const { Token } = require('@squantumengine/horizon');

  return [
    {
      text: 'Huruf kapital',
      icon: 'check-circle',
      color: Token.COLORS.neutral[600],
      passed: false
    },
    {
      text: 'Min. 8 karakter',
      icon: 'check-circle',
      color: Token.COLORS.neutral[600],
      passed: false
    },
    {
      text: 'Angka',
      icon: 'check-circle',
      color: Token.COLORS.neutral[600],
      passed: false
    },
    {
      text: 'Huruf kecil',
      icon: 'check-circle',
      color: Token.COLORS.neutral[600],
      passed: false
    },
    {
      text: 'Simbol',
      icon: 'check-circle',
      color: Token.COLORS.neutral[600],
      passed: false
    }
  ];
};
