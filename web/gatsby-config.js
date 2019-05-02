module.exports = {
  siteMetadata: {
    title: `Workflow`,
    description: `Get things done`,
    author: `@clarsen`,
  },
  plugins: [
    `gatsby-plugin-typescript`,
    `gatsby-plugin-tslint`,
    // `gatsby-plugin-react-helmet`,
    // `gatsby-transformer-sharp`,
    // `gatsby-plugin-sharp`,
    // `gatsby-plugin-emotion`,
    {
      resolve: `gatsby-plugin-create-client-paths`,
      options: { prefixes: [`/app/*`] },
    },
    {
      resolve: `gatsby-source-filesystem`,
      options: {
        name: `src`,
        path: `${__dirname}/src`,
      },
    },
    // {
    //   resolve: `gatsby-transformer-remark`,
    //   options: {
    //     plugins: [
    //       // gatsby-remark-relative-images must
    //       // go before gatsby-remark-images
    //       {
    //         resolve: `gatsby-remark-relative-images`,
    //       },
    //       {
    //         resolve: `gatsby-remark-images`,
    //         options: {
    //         // It's important to specify the maxWidth (in pixels) of
    //         // the content container as this plugin uses this as the
    //         // base for generating different widths of each image.
    //         // see src/components/layout.js
    //         maxWidth: 700,
    //         }
    //       }
    //     ]
    //   }
    // },
    // {
    //   resolve: `gatsby-plugin-typography`,
    //   options: {
    //     pathToConfigModule: `src/utils/typography`,
    //   },
    // },
    {
      resolve: `gatsby-plugin-manifest`,
      options: {
        name: `Workflow`,
        short_name: `Workflow`,
        start_url: `/`,
        background_color: `#663399`,
        theme_color: `#663399`,
        display: `minimal-ui`,
        icon: `src/images/icon.jpg`, // This path is relative to the root of the site.
      },
    },
    // this (optional) plugin enables Progressive Web App + Offline functionality
    // To learn more, visit: https://gatsby.dev/offline
    `gatsby-plugin-offline`,
    // {
    //   resolve: 'gatsby-plugin-sass',
    //   options: {
    //     includePaths: [`${__dirname}/node_modules`, `${__dirname}/src/`],
    //     precision: 8
    //   }
    // },
    `gatsby-plugin-netlify`, // make sure to put last in the array
  ],
}
