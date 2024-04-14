<script lang="ts">
	import { Accordion, AccordionItem, Badge } from 'flowbite-svelte';
	import { BookOpenIcon } from 'svelte-feather-icons';

	export let data;
</script>

<main>
	<h1 class="text-3xl mb-5">Pages</h1>

	<div>
		<Accordion>
			{#each data.pages as page}
				<AccordionItem>
					<span slot="header">
						<div class="flex flex-row gap-3 hover:underline">
							<BookOpenIcon size="17"/>{page.page.name}
                            <Badge color="{ page.status === 'ok' ? 'green' : 'pink' }">{page.status}</Badge>
						</div>
					</span>
					<div class="flex flex-row gap-3 mb-2">
						<Badge color="green" target="_blank" href={page.page.homepage}>Home Page</Badge>
						<span>Tags:</span>
						{#each page.page.tags as tag}
							<Badge href="/pages?t={tag}" color="yellow">{tag}</Badge>
						{/each}
					</div>

					<h3 class="text-xl mb-2">Daily Menu</h3>
					{#if page.status === 'ok'}
						<pre>
                        {page.content}
                    </pre>
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