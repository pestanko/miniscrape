import axios, { type CreateAxiosDefaults } from 'axios';

const DEFAULT_TIMEOUT = 60_000;

export type HttpClient = ReturnType<typeof createHttpClient>;

export const createHttpClient = (opts?: CreateAxiosDefaults) => {
    const httpClient = axios.create({
        timeout: opts?.timeout ?? DEFAULT_TIMEOUT,
        headers: {
            'X-App': 'miniscrape-ui',
        },
    });

    return httpClient;
}