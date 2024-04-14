<script type="ts">
	import Navigation from '@ui/navigation/Navigation.svelte';
	import Header from '@ui/header/Header.svelte';
	import { Toaster } from 'svelte-sonner';

</script>


<!-- Main app structure -->
<div class="app">
	<a class="app__logo" href="/">MiniScrape</a>
	<header class="app__top-nav" aria-label="Secondary Navigation"><Header /></header>
	<nav class="app__menu" aria-label="Menu"><Navigation /></nav>
	<main class="app__main grow">
		<Toaster />
		<slot />
	</main>
</div>

<style>
	.app {
		display: grid;
		grid-template-areas:
			'logo topnav'
			'nav main';
		grid-template-columns: auto 1fr;
		grid-template-rows: 60px 1fr;
		background-color: var(--dt-color-white);
		min-height: 100vh;
	}
	.app__logo {
		grid-area: logo;
		padding: 12px 24px;
		display: flex;
		align-items: center;
		border-right: 1px solid var(--dt-color-neutral-300);
		border-bottom: 1px solid var(--dt-color-neutral-300);
	}
	.app__top-nav {
		grid-area: topnav;
		display: flex;
		justify-content: flex-end;
		align-items: center;
		background-color: var(--dt-color-white);
		border-bottom: 1px solid var(--dt-color-neutral-300);
		padding: 0 24px;
	}
	.app__menu {
		grid-area: nav;
		padding: 24px;
		border-right: 1px solid var(--dt-color-neutral-300);
	}
	.app__main {
		grid-area: main;
		background-color: var(--dt-color-neutral-100);
		padding: 24px;
		min-width: 100vh;
	}
</style>