module.exports = {
  content: ['./src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      fontFamily: {
        display: ['IBM Plex Serif'],
        body: ['IBM Plex Sans'],
      },
      borderRadius: {
        sm: '10px',
        DEFAULT: '20px',
        md: '30px',
      },
    },
  },
  plugins: [],
}
