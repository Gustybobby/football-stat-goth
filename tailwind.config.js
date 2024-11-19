/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.html", "./**/*.templ", "./**/*.go"],
  theme: {
    extend: {
      fontFamily: {
        poppins: ["Poppins"],
      },
      colors: {
        primary: {
          DEFAULT: "#d1d5db",
          background: "white",
          foreground: "#f3f4f6",
        },
        secondary: {
          DEFAULT: "#3730a3",
          background: "#4f46e5",
          foreground: "#6366f1",
        },
      },
      backgroundImage: {
        'footballField': "url('/public/football-field-design-vector.jpg')",
        'player_card': "url('/public/playerCard.png')"
      },
    },
  },
  plugins: [require('tailwindcss-motion')], 
};
