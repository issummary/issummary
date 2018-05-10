const webpack = require('webpack');

module.exports = {
  mode: 'development',
  entry: ['react-hot-loader/patch', './src/index.tsx'],

  // Enable sourcemaps for debugging webpack's output.
  devtool: 'cheap-module-eval-source-map',

  resolve: {
    // Add '.ts' and '.tsx' as resolvable extensions.
    extensions: ['.ts', '.tsx', '.js', '.json']
  },

  module: {
    rules: [
      // All files with a '.ts' or '.tsx' extension will be handled by 'awesome-typescript-loader'.
      {
        test: /\.tsx?$/,
        loader: ['react-hot-loader/webpack', 'ts-loader']
      }
    ]
  },
  externals: [
    {
      react: 'React',
      'react-dom': 'ReactDOM',
      'material-ui': 'MaterialUI'
    }
  ],
  plugins: [new webpack.HotModuleReplacementPlugin()],

  devServer: {
    hot: true,
    historyApiFallback: true
  }
};
