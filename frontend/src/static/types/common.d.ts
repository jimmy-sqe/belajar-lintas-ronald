export type Primitive = string | number | boolean | bigint | symbol | null | undefined | Date;

// Helper for array indices
type ArrayIndex = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9;

// Helper to get keys at depth 0 (immediate properties)
type Level0Keys<T> = keyof T & string;

// Helper to get keys at depth 1 (one level deep)
type Level1Keys<T> = {
  [K in keyof T]: T[K] extends Array<infer U>
    ? U extends object
      ? `${K & string}.${ArrayIndex}`
      : never
    : T[K] extends object | undefined
      ? NonNullable<T[K]> extends object
        ? `${K & string}.${keyof NonNullable<T[K]> & string}`
        : never
      : never;
}[keyof T];

// Helper to get keys at depth 2 (two levels deep)
type Level2Keys<T> = {
  [K in keyof T]: T[K] extends Array<infer U>
    ? U extends object
      ? `${K & string}.${ArrayIndex}.${keyof U & string}`
      : never
    : T[K] extends object | undefined
      ? NonNullable<T[K]> extends object
        ? {
            [K2 in keyof NonNullable<T[K]>]: NonNullable<T[K]>[K2] extends Array<infer U2>
              ? U2 extends object
                ? `${K & string}.${K2 & string}.${ArrayIndex}`
                : never
              : NonNullable<T[K]>[K2] extends object | undefined
                ? NonNullable<NonNullable<T[K]>[K2]> extends object
                  ? `${K & string}.${K2 & string}.${keyof NonNullable<NonNullable<T[K]>[K2]> & string}`
                  : never
                : never;
          }[keyof NonNullable<T[K]>]
        : never
      : never;
}[keyof T];

export type DeepKeyOf<T> = Level0Keys<T> | Level1Keys<T> | Level2Keys<T>;

export type PathValue<T, P extends string> = P extends `${infer K}.${infer Rest}`
  ? K extends keyof T
    ? PathValue<T[K], Rest>
    : never
  : P extends keyof T
    ? T[P]
    : never;
