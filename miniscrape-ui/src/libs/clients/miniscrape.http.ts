import type { AppConfig } from '@src/server/config';
import type { BaseDeps } from '@src/server/deps';
import { createHttpClient, type HttpClient } from './client.http';
import type { Logger } from 'pino';


export type PageSelector ={
    category: string;
    tags: string[];
};

export type CategoryResponse = {
    name: string;
    pages: string[];
    tags: string[];
};

export type PageDetail = {
    name: string;
    category: string;
    tags: string[];
    codename: string;
    homepage: string;
    url: string;
};

type ContentStatus = 'ok' | 'empty' | 'error';

export type PageContentResponse = {
    content: string;
    resolver: string;
    status: ContentStatus;
    page: PageDetail;
};

export class MiniScrapeHttpClient {
    private readonly baseUrl: string;
    private readonly httpClient: HttpClient;
    private readonly log: Logger;

    constructor(dp: Pick<BaseDeps, 'log' | 'httpClient'>, props: AppConfig['services']['miniscrape']) {
        this.baseUrl = props.internalUrl;
        this.log = dp.log;
        this.httpClient = createHttpClient({
            baseURL: this.baseUrl,
        });
    }

    async getCategories(): Promise<CategoryResponse[]> {
        const response = await this.httpClient.get<CategoryResponse[]>(this.baseUrl + '/api/v1/categories');
        this.log.debug({ data: response.data }, 'Categories loaded');
        return response.data;
    }

    async getTags(): Promise<string[]> {
        const categories = await this.getCategories();
        const foodCategory = categories.find(c => c.name === 'food');
        const tags = foodCategory?.tags ?? [];
        this.log.debug({  tags }, 'Tags loaded');
        return tags.sort();
    }

    async getPages(sel: PageSelector) {
        const queryBuilder = this.createQueryBuilder(sel);
        const fullUrl = this.baseUrl + '/api/v1/pages?' + queryBuilder;
        const response = await this.httpClient.get(fullUrl);
        const pages = response.data;
        this.log.debug({pages}, 'Pages loaded');
        return pages;
    }

    async getContent(sel: PageSelector): Promise<PageContentResponse[]> {
        const queryBuilder = this.createQueryBuilder(sel);
        const fullUrl = this.baseUrl + '/api/v1/content?' + queryBuilder;
        const response = await this.httpClient.get<PageContentResponse[]>(fullUrl);
        const pages = response.data;
        this.log.debug({ pages }, 'Pages content loaded');
        return pages.sort((fst, snd) => {
            const status = compareStatus(fst.status, snd.status);
            if(status !== 0) {
                return status;
            }

            return fst.page.codename.localeCompare(snd.page.codename);
        });
    }

    private createQueryBuilder(sel: PageSelector): URLSearchParams {
        const qb = new URLSearchParams();

        if(sel.category) {
            qb.set('c', sel.category);
        }

        if(sel.tags) {
            sel.tags.forEach(t => qb.append('t', t))
        }

        return qb;
    }
}

const compareStatus = (fst: ContentStatus, snd: ContentStatus) => {
    const fstNum = statusToNumber(fst);
    const sndNum = statusToNumber(snd);

    return fstNum - sndNum;
};

const statusToNumber = (status: ContentStatus) => {
    switch(status) {
        case 'ok':
            return 0;
        case 'empty':
            return 1;
        case 'error':
            return 2;
    }
};