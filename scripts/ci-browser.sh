cd ./wrappers/js-web
npm i
npm run build:dist
npm run build:test
./node_modules/.bin/cypress install --force
npm run test
