<template>
  <div class="home">
    <h4>This is a simple demo on how to add Janus as a chain to Metamask</h4>
    <Web3Button
      @web3Connected="setWeb3Connected"
      @sbitConnected="setSbitConnected"
      :connected="web3Connected"
      :sbitConnected="sbitConnected"
    />
  </div>
</template>

<script>
// @ is an alias to /src
import Web3Button from '@/components/Web3Button.vue'

export default {
  name: 'Home',
  components: {
    Web3Button
  },
  data() {
    return {
      web3Connected: false,
      sbitConnected: false,
    }
  },
  methods: {
    setWeb3Connected: function(connected) {
      console.log("setWeb3Connected", connected);
      this.web3Connected = connected;
    },
    setSbitConnected: function(sbitConnected) {
      console.log("setSbitConnected", sbitConnected);
      console.log("TODO: periodically query chain id to detect if the user changes network")
      this.sbitConnected = sbitConnected;
    }
  }
}
</script>
