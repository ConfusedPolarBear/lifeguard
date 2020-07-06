const path = require("path");

const VueLoaderPlugin = require('vue-loader/lib/plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
	entry: { index: path.resolve(__dirname, "web", "src", "app.js") },
	output: { path: path.resolve(__dirname, "web", "dist") },
	devtool: 'source-map',	// TODO: investigate alternatives - this slows down building but avoids eval() [evil() ?]
	module: {
		rules: [
		{
			test: /\.vue$/,
			loader: 'vue-loader'
		},
		{
			test: /\.js$/,
			loader: 'babel-loader'
		},
		{
			test: /\.css$/,
			use: [
				{ loader: 'style-loader', options: { injectType: 'linkTag' } },
				{ loader: 'file-loader' },
			],
		}
		],
	},
	plugins: [
		new VueLoaderPlugin(),
		new HtmlWebpackPlugin({
			inject: true,
			template: 'web/src/index.html'
		})
	],
	resolve: {
		extensions: ['.js'],
		alias: {
			'vue$': 'vue/dist/vue.runtime.esm.js',
			'vue-router$': 'vue-router/dist/vue-router.esm.js',
		}
	}
};
