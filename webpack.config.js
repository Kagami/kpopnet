const path = require("path");
const CleanWebpackPlugin = require("clean-webpack-plugin");
const ExtractTextPlugin = require("extract-text-webpack-plugin");
const HtmlWebpackPlugin = require("html-webpack-plugin");

module.exports = {
  entry: "./ts/index/index.tsx",
  resolve: {
    extensions: [".tsx", ".ts", ".js"],
  },
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: "ts-loader",
        exclude: /node_modules/,
      },
      {
        test: /\.css$/,
        use: ExtractTextPlugin.extract({use: "css-loader"}),
        exclude: /node_modules/,
      },
    ],
  },
  plugins: [
    new CleanWebpackPlugin(["dist"]),
    new ExtractTextPlugin("static/[contenthash:10].css"),
    new HtmlWebpackPlugin({
      template: "ts/index/index.html",
      favicon: "ts/index/favicon.ico",
    }),
  ],
  output: {
    path: path.resolve(__dirname, "dist"),
    filename: "static/[chunkhash:10].js",
    publicPath: "/",
  },
};
