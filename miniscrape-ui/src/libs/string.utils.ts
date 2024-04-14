// converts string to boolean
export const strToBool = (str?: string) => {
	const possibleValues = ['true', 'on', 'yes', '1'];
	return possibleValues.includes((str ?? '').toLowerCase());
};