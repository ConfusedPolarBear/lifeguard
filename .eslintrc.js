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
		'block-spacing': [ 'error' ],
		'brace-style': [ 'error' ],
		'indent': [ 'error', 'tab' ],
		'keyword-spacing': [ 'error' ],
		'linebreak-style': [ 'error', 'unix' ],
		'quotes': [ 'error', 'single' ],
		'semi': [ 'error', 'always' ],
		'space-before-blocks': [ 'error' ],
		'space-infix-ops': [ 'error' ],
		'spaced-comment': [ 'error', 'always', { 'block': { 'balanced': true } } ]
	}
};
