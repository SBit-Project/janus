var SRC20 = artifacts.require("SRC20Token");

module.exports = async function(deployer) {
  await deployer.deploy(SRC20);
};