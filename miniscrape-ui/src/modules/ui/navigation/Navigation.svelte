<script lang="ts">
	import { page } from '$app/stores';
	import { BookIcon, GithubIcon, HomeIcon, SettingsIcon, TagIcon } from 'svelte-feather-icons';

	type Route = {
		name: string;
		icon?: any;
		path: string;
		isDisabled?: boolean;
		items?: Route[];
	};

	type Section = {
		name: string;
		items: Route[];
	};

	let sections: Section[] = [
		{
			name: '',
			items: [{ name: 'Home', path: '/', icon: HomeIcon }],
		},
		{
			name: 'Pages: Food',
			items: [
				{ name: 'Content', path: '/pages', icon: BookIcon },
				{ name: 'Config', path: '/pages/config', icon: SettingsIcon },
			],
		},
		{
			name: 'About',
			items: [
				{
					name: 'GitHub',
					path: 'https://github.com/pestanko/miniscrape',
					icon: GithubIcon,
				},
			],
		},
	];

	$: isActive = (path: string) => path === $page.url.pathname;
</script>

<aside class="menu">
	{#each sections as section}
		<h3 class="menu__caption font-medium text-sm">{section.name}</h3>
		<ul class="menu__container">
			{#each section.items as item}
				<li class="menu__item">
					<a
						href={item.path}
						class:is-disabled={item.isDisabled}
						class:is-active={isActive(item.path)}
					>
						<svelte:component this={item.icon} size="16" />
						{item.name}
					</a>
				</li>
			{/each}
		</ul>
	{/each}
</aside>

<style>
	.menu__caption {
		color: var(--dt-color-neutral-700);
	}
	.menu__container {
		margin-top: 0.5em;
		margin-bottom: 1em;
	}
	.menu__container:last-child {
		margin-bottom: 0;
	}
	.menu__item {
		position: inherit;
	}
	.menu__item a {
		display: flex;
		align-items: center;
		padding: 0.5rem;
		font-size: 0.875rem;
		height: 32px;
		border-radius: var(--dt-border-radius-sm);
		gap: 8px;
	}
	.menu__item a:hover {
		background-color: var(--dt-color-neutral-200);
	}
	.menu__item a.is-disabled {
		color: var(--dt-color-neutral-400);
		pointer-events: none;
	}
	.menu__item a.is-active {
		background-color: var(--dt-color-neutral-200);
	}

	.menu__item a.is-disabled:hover {
		background-color: inherit;
	}
</style>
