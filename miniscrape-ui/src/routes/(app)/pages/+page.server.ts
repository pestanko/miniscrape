import { provideServerDeps } from '@src/server/deps';
import type { Load } from '@sveltejs/kit';


export const load: Load  = async ({ url }) => {
    const dp = await provideServerDeps();
    const { miniscrape } = dp.clients;
    const tags = url.searchParams.getAll('t');
    const pages = await miniscrape.getContent({
        category: 'food',
        tags,
    })
	return {
        pages,
	};
}