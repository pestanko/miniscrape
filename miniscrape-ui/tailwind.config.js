/** @type {import('tailwindcss').Config} */
export default {
	content: [
		'./src/**/*.{html,js,svelte,ts}',
		'./node_modules/flowbite-svelte/**/*.{html,js,svelte,ts}',
	],
	theme: {
		extend: {
			colors: {
				background: 'var(--background)',
				'accent-primary': 'var(--accent-primary)',
				'bg-primary': 'var(--bg-primary)',
				'accent-secondary': 'var(--accent-secondary)',
			},
			width: {
				38: '9.5rem',
				navigation: 'var(--navigation-width)',
				'navigation-thin': 'var(--navigation-width-thin)',
			},
			fontSize: {
				header: 'var(--header-font-size)',
			},
			height: {
				header: 'var(--header-height)',
			},
			flexBasis: {
				header: 'var(--header-height)',
			},
			gridTemplateAreas: {
				layout: ['nav header', 'nav main'],
			},
			gridTemplateColumns: {
				layout: 'var(--navigation-width) 1fr',
				'layout-thin': 'var(--navigation-width-thin) 1fr',
			},
			gridTemplateRows: {
				layout: 'var(--header-height) 1fr',
			},
			transitionProperty: {
				grid: 'grid',
				width: 'width',
			},
		},
		variables: {
			DEFAULT: {
				measure: '65ch', // max width of text for the best readability
				header: {
					height: '6.25rem',
					'font-size': 'clamp(0.6rem, 1.5vw, 1rem)',
				},
				navigation: {
					width: '16rem',
					'width-thin': '3.5rem',
				},
				background: '#f5f7fa',
				'accent-primary': 'var(--dt-color-base-700)',
				'bg-primary': 'var(--dt-color-neutral-200)',
				'accent-secondary': 'var(--dt-color-base-900)',
			},
		},
	},
	variants: {
		gridTemplateAreas: ['responsive'],
	},
	corePlugins: {
		aspectRatio: false,
	},
	plugins: [
		require('@mertasan/tailwindcss-variables'),
		require('@savvywombat/tailwindcss-grid-areas'),
		require('@tailwindcss/forms'),
		require('@tailwindcss/aspect-ratio'),
		require('@tailwindcss/container-queries'),
		require('tailwindcss-animate'),
		require('flowbite/plugin'),
	],
};
