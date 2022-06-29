<template>
  <div class="hello">
    <div v-if="web3Detected">
      <b-button v-if="sbitConnected">Connected to SBIT</b-button>
      <b-button v-else-if="connected" v-on:click="connectToSbit()">Connect to SBIT</b-button>
      <b-button v-else v-on:click="connectToWeb3()">Connect</b-button>
    </div>
    <b-button v-else>No Web3 detected - Install metamask</b-button>
  </div>
</template>

<script>
let SBITMainnet = {
  chainId: '0x51', // 81
  chainName: 'SBIT Mainnet',
  rpcUrls: ['https://janus.siswap.com/api/'],
  blockExplorerUrls: ['https://mainnet.sbit.dev/'],
  iconUrls: [
    'https://mainnet.sbit.dev/images/metamask_icon.svg',
    'https://mainnet.sbit.dev/images/metamask_icon.png',
  ],
  nativeCurrency: {
    decimals: 18,
    symbol: 'SBIT',
  },
};
let SBITTestNet = {
    chainId: '0x22B9', // 8889
  chainName: 'SBIT Testnet',
  rpcUrls: ['https://testnet-janus.siswap.com/api/'],
  blockExplorerUrls: ['https://testnet.sbit.dev/'],
  iconUrls: [
    'https://mainnet.sbit.dev/images/metamask_icon.svg',
    'https://mainnet.sbit.dev/images/metamask_icon.png',
  ],
  nativeCurrency: {
    decimals: 18,
    symbol: 'SBIT',
  },
};
let SBITRegTest = {
  chainId: '0x22BA', // 8890
  chainName: 'SBIT Regtest',
  rpcUrls: ['https://localhost:22402'],
  // blockExplorerUrls: ['https://testnet.sbit.dev/'],
  iconUrls: [
    'https://mainnet.sbit.dev/images/metamask_icon.svg',
    'https://mainnet.sbit.dev/images/metamask_icon.png',
  ],
  nativeCurrency: {
    decimals: 18,
    symbol: 'SBIT',
  },
};
let config = {
  "0x51": SBITMainnet,
  "0x22B9": SBITTestNet,
  "0x22BA": SBITRegTest,
};

export default {
  name: 'Web3Button',
  props: {
    msg: String,
    connected: Boolean,
    sbitConnected: Boolean,
  },
  computed: {
    web3Detected: function() {
      return !!this.Web3;
    },
  },
  methods: {
    getChainId: function() {
      return window.sbit.chainId;
    },
    isOnSbitChainId: function() {
      let chainId = this.getChainId();
      return chainId == SBITMainnet.chainId || chainId == SBITTestNet.chainId;
    },
    connectToWeb3: function(){
      if (this.connected) {
        return;
      }
      let self = this;
      window.sbit.request({ method: 'eth_requestAccounts' })
        .then(() => {
          console.log("Emitting web3Connected event");
          let sbitConnected = self.isOnSbitChainId();
          let currentlySbitConnected = self.sbitConnected;
          self.$emit("web3Connected", true);
          if (currentlySbitConnected != sbitConnected) {
            console.log("ChainID matches SBIT, not prompting to add network to web3, already connected.");
            self.$emit("sbitConnected", true);
          }
        })
        .catch((e) => {
          console.log("Connecting to web3 failed", arguments, e);
        })
    },
    connectToSbit: function() {
      console.log("Connecting to Sbit, current chainID is", this.getChainId());

      let self = this;
      let sbitConfig = config[this.getChainId()] || SBITTestNet;
      console.log("Adding network to Metamask", sbitConfig);
      window.sbit.request({
        method: "wallet_addEthereumChain",
        params: [sbitConfig],
      })
        .then(() => {
          self.$emit("sbitConnected", true);
        })
        .catch(() => {
          console.log("Adding network failed", arguments);
        })
    },
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
