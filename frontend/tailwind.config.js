/** @type {import('tailwindcss').Config} */
export default {
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
        },
        accent: '#f97316',
      },
      fontFamily: {
        display: ['"DM Sans"', 'Inter', 'system-ui', 'sans-serif'],
        body: ['Inter', 'system-ui', 'sans-serif'],
      },
      boxShadow: {
        soft: '0 10px 25px rgba(15, 23, 42, 0.1)',
      },
    },
  },
  plugins: [require('@tailwindcss/forms')],
};
