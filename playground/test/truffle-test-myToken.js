const MyToken = artifacts.require("MyToken");

contract("MyToken", async accounts => {

  it("Has been deployed", async () => {
    const myTokenDeployed = await MyToken.deployed();
    assert(myTokenDeployed, "contract has been deployed");
    console.log("Address is: ", myTokenDeployed.address)
  });

  it("should perform transactions correctly", async () => {
      
    const acc1 = accounts[0];
    console.log("acc1 address: ", acc1);
    const acc2 = accounts[1];
    console.log("acc2 address: ", acc2);

    const amount = 10;

    const myToken = await MyToken.deployed();

    let balance = await myToken.balanceOf(acc1);
    let acc1StartingBalance = balance.toNumber();
    console.log("starting balance of acc1: ", acc1StartingBalance);

    balance = await myToken.balanceOf(acc2);
    let acc2StartingBalance = balance.toNumber();
    console.log("starting balance of acc2: ", acc2StartingBalance);

    await myToken.mint(acc1, 100).then((receipt) => { console.log("receipt: ", receipt)});
    balance = await myToken.balanceOf(acc1);
    let acc1Balance = balance.toNumber();
    console.log("new balance of acc1: ", acc1Balance);

    await myToken.transfer(acc2, amount, {from: acc1}).then((receipt) => { console.log("receipt: ", receipt)});
    balance = await myToken.balanceOf(acc2);
    let acc2Balance = balance.toNumber();
    balance = await myToken.balanceOf(acc1);
    acc1Balance = balance.toNumber();
    console.log("new balance of acc1 after transfer to acc2: ", acc1Balance);
    console.log("new balance of acc2 after transfer to acc1: ", acc2Balance);

  });

});