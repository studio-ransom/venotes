/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        discord: {
          dark: '#2C2F33',
          darker: '#23272A',
          darkest: '#1E2124',
          light: '#99AAB5',
          lighter: '#DCDDDE',
          blue: '#5865F2',
          green: '#57F287'
        }
      }
    },
  },
  plugins: [],
}