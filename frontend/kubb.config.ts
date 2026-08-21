import { defineConfig } from '@kubb/core';
import { pluginOas } from '@kubb/plugin-oas';
import { pluginTs } from '@kubb/plugin-ts';
import { pluginZod } from '@kubb/plugin-zod';
import { pluginReactQuery } from '@kubb/plugin-react-query';
import { URLPath } from '@kubb/core/utils';

export default defineConfig(() => {
  return {
    root: './src',
    input: {
      path: './openapi/openapi.json'
    },
    output: {
      path: './openapi',
      clean: true
    },
    plugins: [
      pluginOas({}),
      pluginTs({
        output: { path: 'types' },
        unknownType: 'unknown',
        group: {
          type: 'tag',
          name: ({ group }) => group?.toLowerCase().replace(/\s+/g, '-').replace(/-+/g, '-')
        }
      }),
      pluginZod({
        output: { path: 'zod' },
        typed: false,
        group: {
          type: 'tag',
          name: ({ group }) => group?.toLowerCase().replace(/\s+/g, '-').replace(/-+/g, '-')
        }
      }),
      pluginReactQuery({
        output: {
          path: './services'
        },
        group: {
          type: 'tag',
          name: ({ group }) => group?.toLowerCase().replace(/\s+/g, '-').replace(/-+/g, '-')
        },
        client: {
          dataReturnType: 'data',
          importPath: '@/services/fetcher'
        },
        mutation: {
          methods: ['post', 'put', 'patch', 'delete']
        },
        query: {
          methods: ['get']
        },
        queryKey: ({ operation, schemas, casing }) => {
          const path = new URLPath(operation.path, { casing });
          const keys = [
            JSON.stringify(path.toURLPath()),
            schemas.queryParams?.name ? '...(params ? [params] : [])' : undefined,
            schemas.request?.name ? '...(data ? [data] : [])' : undefined,
            schemas.pathParams?.keys
              ? `...(${schemas.pathParams.keys} ? [${schemas.pathParams.keys}] : [])`
              : undefined
          ].filter(Boolean);

          return keys;
        }
      })
    ]
  };
});
