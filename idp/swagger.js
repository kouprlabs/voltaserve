const swaggerAutogen = require('swagger-autogen')()

const doc = {
  info: {
    version: '2.0.0',
    title: 'Voltaserve Identity Provider',
  },
}
const outputFile = './swagger.json'
const endpointsFiles = ['./src/app.ts']

swaggerAutogen(outputFile, endpointsFiles, doc)
