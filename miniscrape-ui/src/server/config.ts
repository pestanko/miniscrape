import { env } from '$env/dynamic/private';
import { strToBool } from '@libs/string.utils';

export type AppConfig = ReturnType<typeof createAppConfig>;
export type CfgAvailableServices = keyof AppConfig['services'];

export const createAppConfig = () => {
    const environment = env.ENV_NAME || 'dev';
    return {
        environment,
        logger: {
            level: env.LOG_LEVEL || 'info',
            pretty: strToBool(env.LOG_PRETTY),
        },
        services: createServicesConfig(environment),
    };
};

const createServicesConfig = (envName: string) => {
    return {
        miniscrape: {
            internalUrl:
                env.SERVICE_MINISCRAPE_URL || 'http://localhost:8080'
        },
    };
};