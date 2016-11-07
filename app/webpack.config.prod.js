var debug = process.env.NODE_ENV !== "production";
var webpack = require('webpack');

module.exports = {
  context: __dirname,
  devtool: debug ? "inline-sourcemap" : null,
  entry: [ 'whatwg-fetch', "./main.js"],
  output: {
    path: __dirname + "/../web/js",
    filename: "app.js"
  },
  module: {
    loaders: [
        {
          test: /\.js$/,
          include: __dirname + "/",
          exclude: /node_modules/,
          loader: "babel-loader",
          query: {
            presets: ['es2015', 'react', 'stage-2']
          }
        }
    ]
  },
  plugins: debug ? [] : [
    new webpack.optimize.DedupePlugin(),
    new webpack.optimize.OccurenceOrderPlugin(),
    new webpack.optimize.UglifyJsPlugin({ mangle: true, sourcemap: false }),
    new webpack.DefinePlugin({
      'process.env':{
        'NODE_ENV': JSON.stringify('production')
      }
    }),
    new webpack.optimize.UglifyJsPlugin({
      compress:{
        warnings: true
      }
    })
  ]
};