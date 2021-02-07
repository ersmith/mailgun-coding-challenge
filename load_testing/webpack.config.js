const { CleanWebpackPlugin } = require('clean-webpack-plugin');

module.exports = {
  mode: 'production',
  entry: './src/main.js',
  output: {
      path: __dirname + '/dist',
      filename: 'test.[name].js',
      libraryTarget: 'commonjs'
  },
  module: {
      rules: [
          { test: /\.js$/, use: 'babel-loader' },
      ]
  },
  stats: {
      colors: true,
      warnings: false
  },
  target: "web",
  externals: /k6(\/.*)?/,
  devtool: 'source-map',
  plugins: [
    new CleanWebpackPlugin(),
],
}
