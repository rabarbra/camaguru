import daisyui from "daisyui"
import typography from "@tailwindcss/typography"
import defaultTheme from "tailwindcss/defaultTheme"

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        'sans': ['"Kode Mono"', ...defaultTheme.fontFamily.sans],
        // 'sans': ['"Syne Mono"', ...defaultTheme.fontFamily.sans],
        // 'sans': ['"M PLUS 1 Code"', ...defaultTheme.fontFamily.sans],
        // 'sans': ['"Cutive Mono"', ...defaultTheme.fontFamily.sans],
        // 'sans': ['"Orbitron"', ...defaultTheme.fontFamily.sans],
        // 'sans': ['"Gruppo"', ...defaultTheme.fontFamily.sans],
      },
    },
  },
  plugins: [
    typography,
    daisyui,
  ],
  daisyui: {
    themes: [
        "nord",
        {
          nordDark: {
            "color-scheme": "dark",
            "primary": "#5E81AC",
            "primary-content": "#03060b",
            "secondary": "#8FBCBB",
            "secondary-content": "#070d0d",
            "accent": "#81A1C1",
            "accent-content": "#06090e",
            "neutral": "#3B4252",
            "neutral-content": "#ECEFF4",
            "base-100": "#2E3440",
            "base-200": "#272c36",
            "base-300": "#1f242d",
            "base-content": "#d8dee9",
            "info": "#88C0D0",
            "info-content": "#070e10",
            "success": "#A3BE8C",
            "success-content": "#0a0d07",
            "warning": "#EBCB8B",
            "warning-content": "#130f07",
            "error": "#BF616A",
            "error-content": "#0d0304",
            "--rounded-box": "0.4rem",
            "--rounded-btn": "0.2rem",
            "--rounded-badge": "0.4rem",
            "--tab-radius": "0.2rem",
          },
      },
    ],
    darkTheme: "nordDark",
    base: true,
    styled: true,
    utils: true,
    prefix: "",
    logs: true,
    themeRoot: ":root",
  },
}
