/* eslint-env node */
import { resolve } from 'path';
import ExtractTextPlugin from 'extract-text-webpack-plugin';
import webpack from 'webpack';
import UglifyJSPlugin from 'uglifyjs-webpack-plugin';
import HtmlWebpackPlugin from 'html-webpack-plugin';

const isProduction = () => process.env.NODE_ENV === 'production';

const config = {
  entry: {
    vendor: './src/vendor',
    index: './src/main.js',
  },
  output: {
    path: resolve(__dirname, 'dist'),
    filename: isProduction() ? '[name].js?[hash]' : '[name].js',
    chunkFilename: '[id].js?[chunkhash]',
    publicPath: '/',
  },
  devServer: {
    host: '127.0.0.1',
    port: 8081,
    historyApiFallback: true,
  },
  devtool: !isProduction() ? 'inline-source-map' : undefined,
  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /node_modules(?!\/(vuikit))/,
        loader: 'babel-loader',
      },
      {
        test: /\.vue$/,
        use: 'vue-loader',
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
      __DEBUG__: JSON.stringify(!isProduction()),
    }),
  ],
};

if (isProduction()) {
  config.plugins.push(
    new webpack.LoaderOptionsPlugin({
      minimize: true,
      debug: false,
    }),
    new UglifyJSPlugin(),
  );
} else {
  config.plugins.push(new webpack.NamedModulesPlugin());
}

export default config;
