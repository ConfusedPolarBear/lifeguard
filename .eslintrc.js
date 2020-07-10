module.exports = {
	'env': {
		'browser': true,
		'es2020': true
	},
	'extends': [
		'eslint:recommended',
		'plugin:vue/essential'
	],
	'parserOptions': {
		'ecmaVersion': 11,
		'sourceType': 'module'
	},
	'plugins': [
		'vue'
	],
	'rules': {
		'indent': [
			'error',
			'tab'
		],
		'linebreak-style': [
			'error',
			'unix'
		],
		'quotes': [
			'error',
			'single'
		],
		'semi': [
			'error',
			'always'
		]
	}
};
