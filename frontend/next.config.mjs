/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  eslint: {
    ignoreDuringBuilds: true,
  },
  //   logging: {
  //     fetches: {
  //       fullUrl: true,
  //     },
  //     // This forwards client console logs to your terminal
  //     browserToTerminal: true,
  //   },
  turbopack: {
    // Custom aliases for module resolution
    resolveAlias: {
      underscore: 'lodash'
    },
    // Custom file extensions to resolve
    resolveExtensions: ['.mdx', '.tsx', '.ts', '.jsx', '.js', '.json'],
    // Custom rules for file transformations (loaders)
    rules: {
      '*.svg': {
        loaders: ['@svgr/webpack'],
        as: '*.js'
      }
    }
  }
};

export default nextConfig;
