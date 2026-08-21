export function restoreFormData<T = Record<string, unknown>>(formData: FormData): T {
  const result: Record<string, unknown> = {};

  for (const [rawKey, value] of formData.entries()) {
    const keys = rawKey.replace(/]/g, '').split('[');
    let current: Record<string, unknown> = result;

    for (let i = 0; i < keys.length; i++) {
      const k = keys[i];
      const isLast = i === keys.length - 1;

      if (isLast) {
        let parsed: unknown;

        if (typeof value === 'string') {
          try {
            parsed = JSON.parse(value);
          } catch {
            parsed = value;
          }
        } else {
          parsed = value;
        }
        current[k] = parsed;
      } else {
        const nextKey = keys[i + 1];
        const isNextArray = /^\d+$/.test(nextKey);

        if (current[k] === undefined) {
          current[k] = isNextArray ? [] : {};
        }
        current = current[k] as Record<string, unknown>;
      }
    }
  }

  return result as T;
}

export function buildFormData<T = unknown>(data: T): FormData {
  const formData = new FormData();

  function appendData(prefix: string, value: unknown) {
    if (value === undefined || value === null) return;

    if (value instanceof File || value instanceof Blob) {
      formData.append(prefix, value);

      return;
    }
    if (value instanceof Date) {
      formData.append(prefix, value.toISOString());

      return;
    }
    if (Array.isArray(value)) {
      value.forEach((item, index) => {
        appendData(`${prefix}[${index}]`, item);
      });

      return;
    }
    if (typeof value === 'object') {
      for (const [key, val] of Object.entries(value)) {
        appendData(`${prefix}[${key}]`, val);
      }

      return;
    }
    formData.append(prefix, String(value));
  }

  if (data) {
    Object.entries(data).forEach(([key, value]) => {
      if (value === undefined || value === null) return;
      appendData(key, value);
    });
  }

  return formData;
}
