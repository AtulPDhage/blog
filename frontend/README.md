# Postly Frontend (Quasar SPA)

Postly is a production-grade, highly responsive, and feature-rich blogging application frontend. It has been migrated from Next.js to **Vue 3 (Composition API)**, **Quasar Framework**, and **Pinia Store** for unified global state management.

---

## Tech Stack & Architecture

- **Core Framework**: Vue 3 (Composition API using `<script setup lang="ts">`)
- **UI Component Library**: Quasar Framework (v2)
- **Build Engine**: Vite 3 (`@quasar/app-vite`)
- **State Management**: Pinia Store
- **HTTP Client**: Axios
- **Security & Sanitization**: DOMPurify (for Jodit Rich Text XSS prevention)
- **Auth Protocol**: Google OAuth 2.0 (Authorization Code Grant via Google Identity Services)

---

## Folder Layout

The project folder is organized according to standard Quasar recommendations:

```
frontend/
  ├── public/                   # Static assets (images, logos)
  ├── src/
  │   ├── boot/                 # Startup hooks (axios, dompurify configs)
  │   ├── components/           # Reusable Vue components (BlogCard, LoadingSpinner)
  │   ├── layouts/              # App shells & frames (MainLayout.vue)
  │   ├── pages/                # Route views (IndexPage, LoginPage, ProfilePage, etc.)
  │   ├── router/               # Vue Router configuration & Navigation guards
  │   ├── stores/               # Pinia store state management
  │   ├── App.vue               # Root application component
  │   └── env.d.ts              # TS global environment variables & wrappers modules
  ├── tsconfig.json             # Root TypeScript config
  ├── quasar.config.ts          # Quasar framework configurations
  └── package.json              # Project script runner and dependency manager
```

---

## Features

1. **Google OAuth popup authentication**: Login securely using your Google account.
2. **Dynamic Blog Feed & Filters**: Search blogs or filter them by category via a responsive left drawer.
3. **Save Post (Bookmarks)**: Save your favorite blog posts to read later (stored per user).
4. **Rich Text Content Editor**: Create and edit posts using Jodit editor with full formatting options.
5. **AI Assistants**: Use built-in AI helpers to improve blog titles, optimize descriptions, and automatically fix grammar inside the Jodit editor.
6. **Comments Section**: Authenticated users can leave or delete comments on blog posts.

---

## Development Instructions

### 1. Install the dependencies

```bash
npm install
```

### 2. Startup development mode (HMR, error reporting, etc.)

```bash
npm run dev
```

The application will start in development mode on `http://localhost:8080`.

### 3. Lint & Format check

```bash
npm run lint
```

Formating is powered by **Prettier** and code quality is checked by **ESLint**.

### 4. Build for production compilation

```bash
npm run build
```

The optimized bundled build files will be generated under the `dist/spa` folder.

---

## Security Guidelines

- **XSS Prevention**: Since blog posts are authored in Jodit (HTML output) and rendered dynamically, the application explicitly processes all HTML views through `DOMPurify` before injecting them using `v-html`. Avoid using raw `v-html` directly without sanitizing.
- **Session Security**: The user authentication token is handled via Secure cookies using the `js-cookie` library. Unauthenticated routes redirect to `/login` automatically via synchronous router navigation guards.
