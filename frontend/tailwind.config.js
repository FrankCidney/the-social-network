/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        canvas: '#F9FAFB',
        surface: '#FFFFFF',
        primary: '#4F46E5',
        'text-main': '#111827',
        privacy: {
          public: '#10B981',
          private: '#EF4444',
        }
      },
      borderRadius: {
        'bento': '12px',
      }
    },
  },
  plugins: [],
}
