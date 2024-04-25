<script lang="ts">
	import type { PageContentResponse } from '@src/libs/clients/miniscrape.http';
	import { Accordion, AccordionItem, Badge, Button } from 'flowbite-svelte';
	import { BookOpenIcon } from 'svelte-feather-icons';

	export let data: {
		pages: PageContentResponse[] 
	} = {
		pages: []
	};

	type AccordionState = {
		isOpen: boolean;
	};

	type StateObj = { [key: string]: AccordionState };


	let states: { [key: string]: AccordionState } = data.pages.reduce((acc: StateObj, page) => {
		acc[page.page.codename as string] = { isOpen: false };
		return acc;
	}, {  });

	console.log("States", states);

	const setStateAll = (isOpen: boolean) => {
		Object.keys(states).forEach((key) => {
			states[key].isOpen = isOpen;
		});
	}

</script>

<main>
	<h1 class="text-3xl mb-5">Pages</h1>

	<div id="pages-menu" class="menu flex gap-4">
		<Button color="green" on:click={() => setStateAll(true)}>Expand All</Button>
		<Button color="yellow" on:click={() => setStateAll(false)}>Collapse All</Button>
	</div>

	<div id="pages-list">
		<Accordion color="light" multiple>
			{#each data.pages as page}
				<AccordionItem bind:open={states[page.page.codename].isOpen}>
					<span slot="header">
						<div class="flex flex-row gap-3 hover:underline">
							<BookOpenIcon size="17"/>{page.page.name}
                            <Badge color="{ page.status === 'ok' ? 'green' : 'pink' }">{page.status}</Badge>
						</div>
					</span>
					<div class="flex flex-row gap-3 mb-2">
						<Badge color="green" target="_blank" href={page.page.homepage}>Home Page</Badge>
						<Badge color="blue" target="_blank" href={page.page.url}>Daily Menu URL</Badge>
						<span>Tags:</span>
						{#each page.page.tags as tag}
							<Badge href="/pages?t={tag}" color="yellow">{tag}</Badge>
						{/each}
					</div>

					<h3 class="text-xl mb-2">Daily Menu</h3>
					{#if page.status === 'ok'}
						{#if page.resolver === 'pdf'}
						<embed src={page.content} type="application/pdf" width="100%" height="600px" />
						{:else if page.resolver === 'img'}
						<img src={page.content} alt="Daily Menu: {page.page.name}" />
						{:else if page.resolver === 'iframe'}
						<iframe src={page.content} width="100%" height="600px" title="Daily Menu: {page.page.name}" />
						{:else if page.resolver === 'url_only'}
						<Badge color="blue" target="_blank" href={page.page.url}>Daily Menu URL</Badge>
						{:else}
						<pre>
                        	{page.content}
                    	</pre>
						{/if}

						
					{:else}
						<p>
							Unable to load daily menu for the restaurant {page.page.name}. Please visit
							<a class="hover:underline" target="_blank" href={page.page.homepage}>{page.page.homepage}</a>.
						</p>
					{/if}
				</AccordionItem>
			{/each}
		</Accordion>
	</div>
</main>
