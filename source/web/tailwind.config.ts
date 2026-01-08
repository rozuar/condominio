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
          DEFAULT: '#2D5016',
          light: '#4A7C23',
          dark: '#1F3810',
        },
        tierra: '#8B7355',
        agua: '#3B82A0',
      },
    },
  },
  plugins: [],
}
export default config
