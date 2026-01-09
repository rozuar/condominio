/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  eslint: {
    // Avoid CI/build failures due to eslint tooling mismatch.
    // Lint can still be run explicitly via `npm run lint`.
    ignoreDuringBuilds: true,
  },
}

module.exports = nextConfig
