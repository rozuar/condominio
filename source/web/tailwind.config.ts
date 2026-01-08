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
          DEFAULT: '#00C853',
          light: '#69F0AE',
          dark: '#00A844',
        },
        tierra: {
          DEFAULT: '#FF9800',
          light: '#FFB74D',
          dark: '#F57C00',
        },
        agua: {
          DEFAULT: '#00B0FF',
          light: '#40C4FF',
          dark: '#0091EA',
        },
        cielo: '#00E5FF',
        accent: '#FF4081',
        success: '#00E676',
        warning: '#FFEA00',
      },
    },
  },
  plugins: [],
}
export default config
