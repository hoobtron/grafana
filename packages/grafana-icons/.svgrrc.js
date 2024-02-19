/**
 * Reference: https://react-svgr.com/docs/options/
 */
module.exports = {
  icon: true,
  jsxRuntime: 'automatic',
  outDir: './generated',
  template: require('./templates/icon'),
  svgoConfig: {
    plugins: [
      // Sanitise the SVGs
      'removeScriptElement',
    ],
  },
  jsx: {
    babelConfig: {
      plugins: [
        // Remove fill and id attributes from SVG child elements
        [
          '@svgr/babel-plugin-remove-jsx-attribute',
          {
            elements: ['path', 'g', 'clipPath'],
            attributes: ['id', 'fill'],
          },
        ],
      ],
    },
  },
};