const swaggerAutogen = require('swagger-autogen')()

const doc = {
  info: {
    title: 'Voltaserve Identity Provider',
  },
}
const outputFile = './swagger.json'
const endpointsFiles = ['./src/app.ts']

swaggerAutogen(outputFile, endpointsFiles, doc)
