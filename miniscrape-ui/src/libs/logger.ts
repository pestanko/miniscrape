import pino, { type DestinationStream, type Logger, type LoggerOptions } from 'pino';

type LoggerParams = {
	stream?: DestinationStream;
	level: string;
	browser?: boolean;
};

export const createPinoLogger = (params: LoggerParams): Logger => {
	const { level, stream } = params;
	const redactionRules = defaultRedactionRules();

	const logConfig: LoggerOptions = {
		level,
		redact: redactionRules,
		serializers: pino.stdSerializers,
		formatters: {
			level: (level) => ({ level }),
		},
	};

	const logger = pino(logConfig, stream);
	return logger;
};

interface RedactionRules {
	paths: string[];
	censor: string;
}

const defaultRedactionRules = (): RedactionRules => {
	return {
		paths: [
			'email',
			'password',
			'username',
			'[*].email',
			'[*].password',
			'[*].username',
			'[*].token',
			'[*].secret',
		],
		censor: '[REDACTED]',
	};
};