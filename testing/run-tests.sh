#!/bin/sh
#Test openzeppelin integration
make openzeppelin
#Test Aave integration
cd aave
npm run test
