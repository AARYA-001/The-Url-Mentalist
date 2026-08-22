import { browser } from '$app/environment';
import { get, writable } from 'svelte/store';

type Theme = 'light' | 'dark';

const store = writable<Theme>('light');

function applyTheme(next: Theme): void {
  store.set(next);
  if (browser) {
    localStorage.setItem('theme', next);
    document.documentElement.classList.toggle('dark', next === 'dark');
  }
}

export const theme = {
  subscribe: store.subscribe,

  init(): void {
    if (!browser) return;
    const saved = localStorage.getItem('theme') as Theme | null;
    const t = saved ?? 'light';
    store.set(t);
    document.documentElement.classList.toggle('dark', t === 'dark');
  },

  toggle(event?: MouseEvent): void {
    if (!browser) return;
    const current = get(store);
    const next: Theme = current === 'dark' ? 'light' : 'dark';

    const doc = document as any;
    // Fallback if View Transitions API is not supported or reduced motion is preferred
    if (
      !doc.startViewTransition ||
      window.matchMedia('(prefers-reduced-motion: reduce)').matches
    ) {
      applyTheme(next);
      return;
    }

    // Get origin coordinates from mouse event or fallback to window center
    const x = event ? event.clientX : window.innerWidth / 2;
    const y = event ? event.clientY : window.innerHeight / 2;

    // Calculate maximum radius to cover the furthest corner of the viewport
    const endRadius = Math.hypot(
      Math.max(x, window.innerWidth - x),
      Math.max(y, window.innerHeight - y)
    );

    const transition = doc.startViewTransition(() => {
      applyTheme(next);
    });

    transition.ready.then(() => {
      const clipPath = [
        `circle(0px at ${x}px ${y}px)`,
        `circle(${endRadius}px at ${x}px ${y}px)`
      ];

      doc.documentElement.animate(
        {
          clipPath
        },
        {
          duration: 500,
          easing: 'ease-in-out',
          pseudoElement: '::view-transition-new(root)'
        }
      );
    });
  },
};
