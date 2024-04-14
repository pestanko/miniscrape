<script lang="ts">
	import { Card } from 'flowbite-svelte';

	export let data: {
		tags: string[];
	} = {
		tags: [],
	};
</script>

<main>
	<h1 class="text-3xl leading-10">Tags</h1>
	<div class="flex flex-row gap-2 flex-wrap">
		{#each data.tags as tag}
			<Card href="/pages?t={tag}">
				<h3 class="text-2xl">{tag}</h3>
				<p class="text-xs">Show pages with tag <strong>{tag}</strong></p>
			</Card>
		{/each}
	</div>
</main>
