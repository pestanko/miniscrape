import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vitest/config';


export default defineConfig({
	plugins: [sveltekit()],
	test: {
		include: ['src/**/*.{test,spec}.{js,ts}'],
		watch: false,
		coverage: {
			reporter: ['lcov', 'text'],
			exclude: ['**/*.spec.ts'],
		},
		reporters: ['default', 'junit', 'json'],
		outputFile: {
			junit: 'reports/test-results.xml',
			json: 'reports/test-results.json',
		},
		globals: true,
		setupFiles: [],
	},
	ssr: {
		external: ['reflect-metadata'],
	},
});
