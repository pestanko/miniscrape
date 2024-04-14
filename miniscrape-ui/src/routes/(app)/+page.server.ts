import { provideServerDeps } from '@src/server/deps';
import type { Load } from '@sveltejs/kit';


export const load: Load  = async ({}) => {
    const dp = await provideServerDeps();
    const { miniscrape } = dp.clients;

    const tags = await miniscrape.getTags();

	return {
        tags,
	};
}