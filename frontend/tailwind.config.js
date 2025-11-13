/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'class',
  content: ['./index.html', './src/**/*.{js,jsx,ts,tsx}'],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eef6ff',
          100: '#d9e9ff',
          200: '#b7d4ff',
          300: '#84b7ff',
          400: '#5290ff',
          500: '#3367f0',
          600: '#2349c7',
          700: '#1f3ba0',
          800: '#1d357f',
          900: '#1d315f',
          950: '#0f1a3d',
        },
        accent: {
          light: '#fbbf24',
          DEFAULT: '#f59e0b',
          dark: '#d97706',
        },
        luxury: {
          gold: '#d4af37',
          platinum: '#e5e4e2',
        },
      },
      fontFamily: {
        display: ['"Playfair Display"', '"Cormorant Garamond"', 'serif'],
        body: ['Inter', '"Segoe UI"', 'system-ui', 'sans-serif'],
        elegant: ['"Crimson Text"', 'Georgia', 'serif'],
      },
      boxShadow: {
        soft: '0 10px 25px rgba(15, 23, 42, 0.08)',
        elegant: '0 4px 20px rgba(15, 23, 42, 0.12)',
        luxury: '0 8px 30px rgba(31, 59, 160, 0.15)',
      },
      backgroundImage: {
        'gradient-luxury': 'linear-gradient(135deg, #fafbfc 0%, #f1f5f9 100%)',
        'gradient-luxury-dark': 'linear-gradient(135deg, #0f172a 0%, #1e293b 100%)',
      },
    },
  },
  plugins: [require('@tailwindcss/forms')],
};
