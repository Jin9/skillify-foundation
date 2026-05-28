import type { Config } from 'tailwindcss';

const config: Config = {
  content: [
    './pages/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          50:  '#e6f4f0',
          100: '#c2e2d8',
          300: '#6dbfa3',
          400: '#3da382',
          500: '#1a8060',
          600: '#1a8060',
          700: '#157053',
          900: '#0d4a37',
        },
      },
      fontFamily: {
        sans: ['IBM Plex Sans Thai', 'Sarabun', 'sans-serif'],
      },
    },
  },
  plugins: [],
};

export default config;
