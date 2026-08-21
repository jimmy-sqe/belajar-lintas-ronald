import { z } from 'zod';

export const customZodErrorMap: z.ZodErrorMap = (issue, ctx) => {
  if (issue.code === z.ZodIssueCode.invalid_type) {
    return { message: 'Wajib diisi' };
  }
  if (issue.code === z.ZodIssueCode.invalid_string) {
    if (issue.validation === 'email') return { message: 'Format email salah' };
    return { message: ctx.defaultError };
  }
  if (issue.code === z.ZodIssueCode.too_small) {
    if (issue.type === 'string') return { message: `Panjang karakter minimal ${issue.minimum}` };
    return { message: ctx.defaultError };
  }
  if (issue.code === z.ZodIssueCode.too_big) {
    if (issue.type === 'string') return { message: `Panjang karakter maksimal ${issue.maximum}` };
    if (issue.type === 'array' && issue.path[0] === 'file_attachment_ids')
      return { message: `Maksimal ${issue.maximum} file yang dapat diupload` };
    return { message: ctx.defaultError };
  }

  return { message: ctx.defaultError };
};
