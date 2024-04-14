import { createPinoLogger } from '@libs/logger';
import { createAppConfig, type AppConfig } from './config';
import type { Logger } from 'pino';
import { createHttpClient, type HttpClient } from '@libs/clients/client.http';
import { MiniScrapeHttpClient } from '@libs/clients/miniscrape.http';

export type BaseDeps = {
    log: Logger;
    config: AppConfig;
    httpClient: HttpClient;
}

export const provideServerDeps = async () => {
    const config = createAppConfig();
    const log = createPinoLogger(config.logger);

    log.debug(config, 'Config has been loaded');

    const baseDeps = {
        config,
        log,
        httpClient: createHttpClient(),
    };

    const clients = {
        miniscrape: new MiniScrapeHttpClient(baseDeps, config.services.miniscrape),
    };

    return {
        ...baseDeps,
        clients,
    }
}