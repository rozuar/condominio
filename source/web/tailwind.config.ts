import type { Config } from 'tailwindcss'

const config: Config = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          // Sophisticated brand palette (muted, professional)
          DEFAULT: '#0F3D2E', // deep forest
          light: '#145A43',
          dark: '#0B2A1F',
        },
        tierra: {
          DEFAULT: '#B7791F',
          light: '#D69E2E',
          dark: '#975A16',
        },
        agua: {
          DEFAULT: '#2563EB',
          light: '#60A5FA',
          dark: '#1D4ED8',
        },
        cielo: '#0EA5E9',
        accent: '#7C3AED',
        success: '#16A34A',
        warning: '#CA8A04',
      },
    },
  },
  plugins: [],
}
export default config
