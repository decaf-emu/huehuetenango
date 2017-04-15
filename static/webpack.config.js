const { resolve } = require('path');
const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const ExtractTextPlugin = require('extract-text-webpack-plugin');

const isDevelopmentMode = !(require('yargs').argv.p || false);

module.exports = (options = {}) => ({
  entry: {
    vendor: './src/vendor',
    index: './src/main.js',
  },
  output: {
    path: resolve(__dirname, 'dist'),
    filename: options.dev ? '[name].js' : '[name].js?[hash]',
    chunkFilename: '[id].js?[chunkhash]',
    publicPath: '/',
  },
  module: {
    rules: [
      {
        test: /\.vue$/,
        use: ['vue-loader'],
      },
      {
        test: /\.js$/,
        use: ['babel-loader'],
      },
      {
        test: /\.html$/,
        use: [
          {
            loader: 'html-loader',
            options: {
              root: resolve(__dirname, 'src'),
              attrs: ['img:src', 'link:href'],
            },
          },
        ],
      },
      {
        test: /\.css$/,
        use: ExtractTextPlugin.extract({
          fallback: 'style-loader',
          use: ['css-loader', 'postcss-loader'],
        }),
      },
      {
        test: /\.scss$/,
        use: ExtractTextPlugin.extract({
          fallback: 'style-loader',
          use: ['css-loader', 'sass-loader', 'postcss-loader'],
        }),
      },
      {
        test: /favicon\.png$/,
        use: [
          {
            loader: 'file-loader',
            options: {
              name: '[name].[ext]?[hash]',
            },
          },
        ],
      },
      {
        test: /\.(png|jpg|jpeg|gif|eot|ttf|woff|woff2|svg|svgz)(\?.+)?$/,
        exclude: /favicon\.png$/,
        use: [
          {
            loader: 'url-loader',
            options: {
              limit: 10000,
            },
          },
        ],
      },
    ],
  },
  plugins: [
    new webpack.optimize.CommonsChunkPlugin({
      names: ['vendor', 'manifest'],
    }),
    new HtmlWebpackPlugin({
      template: 'src/index.html',
    }),
    new ExtractTextPlugin('styles.css'),
    new webpack.DefinePlugin({
      __DEBUG__: JSON.stringify(isDevelopmentMode),
    }),
  ],
  devServer: {
    host: '127.0.0.1',
    port: 8010,
    historyApiFallback: true,
  },
  devtool: options.dev ? '#eval-source-map' : undefined,
});
