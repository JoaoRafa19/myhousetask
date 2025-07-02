// tailwind.config.js
/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
      "./internal/web/view/**/*.templ", // Diga para ele ler TODOS os arquivos .templ
    ],
    theme: {
      extend: {},
    },
    plugins: [],
  }