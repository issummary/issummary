const webpack = require('webpack');

module.exports = {
  mode: 'production',
  entry: ['./src/index.tsx'],

  resolve: {
    // Add '.ts' and '.tsx' as resolvable extensions.
    extensions: ['.ts', '.tsx', '.js', '.json']
  },

  module: {
    rules: [
      // All files with a '.ts' or '.tsx' extension will be handled by 'awesome-typescript-loader'.
      {
        test: /\.tsx?$/,
        loader: ['ts-loader']
      }
    ]
  }
};
