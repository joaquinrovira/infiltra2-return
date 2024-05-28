/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.{html,js,templ}"],
  theme: {
    extend: {
      colors: {
        accent: "rgb(var(--color-accent) / <alpha-value>)",
        primary: "rgb(var(--color-primary) / <alpha-value>)",
        primaryDark: "rgb(var(--color-primary-dark) / <alpha-value>)",
        secondary: "rgb(var(--color-secondary) / <alpha-value>)",
        secondaryDark: "rgb(var(--color-secondary-dark) / <alpha-value>)",
      },
      fontFamily: {
        ["cthulu"]: '"Cthulhumbus", system-ui', // Adds a new `font-display` class
      },
    },
  },
  plugins: [],
};

